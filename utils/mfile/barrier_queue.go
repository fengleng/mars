package mfile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gososy/sorpc/log"
	"github.com/gososy/sorpc/utils"
	"github.com/huandu/skiplist"
	"github.com/streadway/handy/atomic"
	"io"
	"os"
	"sync"
	"time"
)

//
// 串行模式下的先入先出队列
//
type RetryItem struct {
	item      *Item
	retryAtMs uint64
}

const MaxPriority = 3
const PriorityMask = 0x3

type PriorityList struct {
	ll [MaxPriority][]*Item
}

type IndexCacheItemList struct {
	list      []*Item
	pl        *PriorityList
	retryList *utils.MinHeap
}

type IndexCache struct {
	minIndexList *skiplist.SkipList
	hashList     map[uint32]*IndexCacheItemList

	itemTotalCnt int
}

func NewIndexCache() *IndexCache {
	return &IndexCache{
		minIndexList: skiplist.New(skiplist.Uint32),
		hashList:     map[uint32]*IndexCacheItemList{},
	}
}

func (p *IndexCacheItemList) appendItem(item *Item) {
	if item.Priority == 0 {
		p.list = append(p.list, item)
	} else {
		pr := (item.Priority & PriorityMask) - 1
		if p.pl == nil {
			p.pl = &PriorityList{}
		}
		p.pl.ll[pr] = append(p.pl.ll[pr], item)
	}
}

func (p *IndexCacheItemList) popItem(check func(item *Item) bool) *Item {
	if p.pl != nil {
		for i := MaxPriority - 1; i >= 0; i-- {
			x := p.pl.ll[i]
			if len(x) > 0 {
				f := x[0]
				if check == nil || check(f) {
					x = x[1:]
					p.pl.ll[i] = x
					return f
				} else {
					// 严格先跑高优的
					return nil
				}
			}
		}
	}

	if len(p.list) > 0 {
		f := p.list[0]
		if check == nil || check(f) {
			p.list = p.list[1:]
			return f
		}
	}

	return nil
}

func (p *IndexCacheItemList) empty() bool {
	if p.pl != nil {
		for i := MaxPriority - 1; i >= 0; i-- {
			x := p.pl.ll[i]
			if len(x) > 0 {
				return false
			}
		}
	}
	return len(p.list) == 0
}

func (p *IndexCacheItemList) peekItem(check func(item *Item) bool) *Item {
	if p.pl != nil {
		for i := MaxPriority - 1; i >= 0; i-- {
			x := p.pl.ll[i]
			if len(x) > 0 {
				f := x[0]
				if check == nil || check(f) {
					return f
				}
			}
		}
	}

	if len(p.list) > 0 {
		f := p.list[0]
		if check == nil || check(f) {
			return f
		}
	}

	return nil
}

func (p *IndexCache) addItems(items []*Item) {
	m := p.minIndexList
	h := p.hashList
	for _, v := range items {
		x := h[v.Hash]
		if x == nil {
			x = &IndexCacheItemList{}
			m.Set(v.Index, v.Hash)
		}
		x.appendItem(v)
		h[v.Hash] = x
	}
	p.itemTotalCnt += len(items)
}

func (p *IndexCache) retry(item *Item, delayMs uint32) {
	m := p.minIndexList
	h := p.hashList
	x := h[item.Hash]
	if x == nil {
		x = &IndexCacheItemList{}
		h[item.Hash] = x
		m.Set(item.Index, item.Hash)
	}
	if x.retryList == nil {
		x.retryList = utils.NewMinHeap()
	}
	retryAtMs := utils.NowMs() + uint64(delayMs)
	x.retryList.PushEl(&utils.MinHeapElement{
		Value: &RetryItem{
			item:      item,
			retryAtMs: retryAtMs,
		},
		Priority: int64(retryAtMs),
	})
}

