package mfile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gososy/sorpc/log"
	"github.com/gososy/sorpc/utils"
	"go.uber.org/atomic"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	itemBegin = uint16(0x1234)
	itemEnd   = uint16(0x5678)
)

var wgPool sync.Pool

const indexItemSize = 32

//
// Data file format is:
// item begin marker | data | item end marker
//       2           |   V  |       2
//
//
// Index file format is
// item begin marker | index | item end marker
//       2           |   28  |       2
//
const (
	dataItemExtraSize = 4
)

func init() {
	wgPool.New = func() interface{} {
		return &sync.WaitGroup{}
	}
}
func wgPoolGet() *sync.WaitGroup {
	return wgPool.Get().(*sync.WaitGroup)
}
func wgPoolPut(wg *sync.WaitGroup) {
	wgPool.Put(wg)
}

const (
	DelayTypeNil    = uint32(0)
	DelayTypeRelate = uint32(1)
)

type Item struct {
	CreatedAt  uint32
	CorpId     uint32
	AppId      uint32
	Hash       uint32
	Data       []byte
	Seq        uint64
	Index      uint32
	offset     uint32
	size       uint32
	RetryCount uint32
	DelayType  uint32
	DelayValue uint32

	// 占2位, (0-3)
	Priority uint8
}
type writeReq struct {
	item *Item
	wg   *sync.WaitGroup
	// response
	err error
}
type fileBase struct {
	name           string
	dataPath       string
	seq            uint64
	indexFile      *os.File
	dataFile       *os.File
	dataCorruption bool
}

func packMisc(delayType, delayValue uint32, priority uint8) (res uint32) {
	if delayType == DelayTypeRelate {
		res = (1 << 31) | delayValue
	}

	if priority > 0 {
		res |= (uint32(priority) & 0x3) << 29
	}

	return
}

func unpackMisc(x uint32) (delayType, delayValue uint32, priority uint8) {
	if x&(1<<31) != 0 {
		delayType = DelayTypeRelate
		delayValue = x & ^uint32(0x7<<29)
	} else {
		delayValue = 0
		delayType = DelayTypeNil
	}

	priority = uint8((x >> 29) & 0x3)

	return
}

func (p *fileBase) indexFilePath() string {
	return path.Join(p.dataPath, fmt.Sprintf("%s.%d.idx", p.name, p.seq))
}

func (p *fileBase) dataFilePath() string {
	return path.Join(p.dataPath, fmt.Sprintf("%s.%d.dat", p.name, p.seq))
}

func (p *fileBase) finishFilePath() string {
	return path.Join(p.dataPath, fmt.Sprintf("%s.%d.finish", p.name, p.seq))
}

type fileWriter struct {
	fileBase
	writeBytes uint32
	writeItems uint32
	idxBuf     []byte
	datBuf     []byte
	needSync   bool
}

