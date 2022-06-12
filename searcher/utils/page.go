package utils

import (
	"math"
)

func CreatePaging(page, pagesize, total int) *Paging {
	if page < 1 {
		page = 1
	}
	if pagesize < 1 {
		pagesize = 10
	}

	page_count := math.Ceil(float64(total) / float64(pagesize))

	paging := new(Paging)
	paging.Page = page
	paging.Pagesize = pagesize
	paging.Total = total
	paging.PageCount = int(page_count)
	paging.NumsCount = 7
	paging.setNums()
	return paging
}

type Paging struct {
	Page      int   //当前页
	Pagesize  int   //每页条数
	Total     int   //总条数
	PageCount int   //总页数
	Nums      []int //分页序数
	NumsCount int   //总页序数
}

func (this *Paging) setNums() {
	this.Nums = []int{}
	if this.PageCount == 0 {
		return
	}

	half := math.Floor(float64(this.NumsCount) / float64(2))
	begin := this.Page - int(half)
	if begin < 1 {
		begin = 1
	}

	end := begin + this.NumsCount - 1
	if end >= this.PageCount {
		begin = this.PageCount - this.NumsCount + 1
		if begin < 1 {
			begin = 1
		}
		end = this.PageCount
	}

	for i := begin; i <= end; i++ {
		this.Nums = append(this.Nums, i)
	}
}
