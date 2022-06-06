package skip_list

import (
	"bytes"
	"goSearcher/searcher/codec"
	"math/rand"
	"sync"
)

const (
	defaultMaxLevel = 48
)

type SkipList struct {
	header *Element

	rand *rand.Rand

	maxLevel int
	length   int
	lock     sync.RWMutex
	size     int64

	curmaxLevel int
}

func NewSkipList() *SkipList {

	header := &Element{
		levels: make([]*Element, defaultMaxLevel),
	}

	return &SkipList{
		header:      header,
		rand:        r,
		maxLevel:    defaultMaxLevel - 1,
		curmaxLevel: 0,
	}

}

// Element 跳表节点结构
type Element struct {
	levels []*Element
	entry  *codec.Entry
	score  float64
}

//创建一个新的跳表元素
func newElement(score float64, entry *codec.Entry, level int) *Element {

	return &Element{
		levels: make([]*Element, level),
		entry:  entry,
		score:  score,
	}
}

func (elem *Element) Entry() *codec.Entry {
	return elem.entry
}

// Add 往跳表加一个对象
func (list *SkipList) Add(data *codec.Entry) error {
	//implement me here!!!
	list.lock.Lock()
	defer list.lock.Unlock()

	//创建list elemt
	keyscore := list.calcScore(data.Key)
	key := data.Key
	level := list.randLevel()
	if level >= list.curmaxLevel {
		list.curmaxLevel = level
	}
	//因为是提升level层，所以，element一共level + 1层
	elem := newElement(keyscore, data, level+1)

	//查找要插入的位置
	prevElemt := list.header
	prevElemts := make([]*Element, list.maxLevel+1)

	for i := list.curmaxLevel; i >= 0; i-- {

		for ne := prevElemt.levels[i]; ne != nil; ne = prevElemt.levels[i] {
			if com := list.compare(keyscore, key, ne); com <= 0 {
				if com == 0 {
					//更新
					ne.entry = data
					return nil
				} else {
					prevElemt = ne
				}
			} else {
				// >0 找到插入的位置
				break
			}
		}
		//找到当前 level的 pre插入点
		prevElemts[i] = prevElemt
	}

	//插入
	for i := level; i >= 0; i-- {
		ne := prevElemts[i].levels[i]
		prevElemts[i].levels[i] = elem
		elem.levels[i] = ne

	}

	list.length++
	return nil
}

func (list *SkipList) Search(key []byte) (e *codec.Entry) {
	//implement me here!!!
	list.lock.RLock()
	defer list.lock.RUnlock()
	// 如果长度是0 则直接返回 nil
	if list.length == 0 {
		return nil
	}
	//从上到下查找
	prevElem := list.header
	keyScore := list.calcScore(key)

	for i := list.curmaxLevel; i >= 0; i-- {

		for ne := prevElem.levels[i]; ne != nil; ne = prevElem.levels[i] {
			if com := list.compare(keyScore, key, ne); com <= 0 {
				if com == 0 {
					return ne.entry
				} else {
					prevElem = ne
				}
			} else {
				break
			}

		}

	}
	return nil
}

func (list *SkipList) Close() error {
	return nil
}

func (list *SkipList) calcScore(key []byte) (score float64) {
	var hash uint64
	l := len(key)

	if l > 8 {
		l = 8
	}

	for i := 0; i < l; i++ {
		shift := uint(64 - 8 - i*8)
		hash |= uint64(key[i]) << shift
	}

	score = float64(hash)
	return
}

//compare 等于0，小于-1， 大于1
func (list *SkipList) compare(score float64, key []byte, next *Element) int {
	//implement me here!!!
	if score == next.score {
		return bytes.Compare(key, next.entry.Key)
	}

	if score < next.score {
		return -1
	} else {
		return 1
	}

}

//任意层数, 返回将此数据最多提升的层数
//提升的层数从1到48随机上升
func (list *SkipList) randLevel() int {
	//implement me here!!!
	level := 1

	for ; level < defaultMaxLevel; level++ {
		if list.rand.Intn(2) == 0 {
			return level
		}
	}

	return list.maxLevel
}

func (list *SkipList) Size() int64 {
	//implement me here!!!
	return list.size
}