func (p *fileWriter) init() error {
	var err error
	idxPath := p.indexFilePath()
	datPath := p.dataFilePath()
	p.indexFile, err = os.OpenFile(idxPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	p.dataFile, err = os.OpenFile(datPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	idxInfo, err := os.Stat(idxPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	idxSize := int(idxInfo.Size())
	if idxSize%indexItemSize != 0 {
		p.dataCorruption = true
		return fmt.Errorf("invalid index file Size %d, index item Size %d", idxSize, indexItemSize)
	}
	p.writeItems = uint32(idxSize / indexItemSize)
	if idxSize == 0 {
		log.Infof("create new file %s", idxPath)
	}
	datInfo, err := os.Stat(datPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	p.writeBytes = uint32(datInfo.Size())
	return nil
}
func (p *fileWriter) batchWrite(reqList []*writeReq, startIndex *uint32) error {
	idxBufSize := len(reqList) * indexItemSize
	if len(p.idxBuf) < idxBufSize {
		p.idxBuf = make([]byte, idxBufSize)
	}
	var datBufSize int
	for _, v := range reqList {
		datBufSize += len(v.item.Data) + dataItemExtraSize
	}
	if len(p.datBuf) < datBufSize {
		p.datBuf = make([]byte, datBufSize)
	}
	{
		offset := p.writeBytes
		b := binary.LittleEndian
		for i, v := range reqList {
			p := p.idxBuf[i*indexItemSize:]
			b.PutUint16(p[0:2], itemBegin)
			p = p[2:]
			it := v.item
			datSize := len(it.Data) + dataItemExtraSize
			delay := packMisc(it.DelayType, it.DelayValue, it.Priority)
			b.PutUint32(p[0:4], utils.Now())       // created_at
			b.PutUint32(p[4:8], it.CorpId)         // corp_id
			b.PutUint32(p[8:12], it.AppId)         // app_id
			b.PutUint32(p[12:16], it.Hash)         // hash
			b.PutUint32(p[16:20], delay)           // delay
			b.PutUint32(p[20:24], offset)          // offset
			b.PutUint32(p[24:28], uint32(datSize)) // size
			offset += uint32(datSize)
			b.PutUint16(p[28:], itemEnd)
		}
		p := p.datBuf
		for _, v := range reqList {
			b.PutUint16(p[0:2], itemBegin)
			copy(p[2:], v.item.Data)
			b.PutUint16(p[2+len(v.item.Data):], itemEnd)
			p = p[len(v.item.Data)+dataItemExtraSize:]
		}
	}
	// write file
	{
		n, err := p.dataFile.Write(p.datBuf[:datBufSize])
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		for n < datBufSize {
			r, err := p.dataFile.Write(p.datBuf[n:datBufSize])
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			n += r
		}
	}
	{
		n, err := p.indexFile.Write(p.idxBuf[:idxBufSize])
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		for n < idxBufSize {
			r, err := p.indexFile.Write(p.idxBuf[n:idxBufSize])
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			n += r
		}
	}
	*startIndex = p.writeItems
	p.writeBytes += uint32(datBufSize)
	p.writeItems += uint32(len(reqList))
	p.needSync = true
	return nil
}
func (p *fileWriter) sync() error {
	if !p.needSync {
		return nil
	}
	p.needSync = false
	if p.indexFile != nil {
		err := p.indexFile.Sync()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	if p.dataFile != nil {
		err := p.dataFile.Sync()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	return nil
}
func (p *fileWriter) check() error {
	if p.indexFile != nil {
		st, err := p.indexFile.Stat()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		curSize := uint32(st.Size())
		if curSize%indexItemSize != 0 {
			p.dataCorruption = true
			return fmt.Errorf("index file size %d not align %d bytes of index item",
				curSize, indexItemSize)
		}
		expSize := p.writeItems * indexItemSize
		if expSize != curSize {
			if curSize < expSize {
				p.dataCorruption = true
				return fmt.Errorf("index file truncated, exp size %d, but cur size %d", expSize, curSize)
			} else {
				// reset index
				p.writeItems = curSize / indexItemSize
			}
		}
	}
	if p.dataFile != nil {
		st, err := p.dataFile.Stat()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		curSize := uint32(st.Size())
		if curSize != p.writeBytes {
			if curSize < p.writeBytes {
				p.dataCorruption = true
				return fmt.Errorf("data file truncated, exp size %d,  but cur size %d", p.writeBytes, curSize)
			}
			// reset write bytes
			p.writeBytes = curSize
		}
	}
	return nil
}
func (p *fileWriter) close() {
	if p.indexFile != nil {
		if p.needSync {
			_ = p.indexFile.Sync()
		}
		_ = p.indexFile.Close()
		p.indexFile = nil
	}
	if p.dataFile != nil {
		if p.needSync {
			_ = p.dataFile.Sync()
		}
		_ = p.dataFile.Close()
		p.dataFile = nil
	}
	p.needSync = false
}

type FileGroupConfig struct {
	MaxBytesPerFile uint32
	MaxItemPerFile  uint32
}

func NewDefaultFileGroupConfig() *FileGroupConfig {
	return &FileGroupConfig{
		MaxBytesPerFile: 256 * 1024 * 1024,
		MaxItemPerFile:  (256 * 1024 * 1024) / indexItemSize,
	}
}

type FileGroup struct {
	name       string
	dataPath   string
	cfg        *FileGroupConfig
	curSeq     uint64
	writer     *fileWriter
	hasInit    bool
	exitChan   chan int
	exitedChan chan int
	exiting    atomic.Uint32
	writeChan  chan *writeReq
}

func NewFileGroup(name, dataPath string, cfg *FileGroupConfig) *FileGroup {
	if name == "" {
		panic("name can not be empty")
	}
	for i := 0; i < len(name); i++ {
		if name[i] == '.' {
			panic("name has invalid char")
		}
	}
	if dataPath == "" {
		panic("data path can not be empty")
	}
	if cfg == nil {
		cfg = NewDefaultFileGroupConfig()
	}
	return &FileGroup{
		name:       name,
		dataPath:   dataPath,
		cfg:        cfg,
		exitChan:   make(chan int),
		exitedChan: make(chan int),
		writeChan:  make(chan *writeReq, 10000),
	}
}
func (p *FileGroup) batchWrite(reqList []*writeReq) {
	if len(reqList) == 0 {
		return
	}

	failAll := func(err error) {
		for _, v := range reqList {
			v.err = err
			v.wg.Done()
		}
	}
	if p.writer != nil && p.writer.dataCorruption {
		failAll(errors.New("data corruption"))
		return
	}
	if p.exiting.Load() > 0 {
		failAll(errors.New("exiting"))
		return
	}
	if p.writer == nil {
		p.writer = &fileWriter{
			fileBase: fileBase{
				name:     p.name,
				dataPath: p.dataPath,
				seq:      p.curSeq,
			},
		}
		err := p.writer.init()
		if err != nil {
			log.Errorf("err:%v", err)
			failAll(err)
			return
		}
	}
	// 算下是否要切换文件分片
	var bs uint32
	var switchToNext bool
	for i := 0; i < len(reqList); i++ {
		bs += uint32(len(reqList[i].item.Data))
		if p.writer.writeBytes+bs > p.cfg.MaxBytesPerFile {
			switchToNext = true
			break
		}
		if p.writer.writeItems+uint32(i) > p.cfg.MaxItemPerFile {
			switchToNext = true
			break
		}
	}
	var startIndex uint32
	err := p.writer.batchWrite(reqList, &startIndex)
	if err != nil {
		failAll(err)
		return
	}
	for j, x := range reqList {
		x.item.Seq = p.writer.seq
		x.item.Index = startIndex + uint32(j)
		x.wg.Done()
	}

	if switchToNext {
		p.writer.close()
		p.writer = nil
		p.curSeq++
	}
}
func (p *FileGroup) ioLoop() {
	log.Infof("start io loop")
	syncTicker := time.NewTicker(5 * time.Second)
	checkTicker := time.NewTicker(30 * time.Second)
	defer func() {
		syncTicker.Stop()
		checkTicker.Stop()
	}()
	doCheck := func() {
		if p.writer != nil {
			err := p.writer.check()
			if err != nil {
				log.Errorf("check err:%v", err)
			}
		}
	}
	// check file before loop
	doCheck()
	for {
		select {
		case <-p.exitChan:
			goto exit
		case <-syncTicker.C:
			if p.writer != nil {
				_ = p.writer.sync()
			}
		case <-checkTicker.C:
			doCheck()
		case req := <-p.writeChan:
			// 聚合一下
			var reqList = []*writeReq{
				req,
			}
			bs := len(req.item.Data)
			for i := 0; i < 100 && bs > 64*1024*1024; i++ {
				select {
				case x := <-p.writeChan:
					reqList = append(reqList, x)
					bs += len(x.item.Data)
				default:
					goto mergeExit
				}
			}
		mergeExit:
			p.batchWrite(reqList)
		}
	}
exit:
	if p.writer != nil {
		p.writer.close()
	}
	select {
	case p.exitedChan <- 1:
	default:
	}
}
func scanInitFileGroup(name, dataPath string, minSeqPtr, maxSeqPtr *uint64) error {
	st, err := os.Stat(dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dataPath, 0777)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		} else {
			log.Errorf("stat %s err:%v", dataPath, err)
			return err
		}
	} else {
		if !st.IsDir() {
			return fmt.Errorf("%s existed, but not dir", dataPath)
		}
		// scan
		fileList, err := ioutil.ReadDir(dataPath)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		pfx := name + "."
		maxSeq := uint64(0)
		minSeq := uint64(0)
		for _, f := range fileList {
			if f.IsDir() {
				continue
			}
			name := f.Name()
			if !strings.HasPrefix(name, pfx) {
				continue
			}
			if !strings.HasSuffix(name, ".idx") {
				continue
			}
			name = name[len(pfx):]
			pos := strings.IndexByte(name, '.')
			if pos <= 0 {
				continue
			}
			seqStr := name[:pos]
			seq, err := strconv.ParseUint(seqStr, 10, 64)
			if err == nil {
				if seq > maxSeq {
					maxSeq = seq
				}
				if minSeq == 0 || seq < minSeq {
					minSeq = seq
				}
			}
		}
		if maxSeqPtr != nil {
			*maxSeqPtr = maxSeq
		}
		if minSeqPtr != nil {
			*minSeqPtr = minSeq
		}
	}
	return nil
}
func (p *FileGroup) Init() error {
	if p.hasInit {
		panic("duplicate init")
	}
	err := scanInitFileGroup(p.name, p.dataPath, nil, &p.curSeq)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if p.curSeq == 0 {
		p.curSeq = 1
	}
	go func() {
		p.ioLoop()
	}()
	p.hasInit = true
	return nil
}
func (p *FileGroup) Write(item *Item) error {
	if p.exiting.Load() > 0 {
		return errors.New("exiting")
	}
	req := &writeReq{
		item: item,
		wg:   wgPoolGet(),
	}
	req.wg.Add(1)
	p.writeChan <- req
	req.wg.Wait()
	if req.err != nil {
		log.Errorf("err:%v", req.err)
		return req.err
	}
	wgPoolPut(req.wg)
	return nil
}
func (p *FileGroup) Close() {
	if !p.hasInit {
		return
	}
	if !p.exiting.CAS(0, 1) {
		// has closed by other routine
		return
	}
	select {
	case p.exitChan <- 1:
	default:
	}
	// block until real Close, avoid concurrency open & op file by other routine
	<-p.exitedChan
}
