package btree

//TODO: 落盘文件过大 一个磁盘块不能存储的下一个节点
import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"goSearcher/errs"
	"os"
	"sort"
	"sync"
	"syscall"
)

const (
	INITADDR = 0xdeadbeef
	//windows 磁盘块大小是4k 内存页大小是64k
	//BLOCKSIZE      = 4096
	MAX_FREEBLOCKS = 100

	ORDER = 4 //树的阶数
)

var err error

type BTreeInfo struct {
	root       uint64
	nodePool   *sync.Pool
	freeBlocks []uint64
	file       *os.File
	blockSize  uint64
	fileSize   uint64
}
type Node struct {
	IsActive bool //节点所在的磁盘空间是否在当前b+树内 (文件中的数据节点 是否已经存放到b+树中)
	Children []uint64

	Self   uint64
	Prev   uint64
	Next   uint64
	Parent uint64

	Keys    []string
	Records []string
	IsLeaf  bool
}

func NewTree(fileName string) (*BTreeInfo, error) {
	var (
		stat  syscall.Statfs_t
		fstat os.FileInfo
	)
	tree := &BTreeInfo{}

	tree.root = INITADDR
	//定义池中可以放什么类型的数据
	tree.nodePool = &sync.Pool{
		New: func() interface{} {
			return &Node{}
		},
	}
	tree.freeBlocks = make([]uint64, 0, MAX_FREEBLOCKS)

	tree.file, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	if err = syscall.Statfs(fileName, &stat); err != nil {
		return nil, err
	}

	tree.blockSize = uint64(stat.Bsize)
	if tree.blockSize == 0 {
		return nil, errors.New("blockSize should be nonzero")
	}

	if fstat, err = tree.file.Stat(); err != nil {
		return nil, err
	}
	tree.fileSize = uint64(fstat.Size())
	//如果文件中已经存储了部分数据 需要对该数据创建b+树
	if tree.fileSize != 0 {
		//将文件中的数据不断地插入树中 重建树
		if err = tree.restructRootNode(); err != nil {
			return nil, err
		}
		if err = tree.checkDiskBlockForFreeNodeList(); err != nil {
			return nil, err
		}
	}
	fmt.Println("输出创建树的详细信息")
	fmt.Println("root: ", tree.root)
	fmt.Printf("freeBlocks: %v\n", tree.freeBlocks)
	fmt.Println("blockSize: ", tree.blockSize)
	fmt.Println("fileSize: ", tree.fileSize)

	return tree, nil
}
func (tree *BTreeInfo) Insert(key string, value string) error {
	var node *Node
	if tree.root == INITADDR {
		if node, err = tree.newNodeFromDisk(); err != nil {
			return err
		}
		tree.root = node.Self
		node.IsActive = true
		node.Keys = append(node.Keys, key)
		node.Records = append(node.Records, value)
		node.IsLeaf = true
		return tree.flushNodeAndPutNodePool(node)
	}
	return tree.insertIntoLeaf(key, value)
}
func (tree *BTreeInfo) Close() error {
	if tree.file != nil {
		tree.file.Sync()
		return tree.file.Close()
	}
	return nil
}
func (tree *BTreeInfo) Find(key string) (string, error) {
	var node *Node
	//没有节点
	if tree.root == INITADDR {
		return "", err
	}
	if node, err = tree.newMappingNodeFromPool(INITADDR); err != nil {
		return "", err
	}
	if err = tree.findLeaf(node, key); err != nil {
		return "", err
	}
	defer tree.putNodePool(node)

	for i, nkey := range node.Keys {
		if nkey == key {
			return node.Records[i], nil
		}
	}
	return "", errs.NotFoundKey
}

