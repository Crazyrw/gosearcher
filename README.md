### 索引创建
    1.内存版 B+树测试
	=== RUN   TestCreateMemoryBtree
	12193026 12193026
	2022/06/11 10:04:09 create memory btree cost  2m45.391381925s
	2022/06/11 10:04:09 query 二胎( 1210 ) cost  10.941µs
	--- PASS: TestCreateMemoryBtree (165.39s)
	PASS
	ok      goSearcher/searcher/core        165.670s

    2.内存版 SkipList测试
	=== RUN   TestCreateSkipList
	12193026
	2022/06/11 10:09:25 create skiplist cost  1m47.692093479s
	2022/06/11 10:09:25 query 二胎( 10891 ) cost  4.74µs
	--- PASS: TestCreateSkipList (107.69s)
	PASS
	ok      goSearcher/searcher/core        107.937s

    落盘B+树存在问题 一个磁盘块不足以放下一个较大的node >4K
	- 考虑按位存储

### skiplist
key，value = (分词msg, mysqlid)

最大层高48
性能测试：   
cpu: AMD Ryzen 7 5800H with Radeon Graphics   
全部的分词索引数据一共**12193028**个key-value对   
插入索引时间：**1min**左右   
查询key时间：**4us**左右
根据docid查询mysql时间：**0.2s**左右   
Todo：   
(1)参考`redis`高效的`randlevel`算法   
(2)更高效的查询算法   

### 正排索引
	=== RUN   TestGetDocuments
	[376.173ms] [rows:1208]
	--- PASS: TestGetDocuments (0.49s)
	PASS