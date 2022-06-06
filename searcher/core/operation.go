package core

import (
	"bufio"
	"fmt"
	"goSearcher/searcher/btree"
	"goSearcher/searcher/codec"
	"goSearcher/searcher/db"
	"goSearcher/searcher/model"
	"goSearcher/searcher/skip_list"
	"goSearcher/searcher/utils"
	"goSearcher/searcher/words"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//将数据库中所有的文本进行分词 存储到dictionary1.txt
/**
//1.读400w条文档数据 时间很短
	var wg sync.WaitGroup
	step := 1000
	for start := 1; start < 400000; start += step {
		wg.Add(1)
		go func(start, step int) {
			defer wg.Done()
			cutTerms(start, step)
		}(start, step)
	}
	wg.Wait()
*/
//step := 50000
//total := 12700000 //20000000
//cutTerms(step, total)
//TODO: 只能在内存中维护一个map 用于去重 9千多万条数据
func cutTerms(step, total int) {
	var wg sync.WaitGroup

	//1.先将文件中的分词结果读出来 放在kv中
	fileReader, _ := os.Open("searcher/data/terms/dictionary1.txt")
	defer fileReader.Close()
	stat, _ := fileReader.Stat()
	if stat.Size() != 0 {
		fmt.Println("初始文件有数据")
		scanner := bufio.NewScanner(fileReader)
		for scanner.Scan() {
			line := scanner.Text()
			data := string(line)
			//会出现空行的情况
			splits := strings.Split(data, ",")
			if len(splits) == 2 {
				//将splits[1] string类型转换为[]int
				arr := make([]int, 0)
				str2 := splits[1][1 : len(splits[1])-1]
				strs := strings.Split(str2, " ")
				for _, str := range strs {
					v, _ := strconv.Atoi(str)
					arr = append(arr, v)
				}
				kv.Store(splits[0], arr)
			}
		}
	} else {
		fmt.Println("初始文件没有数据！！！")
	}

	//2.开始对新数据进行分词
	tokenizer := words.NewTokenizer()

	//92700000
	//for i := 1; i < 92700000; i += 20000000 {

	//当前这批数据的分词结果存放到nowMap
	//var nowMap sync.Map

	initValue := 80000001
	start := 80000001
	fmt.Println(step, total, start+total)
	for ; start < initValue+total; start += step {
		wg.Add(1)
		go func(start, step int) {
			defer wg.Done()
			//1.从数据库中读取文本数据
			var documents []model.Docs
			result := db.MysqlDB.Where("id >= ? limit ?", start, step).Find(&documents)
			if result.Error != nil {
				panic(result.Error)
			}
			//2.对文本进行切词 <term, []docIds>

			for _, document := range documents {
				docId := document.ID
				text := document.Caption
				terms := tokenizer.Cut(text)
				for _, term := range terms {
					value, ok := kv.Load(term)
					if !ok {
						kv.Store(term, []int{docId})
						//nowMap.Store(term, []int{docId})
					} else {
						//interface{}转换为[]int 接口类型转向普通类型称为类型断言
						value := append(value.([]int), docId)
						kv.Store(term, value)
						//nowMap.Store(term, value)
					}
				}
			}

		}(start, step)
	}
	wg.Wait()
	fmt.Println("(", initValue, "~", initValue+total, ") done.....")

	//写入磁盘
	fileWriter := utils.OpenFile("searcher/data/terms/dictionary1.txt")
	defer fileWriter.Close()
	f := func(key, value interface{}) bool {
		result := fmt.Sprintf("%v,%v\n", key, value)
		fileWriter.WriteString(result)
		return true
	}
	kv.Range(f)

	//3.将分词结果 键值对存入磁盘
	//加入写锁 避免不同的协程同时写
	//f := func(key, value interface{}) bool {
	//	result := fmt.Sprintf("%v,%v\n", key, value)
	//	fileWriter.WriteString(result)
	//	return true
	//}
	//kv.Range(f)
}

//根据分词的结果进行创建b+树 存储到data.db
//TODO: 落盘btree 一个磁盘块无法存储一个节点  有的节点数据过大>4K
func createInvertIndex() {
	start := time.Now()
	//1.创建B+树
	tree, err := btree.NewTree("searcher/data/index/invert.db")
	if err != nil {
		panic(err)
	}
	defer tree.Close()
	//2.将dictionary1.txt分词结果插入b+树中 使用mmap读取
	_, data := utils.ReadByMMAP("searcher/data/terms/dictionary1.txt")
	k, v := utils.GetAllKVS(data)
	fmt.Println(len(k), len(v))
	for i := 0; i < len(k); i++ {
		err := tree.Insert(k[i], v[i])
		if err != nil {
			log.Fatalln("btree Insert error")
		}
	}
	cost := time.Since(start)
	log.Fatalln(cost)
}

// CreateMemoryBtree memory btree
func CreateMemoryBtree(path string) *btree.BPlusTree {
	start := time.Now()
	//使用内存版本的btree
	tree, _ := btree.StartDefaultNewTree()
	//2.将dictionary1.txt分词结果插入b+树中 使用mmap读取
	_, data := utils.ReadByMMAP(path)
	k, v := utils.GetAllKVS(data)
	fmt.Println(len(k), len(v))
	for i := 0; i < len(k); i++ {
		tree.Insert(k[i], v[i])
	}
	cost := time.Since(start)
	log.Println("create memory btree cost ", cost)

	return tree
}

// CreateSkipList create skiplist
func CreateSkipList(path string) *skip_list.SkipList {
	start := time.Now()
	skiplist := skip_list.NewSkipList()
	_, data := utils.ReadByMMAP(path)
	k, v := utils.GetAllByteArrayKV(data)
	fmt.Println(len(k))
	for i := 0; i < len(k); i++ {
		skiplist.Add(codec.NewEntry(k[i], v[i]))
	}
	cost := time.Since(start)
	log.Println("create skiplist cost ", cost)

	start2 := time.Now()
	value := skiplist.Search([]byte("二胎")).Value
	cost2 := time.Since(start2)
	strs := utils.SplitDocIdsFromValue(fmt.Sprintf("%v", value))
	log.Println("query 二胎(", len(strs), ") cost ", cost2)
	return skiplist
}

// GetDocuments get documents by docIds
func GetDocuments(docIds []int) []model.Docs {
	var files []model.Docs
	results := db.MysqlDB.Find(&files, docIds)
	if results.Error != nil {
		log.Fatalln("正排索引失败")
	}
	return files
}
