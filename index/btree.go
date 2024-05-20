package index

import (
	"bitcask/data"
	"sync"

	"github.com/google/btree"
)

// BTree 索引，主要封装了google的btree ku
type BTree struct {
	tree *btree.BTree
	lock *sync.RWMutex
}

func NewBTree() *BTree {
	return &BTree{
		tree: btree.New(32),
		lock: &sync.RWMutex{},
	}
}

// 实现索引接口
func (bt *BTree) Put(key []byte, pos *data.LogRecordPos) bool {
	it := &Item{key: key, pos: pos}
	bt.lock.Lock()
	bt.tree.ReplaceOrInsert(it)
	bt.lock.Unlock()
	return true
}

func (bt *BTree) Get(key []byte) *data.LogRecordPos {
	it := &Item{key: key}
	btreeitem := bt.tree.Get(it)
	if btreeitem == nil {
		return nil
	}
	return btreeitem.(*Item).pos
}

func (bt *BTree) Delate(key []byte) bool {
	it := &Item{key: key}
	bt.lock.Lock()
	//返回原来的item
	olditem := bt.tree.Delete(it)
	bt.lock.Unlock()
	return olditem != nil
}
