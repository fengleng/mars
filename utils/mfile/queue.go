package mfile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"github.com/gososy/sorpc/log"
)

//
// 先入先出队列
//
type QueueReaderConfig struct {
	MaxCacheItemCount uint32
	MaxCacheDataBytes uint32
}

func NewDefaultQueueReaderConfig() *QueueReaderConfig {
	return &QueueReaderConfig{
		MaxCacheItemCount: 8192,
		MaxCacheDataBytes: 64 * 1024 * 1024,
	}
}

type QueueReader struct {
	fileBase
	cfg          *QueueReaderConfig
	itemCount    uint32
	readCursor   uint32
	pendingItems []*Item
	finishFile   *os.File
	finishMap    map[uint32]bool
	recountIndex bool

	warnChan chan *WarnMsg
}

func NewQueueReader(name, dataPath string, seq uint64, cfg *QueueReaderConfig, warnChan chan *WarnMsg) *QueueReader {
	if cfg == nil {
		cfg = NewDefaultQueueReaderConfig()
	}
	if cfg.MaxCacheItemCount == 0 {
		cfg.MaxCacheItemCount = 8192
	}
	if cfg.MaxCacheDataBytes == 0 {
		cfg.MaxCacheDataBytes = 64 * 1024 * 1024
	}
	return &QueueReader{
		fileBase: fileBase{
			name:     name,
			dataPath: dataPath,
			seq:      seq,
		},
		cfg:       cfg,
		finishMap: map[uint32]bool{},
		warnChan:  warnChan,
	}
}

