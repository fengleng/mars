package mfile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gososy/sorpc/log"
	"io"
	"math"
	"os"
)

// 这里 copy 了 QueueReader 的代码过来改
// 主要用途就是修数据工具用的，遍历所有MQ数据，做一些事情
// 因为不想把这类功能，也耦合到正常的业务中，所以copy代码出来改了
type QueueScanner struct {
	fileBase
	cfg          *QueueReaderConfig
	itemCount    uint32
	readCursor   uint32
	pendingItems []*Item
	finishFile   *os.File
	finishMap    map[uint32]bool
	recountIndex bool
}

func NewQueueScanner(name, dataPath string, seq uint64, cfg *QueueReaderConfig) *QueueScanner {
	if cfg == nil {
		cfg = NewDefaultQueueReaderConfig()
	}
	if cfg.MaxCacheItemCount == 0 {
		cfg.MaxCacheItemCount = 8192
	}
	if cfg.MaxCacheDataBytes == 0 {
		cfg.MaxCacheDataBytes = 64 * 1024 * 1024
	}
	return &QueueScanner{
		fileBase: fileBase{
			name:     name,
			dataPath: dataPath,
			seq:      seq,
		},
		cfg:       cfg,
		finishMap: map[uint32]bool{},
	}
}
func (p *QueueScanner) statItemCount() error {
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
func (p *QueueScanner) Init() error {
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
	err = p.fillDataCache()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
func (p *QueueScanner) Close() {
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
func (p *QueueScanner) Pop() (*Item, error) {
	if len(p.pendingItems) == 0 {
		if p.recountIndex {
			err := p.statItemCount()
			if err != nil {
				log.Errorf("err:%v", err)
				return nil, err
			}
			p.recountIndex = false
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
			log.Warnf("item %+v load data fail", item)
			return nil, errors.New("load data fail")
		}
	}
	p.pendingItems = p.pendingItems[1:]
	return item, nil
}
func (p *QueueScanner) fillIndexCache() error {
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
				return nil
			}
			log.Errorf("err:%v", err)
			return err
		}
		log.Infof("%s seq %d: load index at %d with size %d, size %d",
			p.name, p.seq, p.readCursor*indexItemSize, size, n)
		n = n / indexItemSize
		b := binary.LittleEndian
		for i := 0; i < n; i++ {
			index := p.readCursor
			p.readCursor++
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
			//idx.reserved1 = b.Uint32(ptr[16:])
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
func (p *QueueScanner) fillDataCache() error {
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
	n, err := f.ReadAt(buf, int64(left))
	if err != nil && err != io.EOF {
		log.Errorf("err:%v", err)
		return err
	}
	if n < int(size) {
		buf = buf[:n]
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
func (p *QueueScanner) loadFinishFile() error {
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
