package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"min=1,max=20"`
	Offset int      `json:"offset" validate:"min=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	// limit
	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}

	// offset
	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = o
	}

	// Sort
	sort := qs.Get("sort")
	if sort == "" {
		fq.Sort = "desc"
	} else {
		fq.Sort = sort
	}

	// tags
	tags := qs.Get("tags")
	if tags != "" {
		// 如果tags是空字符串, 则设置为空切片
		fq.Tags = strings.Split(tags, ",")
	}
	// search
	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}
	// since
	since := qs.Get("since")
	if since != "" {
		sinceTime := parseTime(since)
		fq.Since = sinceTime
	}
	// until
	until := qs.Get("until")
	if until != "" {
		untilTime := parseTime(until)
		fq.Until = untilTime
	}

	return fq, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}
