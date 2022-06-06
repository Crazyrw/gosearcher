package btree

// 查询统计（测试使用）
type stat struct {
	// 查询遍历到的节点数
	count int
	// 树的节点总数
	nodeCount int
	// 树的key/value总数
	kvCount int
	// 数的高度
	level int
}

func (b *stat) incrCount() {
	b.count++
}

func (b *stat) resetCount() {
	b.count = 0
}

// GetCount after Find/FindRange
func (b *stat) GetCount() int {
	return b.count
}

// GetLevel after tree.CountNode
func (b *stat) GetLevel() int {
	return b.level
}

// GetNodeCount after tree.CountNode
func (b *stat) GetNodeCount() int {
	return b.nodeCount
}

// GetKeyCount after tree.CountNode
func (b *stat) GetKeyCount() int {
	return b.kvCount
}