//------------------------------私有方法--------------------------------------------
//-----------------------NewTree start--------------------------------
func (tree *BTreeInfo) restructRootNode() error {
	node := &Node{} //不是nil 是有地址的
	//读取文件中所有地节点 一个文件在一个磁盘块中 一个文件中存放了一整棵树
	for off := uint64(0); off < tree.fileSize; off += tree.blockSize {
		if err := tree.seekNode(node, off); err != nil {
			return err
		}
		if node.IsActive {
			break
		}
	}
	//节点没办法创建 存储在磁盘块
	if !node.IsActive {
		return errs.InvalidDBFormat
	}
	//在b+树插入该节点
	if node.Parent != INITADDR {
		if err = tree.seekNode(node, node.Parent); err != nil {
			return err
		}
	}
	tree.root = node.Self
	return nil
}

//似懂非懂
func (tree *BTreeInfo) checkDiskBlockForFreeNodeList() error {
	node := &Node{}
	bs := tree.blockSize
	for off := uint64(0); off < tree.fileSize && len(tree.freeBlocks) < MAX_FREEBLOCKS; off += bs {
		if off+bs > tree.fileSize {
			break
		}

		if err = tree.seekNode(node, off); err != nil {
			return err
		}

		if !node.IsActive {
			tree.freeBlocks = append(tree.freeBlocks, off)
		}
	}
	//给当前地文件扩容？？
	//存储得话 应该是一整棵树在一个文件 但是这个文件包含很多磁盘块 一个磁盘块存放多个节点
	next_file := ((tree.fileSize + 4095) / 4096) * 4096
	for len(tree.freeBlocks) < MAX_FREEBLOCKS {
		tree.freeBlocks = append(tree.freeBlocks, next_file)
		next_file += bs
	}
	tree.fileSize = next_file
	return nil
}
func (tree *BTreeInfo) seekNode(node *Node, off uint64) error {
	if node == nil {
		return fmt.Errorf("seekNode(): node is nil")
	}

	tree.clearNodeForUsage(node)

	buf := make([]byte, 8) //地址使用的uint64 8个字节
	if n, err := tree.file.ReadAt(buf, int64(off)); err != nil {
		return err
	} else if uint64(n) != 8 {
		return fmt.Errorf("readat %d from %s, expected len = %d but get %d", off, tree.file.Name(), 4, n)
	}

	//1.根据磁盘块地址读取数据
	bs := bytes.NewBuffer(buf)
	dataLen := uint64(0)
	//二进制读取数据
	if err = binary.Read(bs, binary.LittleEndian, &dataLen); err != nil {
		return err
	}
	if uint64(dataLen)+8 > tree.blockSize {
		return fmt.Errorf("flushNode len(node) = %d exceed t.blockSize %d", uint64(dataLen)+4, tree.blockSize)
	}
	buf = make([]byte, dataLen)
	if n, err := tree.file.ReadAt(buf, int64(off)+8); err != nil {
		return err
	} else if uint64(n) != uint64(dataLen) {
		return fmt.Errorf("readat %d from %s, expected len = %d but get %d", int64(off)+4, tree.file.Name(), dataLen, n)
	}

	//2.bs中为读取到的数据 对数据进行判断
	bs = bytes.NewBuffer(buf)

	// IsActive 赋值
	if err = binary.Read(bs, binary.LittleEndian, &node.IsActive); err != nil {
		return err
	}
	// Children 赋值
	childCount := uint8(0)
	if err = binary.Read(bs, binary.LittleEndian, &childCount); err != nil {
		return err
	}
	node.Children = make([]uint64, childCount)
	for i := uint8(0); i < childCount; i++ {
		child := uint64(0)
		if err = binary.Read(bs, binary.LittleEndian, &child); err != nil {
			return err
		}
		node.Children[i] = child
	}
	// Self 赋值
	self := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &self); err != nil {
		return err
	}
	node.Self = self
	// Next 赋值
	next := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &next); err != nil {
		return err
	}
	node.Next = next
	// Prev 赋值
	prev := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &node.Prev); err != nil {
		return err
	}
	node.Prev = prev
	// Parent 赋值
	parent := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &parent); err != nil {
		return err
	}
	node.Parent = parent
	// Keys
	keysCount := uint8(0)
	if err = binary.Read(bs, binary.LittleEndian, &keysCount); err != nil {
		return err
	}
	node.Keys = make([]string, keysCount)
	for i := uint8(0); i < keysCount; i++ {
		//if err = binary.Read(bs, binary.LittleEndian, &node.Keys[i]); err != nil {
		//	return err
		//}
		l := uint8(0)
		if err = binary.Read(bs, binary.LittleEndian, &l); err != nil {
			return err
		}
		v := make([]byte, l)
		if err = binary.Read(bs, binary.LittleEndian, &v); err != nil {
			return err
		}
		node.Keys[i] = string(v)
	}
	// Records
	recordCount := uint8(0)
	if err = binary.Read(bs, binary.LittleEndian, &recordCount); err != nil {
		return err
	}
	node.Records = make([]string, recordCount)
	for i := uint8(0); i < recordCount; i++ {
		l := uint8(0)
		if err = binary.Read(bs, binary.LittleEndian, &l); err != nil {
			return err
		}
		v := make([]byte, l)
		if err = binary.Read(bs, binary.LittleEndian, &v); err != nil {
			return err
		}
		node.Records[i] = string(v)
	}
	// IsLeaf
	if err = binary.Read(bs, binary.LittleEndian, &node.IsLeaf); err != nil {
		return err
	}

	return nil
}