func (p *IndexCache) popSkipHash(skipHash map[uint32]int, barrierCount int, checkCount bool, delay *BarrierQueueReaderDelay, wm **WarnMsg) *Item {
	m := p.minIndexList
	h := p.hashList
	c := m.Front()
	var nowMs uint64
	for ; c != nil; c = c.Next() {
		hash := c.Value.(uint32)
		if checkCount && skipHash != nil {
			count := skipHash[hash]
			if count >= barrierCount {
				continue
			}
		}
		x := h[hash]
		if x == nil {
			*wm = &WarnMsg{
				Label: "pop len(list) == 0",
			}
			return nil
		}
		// 看看 retry list
		var item *Item
		if x.retryList != nil {
			if nowMs == 0 {
				nowMs = utils.NowMs()
			}
			top := x.retryList.PeekEl()
			if top.Priority > int64(nowMs) {
				continue
			}
			item = x.retryList.PopEl().Value.(*RetryItem).item
			if x.retryList.Len() == 0 {
				x.retryList = nil
			}
		} else {
			item = x.popItem(func(item *Item) bool {
				if delay.enableDelay && item.DelayType == DelayTypeRelate {
					if nowMs == 0 {
						nowMs = utils.NowMs()
					}
					var last uint64
					if delay.hash2LastFinishTs != nil {
						delay.hash2LastFinishTsMu.RLock()
						last = delay.hash2LastFinishTs[item.Hash]
						delay.hash2LastFinishTsMu.RUnlock()
					}
					if last+uint64(item.DelayValue) > nowMs {
						return false
					}
				}
				return true
			})
		}
		if item != nil {
			m.Remove(c.Key())
			if x.empty() && x.retryList == nil {
				delete(h, hash)
			} else {
				if x.retryList != nil {
					top := x.retryList.PeekEl()
					m.Set(top.Value.(*RetryItem).item.Index, hash)
				} else {
					f := x.peekItem(nil)
					if f != nil {
						m.Set(f.Index, hash)
					}
				}
			}
		}
		if item != nil {
			return item
		}
	}
	return nil
}
func (p *IndexCache) dump() {
	log.Infof("skip list len %d, total %d", p.minIndexList.Len(), p.itemTotalCnt)
}

type BarrierQueueReaderConfig struct {
	MaxCacheIndexFileCount int
}

func NewDefaultBarrierQueueReaderConfig() *BarrierQueueReaderConfig {
	return &BarrierQueueReaderConfig{
		MaxCacheIndexFileCount: 15,
	}
}

type BarrierQueueReaderDelay struct {
	enableDelay         bool
	hash2LastFinishTs   map[uint32]uint64
	hash2LastFinishTsMu sync.RWMutex
}

type BarrierQueueReader struct {
	fileBase
	cfg          *BarrierQueueReaderConfig
	itemCount    uint32
	readCursor   uint32
	finishFile   *os.File
	finishMap    map[uint32]bool
	recountIndex bool
	indexes      *IndexCache
	popCnt       int
	finishedCnt  atomic.Int
	delay        BarrierQueueReaderDelay

	indexBuf []byte

	openAt uint32
}

func NewBarrierQueueReader(name, dataPath string, seq uint64, cfg *BarrierQueueReaderConfig) *BarrierQueueReader {
	if cfg == nil {
		cfg = NewDefaultBarrierQueueReaderConfig()
	}
	return &BarrierQueueReader{
		fileBase: fileBase{
			name:     name,
			dataPath: dataPath,
			seq:      seq,
		},
		cfg:       cfg,
		finishMap: map[uint32]bool{},
		indexes:   NewIndexCache(),
	}
}

