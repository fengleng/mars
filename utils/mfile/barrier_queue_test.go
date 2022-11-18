package mfile

import (
	"testing"

	"github.com/gososy/sorpc/log"
	"github.com/gososy/sorpc/utils"
)

func TestReaderTest(t *testing.T) {
	ReaderTest()
}

func TestIndexCacheItemList_appendItem(t *testing.T) {
	l := &IndexCacheItemList{}
	l.appendItem(&Item{
		CreatedAt:  utils.Now(),
		CorpId:     1,
		AppId:      2,
		Hash:       3,
		DelayType:  1,
		DelayValue: 100,
		Priority:   1,
	})
	l.appendItem(&Item{
		CreatedAt:  utils.Now(),
		CorpId:     2,
		AppId:      2,
		Hash:       3,
		DelayType:  1,
		DelayValue: 100,
		Priority:   0,
	})

	for {
		i := l.popItem(func(item *Item) bool {
			return true
		})
		log.Infof("%+v", i)
		if i == nil {
			break
		}
	}
}
