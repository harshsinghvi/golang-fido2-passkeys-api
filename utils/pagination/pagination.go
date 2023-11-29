package pagination

import (
	"strconv"
)

type Pagination struct {
	TotalRecords int `json:"total_records"`
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	NextPage     int `json:"next_page"`
	PrevPage     int `json:"prev_page"`
}

func (pag *Pagination) Validate() *Pagination {
	pag.Set(pag.CurrentPage, pag.TotalRecords)
	return pag
}

func (pag *Pagination) Set(current int, totalRec int) *Pagination {
	pag.TotalRecords = totalRec

	if pag.TotalPages = pag.TotalRecords / 10; pag.TotalRecords%10 == 0 {
		pag.TotalPages = pag.TotalPages - 1
	}

	if current <= 0 {
		pag.PrevPage = 0
	} else {
		pag.PrevPage = current - 1
	}

	if current >= pag.TotalPages {
		pag.NextPage = pag.TotalPages
	} else {
		pag.NextPage = current + 1
	}

	pag.CurrentPage = current
	return pag
}

func (pag *Pagination) ParseString(pageString string) {
	pag.CurrentPage = 0

	if pageString == "" {
		return
	}

	pag.CurrentPage, _ = strconv.Atoi(pageString)

	if pag.CurrentPage == -1 {
		return
	}

	if pag.CurrentPage < 0 {
		pag.CurrentPage = 0
		return
	}
}