//guess: 应该是从池中拿出一个node  清空 然后复用
func (tree *BTreeInfo) clearNodeForUsage(node *Node) {
	node.IsActive = false
	node.Children = nil
	node.Self = INITADDR
	node.Prev = INITADDR
	node.Next = INITADDR
	node.Parent = INITADDR
	node.Keys = nil
	node.Records = nil
	node.IsLeaf = false
}

//-----------------------NewTree end--------------------------------

//-----------------------Insert start--------------------------------
func (tree *BTreeInfo) newNodeFromDisk() (*Node, error) {
	var node *Node
	node = tree.nodePool.Get().(*Node)
	if len(tree.freeBlocks) > 0 {
		off := tree.freeBlocks[0]
		tree.freeBlocks = tree.freeBlocks[1:len(tree.freeBlocks)]
		tree.initNodeForUsage(node)
		node.Self = off
		return node, nil
	}
	if err = tree.checkDiskBlockForFreeNodeList(); err != nil {
		return nil, err
	}
	if len(tree.freeBlocks) > 0 {
		off := tree.freeBlocks[0]
		tree.freeBlocks = tree.freeBlocks[1:len(tree.freeBlocks)]
		tree.initNodeForUsage(node)
		node.Self = off
		return node, nil
	}
	return nil, fmt.Errorf("can't not malloc more node")
}

//为新节点 初始化
func (tree *BTreeInfo) initNodeForUsage(node *Node) {
	node.IsActive = true
	node.Children = nil
	node.Self = INITADDR
	node.Next = INITADDR
	node.Prev = INITADDR
	node.Parent = INITADDR
	node.Keys = nil
	node.Records = nil
	node.IsLeaf = false
}