func (p *BarrierQueueReader) statItemCount() error {
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

func (p *BarrierQueueReader) dump() {
	log.Infof("%s.%d: pop %d fin cnt %d total %d sk %d pending %d finished %v",
		p.name, p.seq, p.popCnt, p.finishedCnt,
		p.indexes.itemTotalCnt, p.indexes.minIndexList.Len(),
		p.indexes.itemTotalCnt-p.popCnt, p.isFinished())
}

func isFileNotFoundError(err error) bool {
	return os.IsNotExist(err)
}

func (p *BarrierQueueReader) Init() error {
	var err error
	idxPath := p.indexFilePath()
	datPath := p.dataFilePath()
	finishPath := p.finishFilePath()
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

	p.openAt = utils.Now()
	return nil
}

func (p *BarrierQueueReader) close() {
	log.Infof("%s.%d: close", p.name, p.seq)

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

func ReaderTest() {
	for {
		r := NewBarrierQueueReader("hello", "/home/pinfire/smq", 1, nil)
		err := r.Init()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		if !r.isFinished() {
			log.Infof("not finished")
			return
		}
		r.close()
	}
}

func (p *BarrierQueueReader) fillIndexCache() error {
	f := p.indexFile
	if f == nil {
		return errors.New("file not open")
	}
	if p.dataCorruption {
		return errors.New("data corruption")
	}
	if p.readCursor < p.itemCount {
		// 跳过已处理的
		for p.readCursor < p.itemCount {
			if p.finishMap[p.readCursor] {
				p.readCursor++
			} else {
				break
			}
		}
		if int(p.itemCount-p.readCursor) <= 0 {
			log.Infof("skip load index, all finished")
			return nil
		}

		// init buffer
		if len(p.indexBuf) == 0 {
			log.Infof("alloc index buf for %s.%d", p.name, p.seq)
			p.indexBuf = make([]byte, (1024*1024)/indexItemSize*indexItemSize)
		}

		eof := false
		for p.readCursor < p.itemCount && !eof {
			n, err := f.ReadAt(p.indexBuf, int64(p.readCursor)*indexItemSize)
			if err != nil {
				if err == io.EOF {
					eof = true
				} else {
					log.Errorf("err:%v", err)
					return err
				}
			}
			log.Infof("%s seq %d: load index at %d with size %d, size %d ret",
				p.name, p.seq, p.readCursor*indexItemSize, len(p.indexBuf), n)
			if n <= 0 {
				break
			}
			n = n / indexItemSize
			b := binary.LittleEndian
			var items []*Item
			for i := 0; i < n; i++ {
				index := p.readCursor
				p.readCursor++
				if p.finishMap[index] {
					continue
				}
				ptr := p.indexBuf[i*indexItemSize:]
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
				idx.DelayType, idx.DelayValue, idx.Priority = unpackMisc(b.Uint32(ptr[16:]))
				idx.offset = b.Uint32(ptr[20:])
				idx.size = b.Uint32(ptr[24:])
				idx.Index = index
				idx.Seq = p.seq
				if idx.size < dataItemExtraSize {
					p.dataCorruption = true
					return fmt.Errorf("invalid data size %d, min than data item extra size", idx.size)
				}
				if idx.DelayType == DelayTypeRelate {
					p.delay.enableDelay = true
				}
				items = append(items, &idx)
			}
			if len(items) > 0 {
				p.indexes.addItems(items)
			}
		}
	}
	return nil
}

func (p *BarrierQueueReader) loadFinishFile() error {
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
	log.Infof("finished len %d", len(p.finishMap))
	return nil
}

func (p *BarrierQueueReader) readData(item *Item) error {
	d := p.dataFile
	if d == nil {
		return fmt.Errorf("invalid err")
	}
	dataBuf := make([]byte, item.size)
	var read uint32
	for read < item.size {
		n, err := d.ReadAt(dataBuf[read:], int64(item.offset)+int64(read))
		if err != nil {
			if err == io.EOF {
				return errors.New("data file truncated")
			}
			log.Errorf("err:%v", err)
			return err
		} else if n > 0 {
			read += uint32(n)
		} else {
			return errors.New("data file truncated")
		}
	}
	ptr := dataBuf[:]
	b := binary.LittleEndian
	begMarker := b.Uint16(ptr)
	item.Data = ptr[2 : item.size-2]
	endMarker := b.Uint16(ptr[item.size-2:])
	if begMarker != itemBegin {
		p.dataCorruption = true
		return errors.New("invalid data begin marker")
	}
	if endMarker != itemEnd {
		p.dataCorruption = true
		return errors.New("invalid data end marker")
	}
	return nil
}

func (p *BarrierQueueReader) retry(item *Item, delayMs uint32, wm **WarnMsg) {
	if p.popCnt <= 0 {
		*wm = &WarnMsg{
			Label: fmt.Sprintf("%s.%d: invalid pop cnt", p.name, p.seq),
		}
		return
	}
	p.popCnt--
	p.indexes.retry(item, delayMs)
}

func (p *BarrierQueueReader) popSkipHash(skipHash map[uint32]int, barrierCount int, checkCount bool, wm **WarnMsg) (*Item, error) {
	if p.recountIndex {
		p.recountIndex = false

		err := p.statItemCount()
		if err != nil {
			log.Errorf("err:%v", err)
			p.recountIndex = true
			return nil, err
		}
		err = p.fillIndexCache()
		if err != nil {
			log.Errorf("err:%v", err)
			p.recountIndex = true
			return nil, err
		}
	}

	item := p.indexes.popSkipHash(skipHash, barrierCount, checkCount, &p.delay, wm)
	if item == nil {
		return nil, nil
	}
	p.popCnt++

	err := p.readData(item)
	if err != nil {
		log.Errorf("err:%v", err)
		// drop msg
		p.finishedCnt.Add(1)
		return nil, err
	}
	return item, nil
}

func (p *BarrierQueueReader) isFinished() bool {
	if len(p.indexes.hashList) == 0 {
		res := (p.popCnt <= int(p.finishedCnt.Get())) && !p.recountIndex
		if res {
			return true
		}
	}
	return false
}

func (p *BarrierQueueReader) getQueueLen() int64 {
	var l int64
	t := p.indexes.itemTotalCnt
	pc := p.popCnt
	if t > pc {
		l = int64(t - pc)
	}
	return l
}

type BarrierQueueGroupReader struct {
	name                     string
	dataPath                 string
	readers                  []*BarrierQueueReader
	minSeq                   uint64
	maxSeq                   uint64
	cfg                      *BarrierQueueReaderConfig
	hasInit                  bool
	initMu                   sync.Mutex
	exitChan                 chan int
	notifyWriteChan          chan uint64
	errChan                  chan error
	warnChan                 chan *WarnMsg
	finishChan               chan *finishReq
	closeChan                chan int
	notifyWriteChan4Consumer chan uint64
	asyncMsgConfirmChan      chan string
}

func NewBarrierQueueGroupReader(name, dataPath string, cfg *BarrierQueueReaderConfig) *BarrierQueueGroupReader {
	if cfg == nil {
		cfg = NewDefaultBarrierQueueReaderConfig()
	}
	return &BarrierQueueGroupReader{
		name:                     name,
		dataPath:                 dataPath,
		maxSeq:                   0,
		cfg:                      cfg,
		exitChan:                 make(chan int),
		notifyWriteChan:          make(chan uint64, 1000),
		errChan:                  make(chan error, 1000),
		warnChan:                 make(chan *WarnMsg, 100),
		finishChan:               make(chan *finishReq, 1000),
		closeChan:                make(chan int),
		notifyWriteChan4Consumer: make(chan uint64, 1000),
		asyncMsgConfirmChan:      make(chan string, 1000),
	}
}

func (p *BarrierQueueGroupReader) GetAsyncMsgConfirmChan() chan string {
	return p.asyncMsgConfirmChan
}

func (p *BarrierQueueGroupReader) ConfirmAsyncMsg(msgId string) {
	p.asyncMsgConfirmChan <- msgId
}

func (p *BarrierQueueGroupReader) cleanFinishedReader() {
	if len(p.readers) <= 1 {
		return
	}
	var n []*BarrierQueueReader
	closeCnt := 0
	last := p.readers[len(p.readers)-1]
	for i := 0; i < len(p.readers)-1; i++ {
		skipCheckFinished := false

		if i+1 == len(p.readers)-1 {
			// 这个文件表示 max seq - 1，当快速切换文件时， max seq - 1 的文件，
			// 可能写入了文件，但是 delay 了一会会 notify 过来
			// 导致新写入的 msg 没有被消费，就 close 掉了
			if last.openAt+30 > utils.Now() {
				skipCheckFinished = true
			}
		}
		r := p.readers[i]
		if !skipCheckFinished && r.isFinished() {
			r.close()
			closeCnt++
		} else {
			n = append(n, r)
		}
	}
	if closeCnt > 0 {
		n = append(n, p.readers[len(p.readers)-1])
		p.readers = n
	}
}

func (p *BarrierQueueGroupReader) isReaderOpened(seq uint64) bool {
	for _, v := range p.readers {
		if v.seq == seq {
			return true
		}
	}
	return false
}

func (p *BarrierQueueGroupReader) openReader() error {
	p.cleanFinishedReader()

	if p.minSeq == 0 {
		return nil
	}
	for p.minSeq <= p.maxSeq && len(p.readers) < p.cfg.MaxCacheIndexFileCount {
		if !p.isReaderOpened(p.minSeq) {
			reader := NewBarrierQueueReader(p.name, p.dataPath, p.minSeq, p.cfg)
			err := reader.Init()
			if err != nil {
				reader.close()
				reader = nil
				log.Errorf("err:%v", err)
				if !isFileNotFoundError(err) {
					return err
				}
			}
			if reader != nil {
				if reader.isFinished() && p.minSeq < p.maxSeq {
					log.Infof("%s %d has finished, skip", reader.name, reader.seq)
					reader.close()
					reader = nil
				} else {
					log.Infof("init %s %d finished", reader.name, reader.seq)
					reader.dump()
					p.readers = append(p.readers, reader)
				}
			}
		}
		if p.minSeq < p.maxSeq {
			p.minSeq++
		} else {
			break
		}
	}
	return nil
}

func (p *BarrierQueueGroupReader) GetReaderNum() int {
	if p == nil {
		return 0
	}

	return len(p.readers)
}

func (p *BarrierQueueGroupReader) ioLoop() {
	log.Infof("start io loop")

	const (
		writeBufSize = 10240 * 4
	)
	writeBuf := make([]byte, writeBufSize)

	flushFinish := func(r *finishReq) int {
		p.SetLastFinishTs(r.seq, r.hash)

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
				p.SetLastFinishTs(r.seq, r.hash)
			default:
				r = nil
			}
		}
		var fn int
		for _, v := range all {
			var reader *BarrierQueueReader
			for _, x := range p.readers {
				if x.seq == v.seq {
					reader = x
					break
				}
			}
			if reader != nil {
				err := batchWriteFinishFile(reader.finishFile, reader.finishMap, v.list, writeBuf)
				if err != nil {
					log.Errorf("err:%v", err)
					p.sendErr(err)
				}
				reader.finishedCnt.Add(int64(len(v.list)))
				if reader.isFinished() {
					fn++
				}
			}
		}
		return fn
	}
	tk := time.NewTicker(time.Minute)
	for {
		select {
		case <-tk.C:
			// 兜底
			err := p.openReader()
			if err != nil {
				log.Errorf("err:%v", err)
				p.sendErr(err)
			}

		case seq := <-p.notifyWriteChan:
			if p.minSeq == 0 {
				p.minSeq = 1
			}
			for _, r := range p.readers {
				if r.seq == seq {
					r.recountIndex = true
					break
				}
			}
			if seq > p.maxSeq {
				p.maxSeq = seq
				err := p.openReader()
				if err != nil {
					log.Errorf("err:%v", err)
					p.sendErr(err)
				}
			}
			select {
			case p.notifyWriteChan4Consumer <- seq:
			default:
			}
		case <-p.exitChan:
			goto EXIT
		case r := <-p.finishChan:
			fn := flushFinish(r)
			if fn > 0 {
				err := p.openReader()
				if err != nil {
					log.Errorf("err:%v", err)
					p.sendErr(err)
				}
			}
		}
	}
EXIT:
	select {
	case p.closeChan <- 1:
	default:
	}
}

