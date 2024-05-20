package index

import (
	"bitcask/data"
	"bytes"

	"github.com/google/btree"
)

// Indexer 通用索引接口,通过不同的数据结构实现该接口，就能使用不同的数据结构的索引
type Indexer interface {
	//put 向索引种存储key 对应的数据位置信息
	Put(key []byte, pos *data.LogRecordPos) bool

	//Get根据key取出对应的索引信息
	Get(key []byte) *data.LogRecordPos

	//Delate 根据key删除对应的索引位置信息
	Delate(key []byte) bool
}

// Item 索引结构体
type Item struct {
	key []byte
	pos *data.LogRecordPos
}

func (ai *Item) Less(bi btree.Item) bool {
	return bytes.Compare(ai.key, bi.(*Item).key) == -1
}
