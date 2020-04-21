package main

import (
	"fmt"
	"log"

	"github.com/juliomarcopineda/pubmed-retrieve/search"
)

func main() {
	pubmedQuery := "cancer[majr] AND Cell[ta]"
	pmids, err := search.GetPmids(pubmedQuery)
	if err != nil {
		log.Fatal(err)
	}

	for _, pmid := range pmids {
		fmt.Println(pmid)
	}
}
