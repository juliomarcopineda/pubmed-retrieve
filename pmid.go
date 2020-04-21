package main

import (
	"fmt"
	"log"

	"github.com/juliomarcopineda/pubmed-retrieve/search"
)

func main() {
	query := "cancer[majr] AND Cell[ta]"
	pmids, err := search.GetPmids(query)
	if err != nil {
		log.Fatal(err)
	}

	for _, pmid := range pmids {
		fmt.Println(pmid)
	}
}
