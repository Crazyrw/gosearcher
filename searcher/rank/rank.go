package rank

import (
	"fmt"
	"goSearcher/searcher/utils"
	"goSearcher/searcher/words"
	"math"
	"sort"
	"strings"
)

type rankRes struct {
	docId int
	score float64
}

//sort desc by score
type rankScoresSliceDecrement []rankRes

func (r rankScoresSliceDecrement) Len() int           { return len(r) }
func (r rankScoresSliceDecrement) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r rankScoresSliceDecrement) Less(i, j int) bool { return r[i].score > r[j].score }
func Rank(query string, docIds []int, terms []string, lens []int) (docIdsRank []int) {
	//get all document length in database
	//db.MysqlDB.Where()
	//allDocumentsLength
	documents := utils.GetDocuments(docIds)
	var docIdsForRank []int //存放要排序的docid
	var docs []string       //存放要排序的caption
	for _, item := range documents {
		docs = append(docs, item.Caption)
		docIdsForRank = append(docIdsForRank, item.ID)
	}

	//tokenizer := words.NewTokenizer()
	//queryWords := tokenizer.CutContent(query)         //query的分词结果
	//IDFs := make(map[string]float64, len(terms)) //保存query中每个word的idf值
	//for _, word := range terms {
	//	var num float64 = 0
	//	for _, doc := range docs {
	//		docWords := tokenizer.CutContent(doc) //doc的分词结果
	//		flag := Find(docWords, word)
	//		if flag == true {
	//			num++
	//		}
	//	}
	//	//df := num / float64(len(docs))
	//	idf := math.Log2(float64(len(docs)) / num)
	//	IDFs[word] = idf
	//}
	var docIdsScore []rankRes
	allDocumentsNum := 92700000
	for _, doc := range documents {
		docScoreSum := 0.0
		for i, term := range terms {
			count := strings.Count(doc.Caption, term)
			tf := float64(count) / math.Log10(float64(len(doc.Caption)))
			idf := math.Log10(float64(allDocumentsNum)) / math.Log10(0.01+float64(lens[i]))
			docScoreSum += tf * idf
		}
		docIdsScore = append(docIdsScore, rankRes{
			docId: doc.ID,
			score: docScoreSum,
		})
	}

	//docsScore := make([]float64, len(docs), len(docs))
	//for k, doc := range docs {
	//	score := getScore(query, doc, IDFs)
	//	docsScore[k] = score
	//}
	//
	////保存所有的docId和对应的score
	//var rankScores []rankRes
	//for index, score := range docsScore {
	//	rankSc := rankRes{
	//		docId: docIdsForRank[index],
	//		score: score,
	//	}
	//	rankScores = append(rankScores, rankSc)
	//}
	//
	//sort rankScores desc
	sort.Sort(rankScoresSliceDecrement(docIdsScore))
	fmt.Println(docIdsScore)
	//get docIds by desc sort
	for _, item := range docIdsScore {
		docIdsRank = append(docIdsRank, item.docId)
	}
	return
}

//获取query对于某个doc的tf-idf值
func getScore(query string, doc string, IDFs map[string]float64) (score float64) {
	score = 0
	//var DF float64 = 0
	tokenizer := words.NewTokenizer()
	queryWords := tokenizer.CutContent(query) //query的分词结果
	docWords := tokenizer.CutDoc(doc)         //doc的分词结果
	for _, word := range queryWords {
		count := Count(docWords, word)
		var tf float64 = count / float64(len(docWords))
		idf := IDFs[word]
		score += tf * idf
	}

	return score
}

//判断string切片中含某个string的个数
func Count(slice []string, val string) (count float64) {
	count = 0
	for _, item := range slice {
		if item == val {
			count++
		}
	}
	return count
}

//判断string切片中是否含某个string
func Find(slice []string, val string) (flag bool) {
	flag = false
	for _, item := range slice {
		if item == val {
			flag = true
			break
		}
	}
	return flag
}