func (p *BarrierQueueGroupReader) Init() error {
	p.initMu.Lock()
	defer p.initMu.Unlock()
	if p.hasInit {
		p.warnChan <- &WarnMsg{
			Label: fmt.Sprintf("%s: duplicate init", p.name),
		}
		return errors.New("duplicate init")
	}
	err := scanInitFileGroup(p.name, p.dataPath, &p.minSeq, &p.maxSeq)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if p.minSeq > p.maxSeq {
		p.warnChan <- &WarnMsg{
			Label: fmt.Sprintf("%s: minSeq > maxSeq", p.name),
		}
		return errors.New("unreachable")
	}
	err = p.openReader()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	go p.ioLoop()
	p.hasInit = true
	return nil
}

func (p *BarrierQueueGroupReader) sendErr(err error) {
	select {
	case p.errChan <- err:
	default:
	}
}

func (p *BarrierQueueGroupReader) SetLastFinishTs(seq uint64, hash uint32) {
	for _, reader := range p.readers {
		if reader.seq == seq {
			if reader.delay.enableDelay {
				reader.delay.hash2LastFinishTsMu.Lock()

				if reader.delay.hash2LastFinishTs == nil {
					reader.delay.hash2LastFinishTs = map[uint32]uint64{}
				}
				reader.delay.hash2LastFinishTs[hash] = utils.NowMs()

				reader.delay.hash2LastFinishTsMu.Unlock()
			}
			break
		}
	}
}

