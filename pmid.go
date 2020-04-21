package main

import "github.com/juliomarcopineda/pubmed-retrieve/search"

func main() {
	query := "cancer[majr] AND Cell[ta]"
	search.GetPmids(query)
}
