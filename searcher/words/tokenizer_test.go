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
	//[计算 在 本京 大学 深造 小明 硕士 那 毕业于 科学院 中国 所后 日 都]
}
func TestTokenizer_CutContent(t *testing.T) {
	tokenizer := NewTokenizer()
	results := tokenizer.CutContent("小明硕士那那那那那那那那那那那那那那那那毕业于中国科学院计算所，后在日本京都大学深造")
	fmt.Println(results)
	//[科学院 计算 大学 硕士 中国 本京 深造 毕业于 在 日 小明 那 所后 都]
}