//新建节点之后 将节点地数据写入磁盘  在根节点上
func (tree *BTreeInfo) flushNodeAndPutNodePool(node *Node) error {
	if err := tree.flushNode(node); err != nil {
		return err
	}
	tree.putNodePool(node)
	return nil
}
func (tree *BTreeInfo) flushNode(node *Node) error {
	if node == nil {
		return fmt.Errorf("flushNode == nil")
	}
	if tree.file == nil {
		return fmt.Errorf("flush node into disk, but not open file")
	}

	var length int
	bs := bytes.NewBuffer(make([]byte, 0))
	// IsActive
	if err = binary.Write(bs, binary.LittleEndian, node.IsActive); err != nil {
		return nil
	}
	// Children
	childCount := uint8(len(node.Children))
	if err = binary.Write(bs, binary.LittleEndian, childCount); err != nil {
		return err
	}
	for _, v := range node.Children {
		if err = binary.Write(bs, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	// Self
	if err = binary.Write(bs, binary.LittleEndian, node.Self); err != nil {
		return err
	}

	// Next
	if err = binary.Write(bs, binary.LittleEndian, node.Next); err != nil {
		return err
	}

	// Prev
	if err = binary.Write(bs, binary.LittleEndian, node.Prev); err != nil {
		return err
	}

	// Parent
	if err = binary.Write(bs, binary.LittleEndian, node.Parent); err != nil {
		return err
	}
	// Keys
	keysCount := uint8(len(node.Keys))
	if err = binary.Write(bs, binary.LittleEndian, keysCount); err != nil {
		return err
	}
	for _, v := range node.Keys {
		if err = binary.Write(bs, binary.LittleEndian, uint8(len([]byte(v)))); err != nil {
			return err
		}
		if err = binary.Write(bs, binary.LittleEndian, []byte(v)); err != nil {
			return err
		}
	}
	// Record
	recordCount := uint8(len(node.Records))
	if err = binary.Write(bs, binary.LittleEndian, recordCount); err != nil {
		return err
	}
	for _, v := range node.Records {
		if err = binary.Write(bs, binary.LittleEndian, uint8(len([]byte(v)))); err != nil {
			return err
		}
		if err = binary.Write(bs, binary.LittleEndian, []byte(v)); err != nil {
			return err
		}
	}

	// IsLeaf
	if err = binary.Write(bs, binary.LittleEndian, node.IsLeaf); err != nil {
		return err
	}

	dataLen := len(bs.Bytes())
	if uint64(dataLen)+8 > tree.blockSize {
		return fmt.Errorf("flushNode len(node) = %d exceed t.blockSize %d", uint64(dataLen)+4, tree.blockSize)
	}
	tmpbs := bytes.NewBuffer(make([]byte, 0))
	if err = binary.Write(tmpbs, binary.LittleEndian, uint64(dataLen)); err != nil {
		return err
	}

	data := append(tmpbs.Bytes(), bs.Bytes()...)
	if length, err = tree.file.WriteAt(data, int64(node.Self)); err != nil {
		return err
	} else if len(data) != length {
		return fmt.Errorf("writeat %d into %s, expected len = %d but get %d", int64(node.Self), tree.file.Name(), len(data), length)
	}
	return nil
}
func (tree *BTreeInfo) putNodePool(node *Node) {
	tree.nodePool.Put(node)
}

//新增的节点为叶子节点
func (tree *BTreeInfo) insertIntoLeaf(key string, value string) error {
	var (
		leaf    *Node
		index   int
		newLeaf *Node
	)
	if leaf, err = tree.newMappingNodeFromPool(INITADDR); err != nil {
		return err
	}
	if err = tree.findLeaf(leaf, key); err != nil {
		return err
	}
	if index, err = insertKeyValIntoLeaf(leaf, key, value); err != nil {
		return err
	}

	// update the last key of parent's if necessary
	if err = tree.mayUpdatedLastParentKey(leaf, index); err != nil {
		return err
	}

	// insert key/val into leaf
	if len(leaf.Keys) <= ORDER {
		return tree.flushNodeAndPutNodePool(leaf)
	}

	// split leaf so new leaf node
	if newLeaf, err = tree.newNodeFromDisk(); err != nil {
		return err
	}
	newLeaf.IsLeaf = true
	if err = tree.splitLeafIntoTowleaves(leaf, newLeaf); err != nil {
		return err
	}

	if err = tree.flushNodesAndPutNodesPool(newLeaf, leaf); err != nil {
		return err
	}

	// insert split key into parent
	return tree.insertIntoParent(leaf.Parent, leaf.Self, leaf.Keys[len(leaf.Keys)-1], newLeaf.Self)
}
func (tree *BTreeInfo) newMappingNodeFromPool(off uint64) (*Node, error) {
	node := tree.nodePool.Get().(*Node)
	tree.initNodeForUsage(node)
	if off == INITADDR {
		return node, nil
	}
	tree.clearNodeForUsage(node)
	if err := tree.seekNode(node, off); err != nil {
		return nil, err
	}
	return node, nil
}
func (tree *BTreeInfo) findLeaf(node *Node, key string) error {
	var root *Node
	rootaddr := tree.root
	//根节点没分裂
	if rootaddr == INITADDR {
		return nil
	}
	if root, err = tree.newMappingNodeFromPool(rootaddr); err != nil {
		return err
	}
	defer tree.putNodePool(root)
	*node = *root
	for !node.IsLeaf {
		idx := sort.Search(len(node.Keys), func(i int) bool {
			return key <= node.Keys[i]
		})
		if idx == len(node.Keys) {
			idx = len(node.Keys) - 1
		}
		if err = tree.seekNode(node, node.Children[idx]); err != nil {
			return err
		}
	}
	return nil
}
func insertKeyValIntoLeaf(node *Node, key string, value string) (int, error) {
	index := sort.Search(len(node.Keys), func(i int) bool {
		return key <= node.Keys[i]
	})
	if index < len(node.Keys) && node.Keys[index] == key {
		fmt.Errorf(key)
		return 0, errs.HasExistedKeyError
	}

	node.Keys = append(node.Keys, key)
	node.Records = append(node.Records, value)
	for i := len(node.Keys) - 1; i > index; i-- {
		node.Keys[i] = node.Keys[i-1]
		node.Records[i] = node.Records[i-1]
	}
	node.Keys[index] = key
	node.Records[index] = value
	return index, nil
}
func (tree *BTreeInfo) mayUpdatedLastParentKey(leaf *Node, index int) error {
	// update the last key of parent's if necessary
	if index == len(leaf.Keys)-1 && leaf.Parent != INITADDR {
		key := leaf.Keys[len(leaf.Keys)-1]
		updateNodeOff := leaf.Parent
		var (
			updateNode *Node
			node       *Node
		)

		if node, err = tree.newMappingNodeFromPool(leaf.Self); err != nil {
			return err
		}
		*node = *leaf
		defer tree.putNodePool(node)

		for updateNodeOff != INITADDR && index == len(node.Keys)-1 {
			if updateNode, err = tree.newMappingNodeFromPool(updateNodeOff); err != nil {
				return err
			}
			for i, v := range updateNode.Children {
				if v == node.Self {
					index = i
					break
				}
			}
			updateNode.Keys[index] = key
			if err = tree.flushNodeAndPutNodePool(updateNode); err != nil {
				return err
			}
			updateNodeOff = updateNode.Parent
			*node = *updateNode
		}
	}
	return nil
}
func (tree *BTreeInfo) splitLeafIntoTowleaves(leaf *Node, new_leaf *Node) error {
	var (
		i, split int
	)
	split = cut(ORDER)

	for i = split; i <= ORDER; i++ {
		new_leaf.Keys = append(new_leaf.Keys, leaf.Keys[i])
		new_leaf.Records = append(new_leaf.Records, leaf.Records[i])
	}

	// adjust relation
	leaf.Keys = leaf.Keys[:split]
	leaf.Records = leaf.Records[:split]

	new_leaf.Next = leaf.Next
	leaf.Next = new_leaf.Self
	new_leaf.Prev = leaf.Self

	new_leaf.Parent = leaf.Parent

	if new_leaf.Next != INITADDR {
		var (
			nextNode *Node
			err      error
		)
		if nextNode, err = tree.newMappingNodeFromPool(new_leaf.Next); err != nil {
			return err
		}
		nextNode.Prev = new_leaf.Self
		if err = tree.flushNodesAndPutNodesPool(nextNode); err != nil {
			return err
		}
	}

	return err
}
func cut(length int) int {
	return (length + 1) / 2
}
func (tree *BTreeInfo) flushNodesAndPutNodesPool(nodes ...*Node) error {
	for _, n := range nodes {
		if err := tree.flushNodeAndPutNodePool(n); err != nil {
			return err
		}
	}
	return err
}
func (tree *BTreeInfo) insertIntoParent(parent_off uint64, left_off uint64, key string, right_off uint64) error {
	var (
		index  int
		parent *Node
		left   *Node
		right  *Node
	)
	if parent_off == uint64(INITADDR) {
		if left, err = tree.newMappingNodeFromPool(left_off); err != nil {
			return err
		}
		if right, err = tree.newMappingNodeFromPool(right_off); err != nil {
			return err
		}
		if err = tree.newRootNode(left, right); err != nil {
			return err
		}
		return tree.flushNodesAndPutNodesPool(left, right)
	}

	if parent, err = tree.newMappingNodeFromPool(parent_off); err != nil {
		return err
	}

	index = getIndex(parent.Keys, key)
	insertIntoNode(parent, index, left_off, key, right_off)

	if len(parent.Keys) <= ORDER {
		return tree.flushNodesAndPutNodesPool(parent)
	}

	return tree.insertIntoNodeAfterSplitting(parent)
}
func (tree *BTreeInfo) newRootNode(left *Node, right *Node) error {
	var (
		root *Node
		err  error
	)

	if root, err = tree.newNodeFromDisk(); err != nil {
		return err
	}
	root.Keys = append(root.Keys, left.Keys[len(left.Keys)-1])
	root.Keys = append(root.Keys, right.Keys[len(right.Keys)-1])
	root.Children = append(root.Children, left.Self)
	root.Children = append(root.Children, right.Self)
	left.Parent = root.Self
	right.Parent = root.Self

	tree.root = root.Self
	return tree.flushNodeAndPutNodePool(root)
}
func getIndex(keys []string, key string) int {
	idx := sort.Search(len(keys), func(i int) bool {
		return key <= keys[i]
	})
	return idx
}
func insertIntoNode(parent *Node, index int, left_off uint64, key string, right_off uint64) {
	var (
		i int
	)
	parent.Keys = append(parent.Keys, key)
	for i = len(parent.Keys) - 1; i > index; i-- {
		parent.Keys[i] = parent.Keys[i-1]
	}
	parent.Keys[index] = key

	if index == len(parent.Children) {
		parent.Children = append(parent.Children, right_off)
		return
	}
	tmpChildren := append([]uint64{}, parent.Children[index+1:]...)
	parent.Children = append(append(parent.Children[:index+1], right_off), tmpChildren...)
}
func (tree *BTreeInfo) insertIntoNodeAfterSplitting(old_node *Node) error {
	var (
		newNode, child, nextNode *Node
		err                      error
		i, split                 int
	)

	if newNode, err = tree.newNodeFromDisk(); err != nil {
		return err
	}

	split = cut(ORDER)

	for i = split; i <= ORDER; i++ {
		newNode.Children = append(newNode.Children, old_node.Children[i])
		newNode.Keys = append(newNode.Keys, old_node.Keys[i])

		// update new_node children relation
		if child, err = tree.newMappingNodeFromPool(old_node.Children[i]); err != nil {
			return err
		}
		child.Parent = newNode.Self
		if err = tree.flushNodesAndPutNodesPool(child); err != nil {
			return err
		}
	}
	newNode.Parent = old_node.Parent

	old_node.Children = old_node.Children[:split]
	old_node.Keys = old_node.Keys[:split]

	newNode.Next = old_node.Next
	old_node.Next = newNode.Self
	newNode.Prev = old_node.Self

	if newNode.Next != INITADDR {
		if nextNode, err = tree.newMappingNodeFromPool(newNode.Next); err != nil {
			return err
		}
		nextNode.Prev = newNode.Self
		if err = tree.flushNodesAndPutNodesPool(nextNode); err != nil {
			return err
		}
	}

	if err = tree.flushNodesAndPutNodesPool(old_node, newNode); err != nil {
		return err
	}

	return tree.insertIntoParent(old_node.Parent, old_node.Self, old_node.Keys[len(old_node.Keys)-1], newNode.Self)
}

//-----------------------Insert end--------------------------------
