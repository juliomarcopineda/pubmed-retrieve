package search

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const eSearchURLString string = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?"

type eSearchResult struct {
	XMLName xml.Name `xml:"eSearchResult"`
	Count   int      `xml:"Count"`
	Pmids   pmids    `xml:"IdList"`
}

type pmids struct {
	PmidSlice []int `xml:"Id"`
}

// GetPmids returns a slice of PMIDs given a PubMed query
func GetPmids(pubmedQuery string) ([]int, error) {
	pmidQuery := map[string]string{
		"db":     "pubmed",
		"term":   pubmedQuery,
		"retmax": "100000",
	}

	eSearchURL, err := url.Parse(eSearchURLString)
	if err != nil {
		return nil, fmt.Errorf("URL Parse Error: %v", err)
	}

	eSearchQuery := eSearchURL.Query()
	for key, val := range pmidQuery {
		eSearchQuery.Set(key, val)
	}
	eSearchURL.RawQuery = eSearchQuery.Encode()

	body, err := getXML(eSearchURL.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to get XML Body: %v", err)
	}

	var result eSearchResult
	err = xml.NewDecoder(body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body: %v", err)
	}

	body.Close()
	return result.Pmids.PmidSlice, nil
}

// getXML returns the XML body of a given URL string as an io.ReadCloser
func getXML(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", response.StatusCode)
	}

	return response.Body, nil
}
