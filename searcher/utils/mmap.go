package utils

import (
	"os"
	"strconv"
	"strings"
	"syscall"
)

const defaultMenuMapSize = 128 * (1 << 20) //假设映射的内存为128M

func ReadByMMAP(path string) (file *os.File, data []byte) {
	stat, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	file, err = os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	//creating a memory of files
	//func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
	/*
		fd: 映射的文件描述符
		offset: 映射到内存区域的起始位置，0表示由内核指定内存地址
		length: 要映射的内存区域的大小
		prot: 内存保护标志位 以通过或运算符`|`组合
			- PROT_EXEC  页内容可以被执行
			- PROT_READ  页内容可以被读取
			- PROT_WRITE 页可以被写入
			- PROT_NONE  页不可访问
		flags：映射对象的类型，常用的是以下两类
			- MAP_SHARED  共享映射，写入数据会复制回文件, 与映射该文件的其他进程共享。
			- MAP_PRIVATE 建立一个写入时拷贝的私有映射，写入数据不影响原文件。
	*/
	data, err = syscall.Mmap(
		int(file.Fd()),
		0,
		int(stat.Size()),
		syscall.PROT_READ,
		syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return file, data
}

// GetAllKVS 返回所有的kv键值对
func GetAllKVS(data []byte) (k, v []string) {
	kvs := strings.Split(string(data), "\n")
	for _, item := range kvs {
		str := strings.Split(item, ",")
		if len(str) == 2 {
			key := str[0]
			value := str[1]
			k = append(k, key)
			v = append(v, value)
		}
	}
	return k, v
}

//-------------------skipList--------------------

// GetAllByteArrayKV
func GetAllByteArrayKV(data []byte) (k [][]byte, v [][]byte) {
	kvs := strings.Split(string(data), "\n")
	for _, item := range kvs {
		str := strings.Split(item, ",")
		if len(str) == 2 {
			key := []byte(str[0])
			value := []byte(str[1])
			k = append(k, key)
			v = append(v, value)
		}
	}
	return k, v
}

// SplitDocIdsFromValue get docIds from value
func SplitDocIdsFromValue(value string) []int {
	split1 := value[1 : len(value)-1]
	split2 := strings.Split(split1, " ")
	var docIds []int
	for _, item := range split2 {
		docId, _ := strconv.Atoi(item)
		docIds = append(docIds, docId)
	}
	return docIds
}
