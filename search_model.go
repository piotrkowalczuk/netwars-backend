package main

const (
	SEARCH_RANGE_POSTS  = 1
	SEARCH_RANGE_TOPICS = 2
	SEARCH_RANGE_BOTH   = 3
)

type SearchParams struct {
	Text  string
	Range int64
}
