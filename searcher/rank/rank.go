package rank

import (
	"goSearcher/searcher/model"
	"goSearcher/searcher/utils"
	"math"
	"sort"
	"strings"
)

type rankRes struct {
	document model.Docs
	score    float64
}
type DocumentPos struct {
	Document model.Docs
	Pos      []Position
}
type Position struct {
	Start int
	End   int
}

//sort desc by score
type rankScoresSliceDecrement []rankRes

func (r rankScoresSliceDecrement) Len() int           { return len(r) }
func (r rankScoresSliceDecrement) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r rankScoresSliceDecrement) Less(i, j int) bool { return r[i].score > r[j].score }
func Rank(docIds []int, terms []string, lens []int) []model.Docs {
	documents := utils.GetDocuments(docIds)
	var docIdsForRank []int //存放要排序的docid
	var docs []string       //存放要排序的caption
	for _, item := range documents {
		docs = append(docs, item.Caption)
		docIdsForRank = append(docIdsForRank, item.ID)
	}
	var documentsScore []rankRes
	allDocumentsNum := 92700000
	for _, doc := range documents {
		docScoreSum := 0.0
		for i, term := range terms {
			count := strings.Count(doc.Caption, term)
			tf := float64(count) / math.Log10(float64(len(doc.Caption)))
			idf := math.Log10(float64(allDocumentsNum)) / math.Log10(0.01+float64(lens[i]))
			docScoreSum += tf * idf
		}
		documentsScore = append(documentsScore, rankRes{
			document: doc,
			score:    docScoreSum,
		})
	}
	//sort rankScores desc
	sort.Sort(rankScoresSliceDecrement(documentsScore))
	//fmt.Println(docIdsScore)
	//get docIds by desc sort
	// var docIdsRank []int
	// for _, item := range docIdsScore {
	// 	docIdsRank = append(docIdsRank, item.docId)
	// }
	documentPos := hightLight(terms, documentsScore)
	// fmt.Println(documentPos)
	spanStart := "<span style=\"color:red\">"
	spanEnd := "</span>"
	finalDocs := make([]model.Docs, 0)
	for _, docPos := range documentPos {
		for _, pos := range docPos.Pos {
			start := pos.Start
			end := pos.End
			caption := docPos.Document.Caption
			docPos.Document.Caption = caption[:start] + spanStart + caption[start:end+1] + spanEnd
		}
		finalDocs = append(finalDocs, docPos.Document)
	}
	// return documentPos
	return finalDocs
}

func hightLight(terms []string, documentsRank []rankRes) []DocumentPos {
	// fmt.Println(documentsRank)
	var documentPos []DocumentPos
	for _, item := range documentsRank {
		content := item.document.Caption
		var pos []Position
		for _, term := range terms {
			start := strings.Index(content, term)
			if start == -1 {
				continue
			}
			end := start + len(term) - 1
			p := Position{
				Start: start,
				End:   end,
			}
			pos = append(pos, p)
		}
		if pos == nil {
			continue
		}
		documentPosItem := DocumentPos{
			Document: item.document,
			Pos:      pos,
		}
		// fmt.Println(documentPosItem)
		documentPos = append(documentPos, documentPosItem)
	}
	//fmt.Println(documentPos)
	return documentPos
}
