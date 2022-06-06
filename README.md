### 索引创建
    1.内存版 B+树测试
    === RUN   TestCreateMemoryBtree
	12193026 12193026
	1m27.607646013s
	2022/06/06 11:03:40 create memory btree cost  1m27.607646013s
	2022/06/06 11:03:40 query 二胎( 1210 ) cost  18.483µs
	--- PASS: TestCreateMemoryBtree (87.70s)
	PASS

    2.内存版 SkipList测试(不知道什么原因 十分的慢)
    === RUN   TestCreateSkipList
	12193026
	2022/06/06 14:14:48 create skiplist cost  19m10.43155128s
	2022/06/06 14:14:57 query 二胎( 10891 ) cost  20.797µs
	--- PASS: TestCreateSkipList (1159.41s)
	PASS

    落盘B+树存在问题 一个磁盘块不足以放下一个较大的node >4K

### skiplist
key，value = (分词msg，mysqlid)

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