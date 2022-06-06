package words

import (
	"fmt"
	"testing"
)

func TestTokenizer_Cut(t *testing.T) {
	//segs.LoadDictionary("dict.txt")
	tokenizer := NewTokenizer()
	results := tokenizer.Cut("小明硕士那那那那那那那那那那那那那那那那毕业于中国科学院计算所，后在日本京都大学深造")
	fmt.Println(results)
}
