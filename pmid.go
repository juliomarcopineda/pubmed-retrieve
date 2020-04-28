package main

import (
	"github.com/juliomarcopineda/pubmed-retrieve/search"
)

func main() {
	pubmedQuery := "cancer[majr] AND Cell[ta]"
	search.Retrieve(pubmedQuery)
}
