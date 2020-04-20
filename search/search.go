package search

var esearch string = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?"

type Query struct {
}

// GetPmids ...
func GetPmids(query string) []int {
	pmids := make([]int, 3)

	return pmids
}

func (q *Query) BuildQuery(stringQuery string) Query {
	return Query{}
}
