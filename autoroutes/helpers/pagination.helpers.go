package helpers

import (
	"strconv"
)

type Pagination struct {
	TotalRecords int64
	CurrentPage  int
	TotalPages   int
	NextPage     int
	PrevPage     int
	Limit        int
	Offset       int
}

const DEFAULT_LIMIT = 20

func (pag *Pagination) Validate() *Pagination {
	pag.Set(pag.CurrentPage, pag.TotalRecords)
	return pag
}

func (pag *Pagination) Set(current int, totalRec int64) *Pagination {
	pag.TotalRecords = totalRec
	pag.CurrentPage = current

	if pag.TotalPages = int(pag.TotalRecords / int64(pag.Limit)); pag.TotalRecords%10 == 0 {
		pag.TotalPages = pag.TotalPages - 1
	}

	if current <= 0 {
		pag.PrevPage = 0
	} else {
		pag.PrevPage = current - 1
	}

	if current >= pag.TotalPages {
		pag.NextPage = pag.TotalPages
		pag.PrevPage = pag.TotalPages
		pag.CurrentPage = pag.TotalPages
	} else {
		pag.NextPage = current + 1
	}

	pag.Offset = pag.Limit * pag.CurrentPage
	return pag
}

func (pag *Pagination) ParseString(pageString string) {
	var err error
	pag.CurrentPage = 0

	if pageString == "" {
		return
	}

	pag.CurrentPage, err = strconv.Atoi(pageString)

	if err != nil {
		pag.CurrentPage = -1
		return
	}

	if pag.CurrentPage < 0 {
		pag.CurrentPage = 0
		return
	}
}

func New(pageString string, args ...int) Pagination {
	var pag Pagination
	pag.Limit = DEFAULT_LIMIT

	if args != nil {
		pag.Limit = args[0]
	}

	pag.ParseString(pageString)
	// pag.Validate()
	return pag
}
