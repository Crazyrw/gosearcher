package tests

import (
	"bufio"
	"encoding/json"
	"goSearcher/searcher/btree"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestBtree_check(t *testing.T) {
	//查看初始堆内存
	//runtime.ReadMemStats(&stat)
	//fmt.Println(stat.HeapSys)

	//0.初始化一颗b+树
	btree, err := btree.StartDefaultNewTree()
	if err != nil {
		panic("初始化btree失败")
	}
	////1.读30w条文档数据
	//var documents []model.Docs
	//start := 1
	//end := start + 1
	//result := db.MysqlDB.Where("id >= ? && id <= ?", start, end).Find(&documents)
	//if result.Error != nil {
	//	panic("查询失败")
	//}
	////2.对文本进行切词 <term, []docIds>
	//tokenizer := words.NewTokenizer()
	//kv := make(map[string][]int)
	//for _, document := range documents {
	//	docId := document.ID
	//	text := document.Caption
	//	terms := tokenizer.Cut(text)
	//	for _, term := range terms {
	//		value, ok := kv[term]
	//		if !ok {
	//			//不存在 直接插入
	//			kv[term] = []int{docId}
	//		} else {
	//			value := append(value, docId)
	//			kv[term] = value
	//		}
	//	}
	//}
	//3.将分词结果 键值对存入磁盘
	termsFileName := "dictionary1.txt"
	//fileWriter := utils.OpenFile("../../searcher/data/terms/" + termsFileName)
	//defer fileWriter.Close()
	//for key, value := range kv {
	//	result := fmt.Sprintf("%s,%v\n", key, value)
	//	fileWriter.WriteString(result)
	//}
	//4.插入B+树 key按照字典进行排序
	//读文件 根据文件中的分词结果进行创建b+树
	fileReader, err := os.Open("../../searcher/data/terms/" + termsFileName)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewReader(fileReader)
	for {
		line, _, err := scanner.ReadLine()
		if err == io.EOF {
			break
		}
		data := string(line)
		splits := strings.Split(data, ",")
		btree.Insert(splits[0], splits[1])
	}
	////4.查找"人名"
	//if result, ok := btree.Find("人名"); ok {
	//	fmt.Println(result)
	//}
	//5.将创建的btree树写入磁盘
	fileName := "../../searcher/data/index/index_1.txt"
	//utils.BtreeToDisk(btree, "../../searcher/data/index/"+fileName)
	marshal, err := json.Marshal(btree)
	if err != nil {
		log.Fatalln("json序列化失败 ", err.Error())
	}
	log.Println("json序列化结果 ", marshal)
	fileWriter, err := os.Open(fileName)
	defer fileWriter.Close()
	if err != nil {
		log.Fatalln("序列化结果写入失败...", err.Error())
	}
	fileWriter.Write(marshal)
	//6.使用mmap读文件到内存

	//7.查找
}
