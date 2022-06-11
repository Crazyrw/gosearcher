package rank

import (
	"fmt"
	"goSearcher/searcher/utils"
	"math"
	"sort"
	"strings"
)

type rankRes struct {
	docId int
	score float64
}
type DocIdPos struct {
	DocId int
	Pos   []position
}
type position struct {
	start int
	end   int
}

//sort desc by score
type rankScoresSliceDecrement []rankRes

func (r rankScoresSliceDecrement) Len() int           { return len(r) }
func (r rankScoresSliceDecrement) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r rankScoresSliceDecrement) Less(i, j int) bool { return r[i].score > r[j].score }
func Rank(docIds []int, terms []string, lens []int) []DocIdPos {
	documents := utils.GetDocuments(docIds)
	var docIdsForRank []int //存放要排序的docid
	var docs []string       //存放要排序的caption
	for _, item := range documents {
		docs = append(docs, item.Caption)
		docIdsForRank = append(docIdsForRank, item.ID)
	}
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
	//sort rankScores desc
	sort.Sort(rankScoresSliceDecrement(docIdsScore))
	fmt.Println(docIdsScore)
	//get docIds by desc sort
	var docIdsRank []int
	for _, item := range docIdsScore {
		docIdsRank = append(docIdsRank, item.docId)
	}
	return hightLight(terms, docIdsRank)
}

func hightLight(terms []string, docIdsRank []int) (docIdPosition []DocIdPos) {
	documents := utils.GetDocuments(docIdsRank)
	fmt.Println(docIdsRank)

	var newDocIds []int
	for _, item := range documents {
		content := item.Caption
		var pos []position
		for _, term := range terms {
			start := strings.Index(content, term)
			end := start + len(term) - 1
			pos = append(pos, position{
				start: start,
				end:   end,
			})
		}
		docIdPosition = append(docIdPosition, DocIdPos{
			DocId: item.ID,
			Pos:   pos,
		})
		newDocIds = append(newDocIds, item.ID)
	}
	return
}
