package tests

import (
	"bufio"
	"fmt"
	"github.com/roy2220/plainkv"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBPlusTree_create(t *testing.T) {
	dict, err := plainkv.OpenOrderedDict("../../searcher/data/index/ordereddict.txt", true)
	if err != nil {
		panic(err)
	}
	defer dict.Close()

	fileReader, err := os.Open("../../searcher/data/terms/dictionary1.txt")
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
		dict.Set([]byte(splits[0]), []byte(splits[1]), false)
	}

}

func TestBPlusTree_query(t *testing.T) {
	dict, err := plainkv.OpenOrderedDict("../../searcher/data/index/ordereddict.txt", true)
	if err != nil {
		panic(err)
	}
	defer dict.Close()
	value, ok := dict.Test([]byte("队员"), true)
	fmt.Printf("%v %q\n", ok, value)
}
