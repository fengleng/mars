package mfile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gososy/sorpc/log"
	"github.com/gososy/sorpc/utils"
	"go.uber.org/atomic"
)

//
// 提供给对象存储用的随机读
//
type RandReader struct {
	fileBase
	lastUsedAt uint32
}

//go:generate pie RandReaders.*
type RandReaders []*RandReader

func NewRandReader(name, dataPath string, seq uint64) *RandReader {
	return &RandReader{
		fileBase: fileBase{
			name:     name,
			dataPath: dataPath,
			seq:      seq,
		},
	}
}

var ErrFileNotExisted = errors.New("file not existed")

func (p *RandReader) Init() error {
	var err error
	idxPath := p.indexFilePath()
	datPath := p.dataFilePath()
	p.indexFile, err = os.OpenFile(idxPath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExisted
		}
		log.Errorf("err:%v", err)
		return err
	}
	p.dataFile, err = os.OpenFile(datPath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExisted
		}
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
func (p *RandReader) close() {
	if p.indexFile != nil {
		_ = p.indexFile.Close()
		p.indexFile = nil
	}
	if p.dataFile != nil {
		_ = p.dataFile.Close()
		p.dataFile = nil
	}
}

var ErrRecordNotFound = errors.New("record not found")

func (p *RandReader) Get(index uint32) (*Item, error) {
	return p.get(index)
}

func (p *RandReader) get(index uint32) (*Item, error) {
	i := p.indexFile
	d := p.dataFile
	if i == nil || d == nil {
		return nil, fmt.Errorf("file not opened")
	}
	var buf [indexItemSize]byte
	pos := indexItemSize * index
	n, err := i.ReadAt(buf[:], int64(pos))
	if err != nil {
		if err == io.EOF {
			return nil, ErrRecordNotFound
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	if n != indexItemSize {
		return nil, ErrRecordNotFound
	}
	b := binary.LittleEndian
	ptr := buf[:]
	var begMarker, endMarker uint16
	begMarker = b.Uint16(ptr)
	ptr = ptr[2:]
	endMarker = b.Uint16(ptr[28:])
	if begMarker != itemBegin {
		p.dataCorruption = true
		return nil, errors.New("invalid index item begin marker")
	}
	if endMarker != itemEnd {
		p.dataCorruption = true
		return nil, errors.New("invalid index item end marker")
	}
	var item Item
	item.CreatedAt = b.Uint32(ptr)
	item.CorpId = b.Uint32(ptr[4:])
	item.AppId = b.Uint32(ptr[8:])
	item.Hash = b.Uint32(ptr[12:])
	//item.reserved1 = b.Uint32(ptr[16:])
	item.offset = b.Uint32(ptr[20:])
	item.size = b.Uint32(ptr[24:])
	item.Index = index
	item.Seq = p.seq
	if item.size == 0 {
		return nil, errors.New("invalid index, buf size = 0")
	}
	dataBuf := make([]byte, item.size)
	var read uint32
	for read < item.size {
		n, err = d.ReadAt(dataBuf[read:], int64(item.offset))
		if err != nil {
			if err == io.EOF {
				return nil, errors.New("data file truncated")
			}
			log.Errorf("err:%v", err)
			return nil, err
		}
		read += uint32(n)
	}
	ptr = dataBuf[:]
	begMarker = b.Uint16(ptr)
	item.Data = ptr[2 : item.size-2]
	endMarker = b.Uint16(ptr[item.size-2:])
	if begMarker != itemBegin {
		p.dataCorruption = true
		return nil, errors.New("invalid data begin marker")
	}
	if endMarker != itemEnd {
		p.dataCorruption = true
		return nil, errors.New("invalid data end marker")
	}
	p.lastUsedAt = utils.Now()
	return &item, nil
}

type RandGroupReader struct {
	name                string
	dataPath            string
	readers             map[uint64]*RandReader
	readersMu           sync.RWMutex
	cleanRoutineStarted atomic.Bool
}

func NewRandGroupReader(name, dataPath string) *RandGroupReader {
	g := &RandGroupReader{
		name:     name,
		dataPath: dataPath,
		readers:  map[uint64]*RandReader{},
	}
	return g
}
func (p *RandGroupReader) getFast(seq uint64, index uint32) (*Item, error) {
	p.readersMu.RLock()
	defer p.readersMu.RUnlock()
	r := p.readers[seq]
	if r != nil {
		return r.get(index)
	}
	return nil, nil
}
func (p *RandGroupReader) tryOpen(seq uint64) (int, error) {
	p.readersMu.Lock()
	defer p.readersMu.Unlock()
	r := p.readers[seq]
	if r != nil {
		return len(p.readers), nil
	}
	r = NewRandReader(p.name, p.dataPath, seq)
	err := r.Init()
	if err != nil {
		if err != ErrFileNotExisted {
			log.Errorf("err:%v", err)
		}
		return len(p.readers), err
	}
	p.readers[seq] = r
	return len(p.readers), nil
}
func (p *RandGroupReader) cleanOpenedFiles() {
	p.readersMu.Lock()
	defer p.readersMu.Unlock()
	// sort by last used at
	list := RandReaders{}
	for _, v := range p.readers {
		list = append(list, v)
	}
	list = list.SortUsing(func(a, b *RandReader) bool {
		return a.lastUsedAt > b.lastUsedAt
	})
	for i := 50; i < len(list); i++ {
		r := list[i]
		r.close()
		delete(p.readers, r.seq)
	}
}
func (p *RandGroupReader) Get(seq uint64, index uint32) (*Item, error) {
	item, err := p.getFast(seq, index)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	cnt, err := p.tryOpen(seq)
	if err != nil {
		if err != ErrFileNotExisted {
			log.Errorf("err:%v", err)
		}
		return nil, err
	}
	if cnt > 100 {
		defer func() {
			if !p.cleanRoutineStarted.Load() {
				if p.cleanRoutineStarted.CAS(false, true) {
					go func() {
						p.cleanOpenedFiles()
						p.cleanRoutineStarted.Store(false)
					}()
				}
			}
		}()
	}
	// again
	item, err = p.getFast(seq, index)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	return nil, errors.New("get again fail")
}