func (p *BarrierQueueGroupReader) SetFinish(item *Item) {
	p.finishChan <- &finishReq{
		seq:   item.Seq,
		index: item.Index,
		hash:  item.Hash,
	}
}

func (p *BarrierQueueGroupReader) NotifyWrite(seq uint64) {
	select {
	case p.notifyWriteChan <- seq:
	default:
	}
}

func (p *BarrierQueueGroupReader) Close() {
	select {
	case p.exitChan <- 1:
	default:
	}
}

func (p *BarrierQueueGroupReader) popSkipHash(skipHash map[uint32]int, barrierCount int, checkCount bool) (*Item, error) {
	for _, reader := range p.readers {
		var wm *WarnMsg
		item, err := reader.popSkipHash(skipHash, barrierCount, checkCount, &wm)
		if wm != nil {
			wm.Label = fmt.Sprintf("%s.%d: %s", p.name, reader.seq, wm.Label)
			p.warnChan <- wm
		}
		if err != nil {
			log.Errorf("err:%v", err)
			p.sendErr(err)
			return nil, err
		}
		if item != nil {
			return item, nil
		}
	}
	return nil, nil
}

func (p *BarrierQueueGroupReader) PopSkipHash(skipHash map[uint32]int, barrierCount int) (*Item, error) {
	if barrierCount <= 0 {
		barrierCount = 1
	}

	item, err := p.popSkipHash(skipHash, barrierCount, true)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}

	if barrierCount > 1 {
		return p.popSkipHash(skipHash, barrierCount, false)
	}

	return nil, nil
}

