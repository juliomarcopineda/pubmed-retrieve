package main

import (
	"fmt"

	"github.com/juliomarcopineda/pubmed-retrieve/search"
)

func main() {
	pubmedQuery := "cancer[majr] AND Cell[ta]"
	pubmedJSON, err := search.Retrieve(pubmedQuery)
	_ = err
	fmt.Println(string(pubmedJSON))
}
