package utils

import (
	"fmt"
	"testing"
)

func TestReadByMMAP(t *testing.T) {
	file, _ := ReadByMMAP("../data/terms/dictionary1.txt")
	fmt.Println(file)
	//fmt.Println(data)
}

func TestGetAllKVS(t *testing.T) {
	_, data := ReadByMMAP("../data/terms/dictionary1.txt")
	k, v := GetAllKVS(data)
	fmt.Println(len(k), len(v))
	//for i := 0; i < len(k); i++ {
	//	fmt.Println(k[i], v[i])
	//}
}
func TestSplitDocIdsFromValue(t *testing.T) {
	str := "[85143703 84065217 91118832 86303864 80433727 92413791 89688037 89637972 84730563 81593687 88697050 87845263 84398456 86250428 86293710 82034631 84621022 82498253]"
	value := SplitDocIdsFromValue(str)
	for _, item := range value {
		fmt.Println(item)
	}
}
func TestGetAllByteArrayKV(t *testing.T) {
	_, data := ReadByMMAP("../data/terms/dictionary1.txt")
	k, v := GetAllByteArrayKV(data)
	fmt.Println(len(k), len(v))
}
