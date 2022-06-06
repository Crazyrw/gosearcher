package words

import (
	"github.com/wangbin/jiebago"
	"goSearcher/searcher/utils"
	"strings"
)

type Tokenizer struct {
	seg jiebago.Segmenter
}

func NewTokenizer() *Tokenizer {
	tokenizer := &Tokenizer{}
	tokenizer.seg.LoadDictionary("dic.txt")
	return tokenizer
}

// Cut 对文本进行分词
func (tokenizer *Tokenizer) Cut(text string) []string {
	//不区分大小写
	text = strings.ToLower(text)
	//移除标点符号
	text = utils.RemovePunctuation(text)
	//移除所有的空格
	text = utils.RemoveSpace(text)

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
	return wordsSlice
}
