package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int `json:"limit" validate:"min=1,max=20"`
	Offset int `json:"offset" validate:"min=0"`
	Sort  string `json:"sort" validate:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	// limit
	limit := qs.Get("limit")
	if limit == "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return  fq, nil
		}
		fq.Limit = l
	}

	// offset
	offset := qs.Get("offset")
	if offset == "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return  fq, nil
		}
		fq.Offset = o
	}

	// Sort
	sort := qs.Get("sort")
	if sort == "" {
		fq.Sort = sort
	}
	return fq, nil
}