func (p *QueueReader) statItemCount() error {
	idxPath := p.indexFilePath()
	idxInfo, err := os.Stat(idxPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	idxSize := int(idxInfo.Size())
	c := uint32(idxSize / indexItemSize)
	if c < p.itemCount {
		return fmt.Errorf("index file truncated, origin %d, cur %d", p.itemCount, c)
	}
	p.itemCount = c
	return nil
}

func (p *QueueReader) Init() error {
	var err error
	idxPath := p.indexFilePath()
	datPath := p.dataFilePath()
	finishPath := p.finishFilePath()

	defer func() {
		if err != nil && p.warnChan != nil {
			p.warnChan <- &WarnMsg{
				Label: fmt.Sprintf("%s: init err %v", p.name, err),
			}
		}
	}()

	p.indexFile, err = os.OpenFile(idxPath, os.O_RDONLY, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	p.dataFile, err = os.OpenFile(datPath, os.O_RDONLY, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	p.finishFile, err = os.OpenFile(finishPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	err = p.statItemCount()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	err = p.loadFinishFile()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	err = p.fillIndexCache()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	err = p.fillDataCache()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (p *QueueReader) close() {
	if p.indexFile != nil {
		_ = p.indexFile.Close()
		p.indexFile = nil
	}
	if p.dataFile != nil {
		_ = p.dataFile.Close()
		p.dataFile = nil
	}
	if p.finishFile != nil {
		_ = p.finishFile.Close()
		p.finishFile = nil
	}
}

func (p *QueueReader) pop() (*Item, error) {
	if len(p.pendingItems) == 0 {
		if p.recountIndex {
			p.recountIndex = false

			err := p.statItemCount()
			if err != nil {
				log.Errorf("err:%v", err)
				p.recountIndex = true
				return nil, err
			}
		}
		if err := p.fillIndexCache(); err != nil {
			return nil, err
		}
		if len(p.pendingItems) == 0 {
			return nil, nil
		}
	}
	item := p.pendingItems[0]
	if item.size > 0 && len(item.Data) == 0 {
		if err := p.fillDataCache(); err != nil {
			return nil, err
		}
		if len(item.Data) == 0 {
			p.warnChan <- &WarnMsg{
				Label: "load data fail",
			}
			log.Warnf("item %+v load data fail", item)
			return nil, errors.New("load data fail")
		}
	}
	p.pendingItems = p.pendingItems[1:]
	return item, nil
}

func (p *QueueReader) fillIndexCache() error {
	f := p.indexFile
	if f == nil {
		return errors.New("file not open")
	}
	if p.dataCorruption {
		return errors.New("data corruption")
	}
	if p.readCursor < p.itemCount {
		need := int(p.cfg.MaxCacheItemCount) - len(p.pendingItems)
		if need <= 0 {
			return nil
		}
		// 跳过已处理的
		for p.readCursor < p.itemCount {
			if p.finishMap[p.readCursor] {
				p.readCursor++
			} else {
				break
			}
		}
		pending := int(p.itemCount - p.readCursor)
		if need > pending {
			need = pending
		}
		if need <= 0 {
			return nil
		}
		size := need * indexItemSize
		var buf = make([]byte, size)
		n, err := f.ReadAt(buf, int64(p.readCursor)*indexItemSize)
		if err != nil {
			if err == io.EOF {
				// n 可能 > 0
			} else {
				log.Errorf("err:%v", err)
				return err
			}
		}
		log.Debugf("%s seq %d: load index at %d with size %d, size %d",
			p.name, p.seq, p.readCursor*indexItemSize, size, n)
		if n <= 0 {
			return nil
		}
		n = n / indexItemSize
		b := binary.LittleEndian
		for i := 0; i < n; i++ {
			index := p.readCursor
			p.readCursor++
			if p.finishMap[index] {
				continue
			}
			ptr := buf[i*indexItemSize:]
			var begMarker, endMarker uint16
			begMarker = b.Uint16(ptr)
			ptr = ptr[2:]
			endMarker = b.Uint16(ptr[28:])
			if begMarker != itemBegin {
				p.dataCorruption = true
				return errors.New("invalid index item begin marker")
			}
			if endMarker != itemEnd {
				p.dataCorruption = true
				return errors.New("invalid index item end marker")
			}
			var idx Item
			idx.CreatedAt = b.Uint32(ptr)
			idx.CorpId = b.Uint32(ptr[4:])
			idx.AppId = b.Uint32(ptr[8:])
			idx.Hash = b.Uint32(ptr[12:])
			// idx.reserved1 = b.Uint32(ptr[16:])
			idx.offset = b.Uint32(ptr[20:])
			idx.size = b.Uint32(ptr[24:])
			idx.Index = index
			idx.Seq = p.seq
			if idx.size < dataItemExtraSize {
				p.dataCorruption = true
				return fmt.Errorf("invalid data size %d, min than data item extra size", idx.size)
			}
			p.pendingItems = append(p.pendingItems, &idx)
		}
	}
	return nil
}

func (p *QueueReader) fillDataCache() error {
	f := p.dataFile
	if f == nil {
		return errors.New("file not open")
	}
	if p.dataCorruption {
		return errors.New("data corruption")
	}
	maxBytes := p.cfg.MaxCacheDataBytes
	curCacheSize := uint32(0)
	left := uint32(math.MaxUint32)
	right := uint32(0)
	var i int
	var item *Item
	for i, item = range p.pendingItems {
		if len(item.Data) > 0 {
			curCacheSize += uint32(len(item.Data))
			if curCacheSize > maxBytes {
				return nil
			}
		} else if item.size > 0 {
			if item.offset < left {
				left = item.offset
			}
			r := item.offset + item.size
			if r > right {
				right = r
			}
			if left > right {
				panic("Unreachable")
			}
			size := right - left
			if size+curCacheSize >= maxBytes {
				break
			}
		}
	}
	if left >= right {
		return nil
	}
	size := right - left
	log.Infof("%s seq %d: load data range at %d - %d, size %d", p.name, p.seq, left, right, size)
	var buf = make([]byte, size)
	var cnt int64
	for cnt < int64(size) {
		n, err := f.ReadAt(buf[cnt:], int64(left)+cnt)
		if err != nil && err != io.EOF {
			log.Errorf("err:%v", err)
			return err
		}
		if n > 0 {
			cnt += int64(n)
		} else {
			log.Errorf("data truncated, required %d, got %d", size, cnt)
			return errors.New("data truncated")
		}
	}
	b := binary.LittleEndian
	for j := 0; j <= i; j++ {
		item = p.pendingItems[j]
		if len(item.Data) > 0 || item.size == 0 {
			continue
		}
		if item.offset < left || item.offset+item.size > right {
			continue
		}
		o := item.offset - left
		ptr := buf[o:]
		var begMarker, endMarker uint16
		begMarker = b.Uint16(ptr)
		item.Data = ptr[2 : item.size-2]
		endMarker = b.Uint16(ptr[item.size-2:])
		if begMarker != itemBegin {
			p.dataCorruption = true
			return errors.New("invalid data begin marker")
		}
		if endMarker != itemEnd {
			p.dataCorruption = true
			return errors.New("invalid data end marker")
		}
	}
	return nil
}

func (p *QueueReader) loadFinishFile() error {
	if p.finishFile == nil {
		panic("Unreachable")
	}
	const (
		loadBufSize = 10240 * 4
	)
	buf := make([]byte, loadBufSize)
	_, err := p.finishFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Errorf("seek err:%v", err)
		return err
	}
	for {
		n, err := p.finishFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("err:%v", err)
			return err
		}
		if n > 0 {
			if n%4 != 0 {
				log.Warnf("invalid size of finish file read return %d", n)
				p.dataCorruption = true
				return errors.New("invalid size of finish file")
			}
			b := binary.LittleEndian
			for i := 0; i < n; i += 4 {
				idx := b.Uint32(buf[i:])
				p.finishMap[idx] = true
			}
		}
	}
	return nil
}

func batchWriteFinishFile(finishFile *os.File, finishMap map[uint32]bool, idxList []uint32, writeBuf []byte) error {
	if finishFile == nil {
		return errors.New("finish file not open")
	}
	var buf []byte
	if len(writeBuf) > 0 {
		buf = writeBuf
	} else {
		const (
			writeBufSize = 10240 * 4
		)
		buf = make([]byte, writeBufSize)
	}
	i := 0
	b := binary.LittleEndian
	for i < len(idxList) {
		j := 0
		m := map[uint32]bool{}
		for ; i < len(idxList); i++ {
			idx := idxList[i]
			if (finishMap != nil && finishMap[idx]) || m[idx] {
				continue
			}
			m[idx] = true
			b.PutUint32(buf[j:], idx)
			j += 4
			if j >= len(buf) {
				break
			}
		}
		if j > 0 {
			tmp := buf[:j]
			n, err := finishFile.Write(tmp)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			for n < len(tmp) {
				tmp = tmp[n:]
				n, err = finishFile.Write(tmp)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
			}
		}
		// if finishMap != nil {
		//	for k := range m {
		//		finishMap[k] = true
		//	}
		// }
	}
	return nil
}

type finishReq struct {
	seq   uint64
	index uint32
	hash  uint32
}

type QueueGroupReader struct {
	name            string
	dataPath        string
	reader          *QueueReader
	minSeq          uint64
	maxSeq          uint64
	cfg             *QueueReaderConfig
	hasInit         bool
	initMu          sync.Mutex
	exitChan        chan int
	notifyWriteChan chan uint64
	msgChan         chan *Item
	msg             *Item
	errChan         chan error
	warnChan        chan *WarnMsg
	finishChan      chan *finishReq
	closeChan       chan int

	ReportQueueLenFunc func(ql int64)
}

func (p *QueueGroupReader) GetReaderNum() int {
	if p == nil {
		return 0
	}

	// 并发模式下的，reader只有1个
	return 1
}

func NewQueueGroupReader(name, dataPath string, cfg *QueueReaderConfig, msgChanSize uint32) *QueueGroupReader {
	if cfg == nil {
		cfg = NewDefaultQueueReaderConfig()
	}
	if msgChanSize <= 0 {
		msgChanSize = 8000
	}
	return &QueueGroupReader{
		name:            name,
		dataPath:        dataPath,
		reader:          nil,
		maxSeq:          0,
		cfg:             cfg,
		exitChan:        make(chan int),
		notifyWriteChan: make(chan uint64, 1000),
		msgChan:         make(chan *Item, msgChanSize),
		errChan:         make(chan error, 1000),
		warnChan:        make(chan *WarnMsg, 100),
		finishChan:      make(chan *finishReq, msgChanSize),
		closeChan:       make(chan int),
	}
}

func (p *QueueGroupReader) getQueueLen() int64 {
	r := p.reader
	if r == nil {
		return 0
	}
	cur := r.readCursor
	t := r.itemCount
	var d uint32
	if t > cur {
		d = t - cur
	}
	return int64(len(p.msgChan) + len(r.pendingItems) + int(d))
}

func (p *QueueGroupReader) ConfirmAsyncMsg(_ string) {
}

func (p *QueueGroupReader) GetAsyncMsgConfirmChan() chan string {
	return nil
}

func (p *QueueGroupReader) SetLastFinishTs(seq uint64, hash uint32) {
}

func (p *QueueGroupReader) openReader() error {
	if p.reader != nil {
		p.reader.close()
		p.reader = nil
	}
	if p.minSeq == 0 {
		return nil
	}
	p.reader = NewQueueReader(p.name, p.dataPath, p.minSeq, p.cfg, p.warnChan)
	err := p.reader.Init()
	if err != nil {
		p.reader.close()
		p.reader = nil
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (p *QueueGroupReader) fillMsgChan() error {
	for {
		var msg *Item
		if p.msg != nil {
			msg = p.msg
			p.msg = nil
		} else {
			if p.reader != nil {
				var err error
				msg, err = p.reader.pop()
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				if msg == nil {
					if p.minSeq < p.maxSeq {
						p.minSeq++
						err = p.openReader()
						if err != nil {
							log.Errorf("err:%v", err)
							return err
						}
						if p.reader != nil {
							continue
						}
					}
				}
			} else {
				err := p.openReader()
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				if p.reader != nil {
					continue
				}
			}
		}
		if msg == nil {
			// 切换文件
			break
		}
		select {
		case p.msgChan <- msg:
		default:
			p.msg = msg
			goto OUT
		}
	}
OUT:
	return nil
}

func (p *QueueGroupReader) sendErr(err error) {
	select {
	case p.errChan <- err:
	default:
	}
}

func (p *QueueGroupReader) ioLoop() {
	log.Infof("start io loop")

	const (
		writeBufSize = 10240 * 4
	)
	writeBuf := make([]byte, writeBufSize)

	// fill msg chan
	err := p.fillMsgChan()
	if err != nil {
		log.Errorf("err:%v", err)
		p.sendErr(err)
	}
	flushFinish := func(r *finishReq) {
		type wrap struct {
			seq  uint64
			list []uint32
		}
		all := map[uint64]*wrap{}
		var last *wrap
		for r != nil {
			if last == nil || last.seq != r.seq {
				last = all[r.seq]
				if last == nil {
					last = &wrap{
						seq: r.seq,
					}
					all[r.seq] = last
				}
			}
			last.list = append(last.list, r.index)
			select {
			case r = <-p.finishChan:
			default:
				r = nil
			}
		}
		for _, v := range all {
			if p.reader != nil && p.reader.seq == v.seq {
				err := batchWriteFinishFile(p.reader.finishFile, p.reader.finishMap, v.list, writeBuf)
				if err != nil {
					log.Errorf("err:%v", err)
					p.sendErr(err)
				}
			} else {
				fb := fileBase{
					name:     p.name,
					dataPath: p.dataPath,
					seq:      v.seq,
				}
				finishPath := fb.finishFilePath()
				finishFile, err := os.OpenFile(finishPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					log.Errorf("err:%v", err)
					p.sendErr(err)
				} else {
					err = batchWriteFinishFile(finishFile, nil, v.list, writeBuf)
					if err != nil {
						log.Errorf("err:%v", err)
						p.sendErr(err)
					}
					_ = finishFile.Close()
				}
			}
		}
	}

	repTk := time.NewTicker(time.Second * 5)
	defer repTk.Stop()

	rep := func() {
		f := p.ReportQueueLenFunc
		if f != nil {
			r := p.reader
			if r != nil {
				f(p.getQueueLen())
			} else {
				f(0)
			}
		}
	}

	for {
		if p.msg != nil {
			select {
			case p.msgChan <- p.msg:
				p.msg = nil
			case <-p.exitChan:
				goto EXIT
			case r := <-p.finishChan:
				flushFinish(r)
			case <-repTk.C:
				rep()
			}
		} else {
			select {
			case seq := <-p.notifyWriteChan:
				if p.minSeq == 0 {
					p.minSeq = 1
				}
				for {
					if p.reader != nil {
						if p.minSeq == seq && p.reader != nil {
							p.reader.recountIndex = true
						}
					}
					if seq > p.maxSeq {
						p.maxSeq = seq
					}
					select {
					case seq = <-p.notifyWriteChan:
					default:
						goto FILL
					}
				}
			case <-p.exitChan:
				goto EXIT
			case r := <-p.finishChan:
				flushFinish(r)
			case <-repTk.C:
				rep()
			}
		}
	FILL:
		err = p.fillMsgChan()
		if err != nil {
			log.Errorf("err:%v", err)
			p.sendErr(err)
		}
	}
EXIT:
	select {
	case p.closeChan <- 1:
	default:
	}
}

func (p *QueueGroupReader) Init() error {
	p.initMu.Lock()
	defer p.initMu.Unlock()
	if p.hasInit {
		panic("duplicate init")
	}
	err := scanInitFileGroup(p.name, p.dataPath, &p.minSeq, &p.maxSeq)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if p.minSeq > p.maxSeq {
		panic("Unreachable")
	}
	for p.minSeq > 0 && p.minSeq <= p.maxSeq {
		p.reader = NewQueueReader(p.name, p.dataPath, p.minSeq, p.cfg, p.warnChan)
		err = p.reader.Init()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		if len(p.reader.pendingItems) > 0 {
			break
		}
		if p.minSeq < p.maxSeq {
			// 当前的文件已经处理完了
			log.Infof("seq %d has all finish, move to next", p.minSeq)
			p.reader.close()
			p.reader = nil
			p.minSeq++
		} else {
			// 已经是最后一个当前的文件了
			break
		}
	}
	go p.ioLoop()
	p.hasInit = true
	return nil
}

func (p *QueueGroupReader) Close() {
	select {
	case p.exitChan <- 1:
	default:
	}
}

func (p *QueueGroupReader) NotifyWrite(seq uint64) {
	select {
	case p.notifyWriteChan <- seq:
	default:
	}
}

func (p *QueueGroupReader) GetMsgChan() chan *Item {
	return p.msgChan
}

func (p *QueueGroupReader) GetErrChan() chan error {
	return p.errChan
}

func (p *QueueGroupReader) GetWarnChan() chan *WarnMsg {
	return p.warnChan
}

func (p *QueueGroupReader) GetCloseChan() chan int {
	return p.closeChan
}

func (p *QueueGroupReader) GetQueueLen() int64 {
	return p.getQueueLen()
}

func (p *QueueGroupReader) SetFinish(item *Item) {
	p.finishChan <- &finishReq{
		seq:   item.Seq,
		index: item.Index,
	}
}