func (p *BarrierQueueGroupReader) Retry(item *Item, delayMs uint32) {
	if delayMs > 3600*1000 {
		log.Warnf("%s retry delay too large %d", p.name, delayMs)

		// 拉群经常报这个，但又是正常的，先屏蔽掉
		//p.warnChan <- &WarnMsg{
		//	Label: "retry delay too large",
		//	Msg:   fmt.Sprintf("%s: retry delay %d ms", p.name, delayMs),
		//}
	}
	item.Data = nil
	for _, reader := range p.readers {
		if reader.seq == item.Seq {
			var wm *WarnMsg
			reader.retry(item, delayMs, &wm)
			if wm != nil {
				p.warnChan <- wm
			}
			return
		}
	}

	log.Warnf("not found reader %s seq %d", p.name, item.Seq)

	p.warnChan <- &WarnMsg{
		Label: fmt.Sprintf("%s: not found reader", p.name),
	}
}

func (p *BarrierQueueGroupReader) GetErrChan() chan error {
	return p.errChan
}

func (p *BarrierQueueGroupReader) GetWarnChan() chan *WarnMsg {
	return p.warnChan
}

func (p *BarrierQueueGroupReader) GetCloseChan() chan int {
	return p.closeChan
}

func (p *BarrierQueueGroupReader) GetQueueLen() int64 {
	var l int64
	for _, r := range p.readers {
		l += r.getQueueLen()
	}
	return l
}

func (p *BarrierQueueGroupReader) GetNotifyWriteChan4Consumer() chan uint64 {
	return p.notifyWriteChan4Consumer
}

func (p *BarrierQueueGroupReader) DumpReaders() {
	for _, v := range p.readers {
		log.Infof("dump reader with seq %d", v.seq)
		v.dump()
		log.Infof("===")
	}
}
