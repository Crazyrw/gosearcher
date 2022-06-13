package words

import (
	"fmt"
	"goSearcher/searcher/utils"
	"testing"
)

func TestTokenizer_Cut(t *testing.T) {
	//segs.LoadDictionary("dict.txt")
	tokenizer := NewTokenizer()
	results := tokenizer.Cut("今日头条是字节的吗")
	fmt.Println(results)
	//[今日头 条 字节]
}
func TestTokenizer_CutContent(t *testing.T) {
	tokenizer := NewTokenizer()
	results := tokenizer.CutContent("今日头条")
	fmt.Println(results)
	//[今日头 条 字节]
}

func TestJieba(t *testing.T) {
	text := "杨紫230米云端相聚##杨紫##花木星球"
	tokenizer := NewTokenizer()
	var wordMap = make(map[string]int)
	resultChan := tokenizer.seg.CutForSearch(text, true)
	for word := range resultChan {
		//去除停用词
		_, ok := utils.AllStopsWords[word]
		if !ok {
			_, found := wordMap[word]
			if !found {
				//去除重复的词  只是针对单个文本去重了
				wordMap[word] = 1
			}
		}
	}
	var wordsSlice []string
	for k := range wordMap {

		wordsSlice = append(wordsSlice, k)
	}
	fmt.Println(wordsSlice)
	//[字节 今日头 条]
}
