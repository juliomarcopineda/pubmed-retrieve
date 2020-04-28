package search

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	eSearch = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?"
	eFetch  = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?"
)

// PubmedArticleSet ...
type PubmedArticleSet struct {
	XMLName        xml.Name        `xml:"PubmedArticleSet"`
	PubmedArticles []pubmedArticle `xml:"PubmedArticle"`
}

type pubmedArticle struct {
	Pmid        int      `xml:"MedlineCitation>PMID"`
	JournalName string   `xml:"MedlineCitation>Article>Journal>Title"`
	Abstract    string   `xml:"MedlineCitation>Article>Abstract>AbstractText"`
	PubDate     pubDate  `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate"`
	AuthorList  []author `xml:"MedlineCitation>Article>AuthorList>Author"`
}

type pubDate struct {
	Year  int `xml:"Year"`
	Month int `xml:"Month,omitempty"`
	Day   int `xml:"Day"`
}

type author struct {
	LastName    string `xml:"LastName"`
	ForeName    string `xml:"ForeName,omitempty"`
	Affiliation string `xml:"AffiliationInfo>Affiliation,omitempty"`
}

type eSearchResult struct {
	XMLName xml.Name `xml:"eSearchResult"`
	Count   int      `xml:"Count"`
	Pmids   []string `xml:"IdList>Id"`
}

// Retrieve TODO comment
func Retrieve(pubmedQuery string) {
	pmids, err := GetPmids(pubmedQuery)
	if err != nil {
		return
	}

	eFetchParams := map[string]string{
		"db":      "pubmed",
		"id":      strings.Join(pmids, ","),
		"retmode": "xml",
	}

	eFetchURL, err := setupURL(eFetch, eFetchParams)
	if err != nil {
		return
	}

	body, err := getXML(eFetchURL)
	if err != nil {
		return
	}

	var pubmedArticleSet PubmedArticleSet
	err = xml.NewDecoder(body).Decode(&pubmedArticleSet)
	if err != nil {
		fmt.Println(fmt.Errorf("Failed to read body: %v", err))
		return
	}

	_ = pubmedArticleSet
	body.Close()
}

// GetPmids returns a slice of PMIDs given a PubMed query
func GetPmids(pubmedQuery string) ([]string, error) {
	eSearchParams := map[string]string{
		"db":     "pubmed",
		"term":   pubmedQuery,
		"retmax": "10",
	}

	eSearchURL, err := setupURL(eSearch, eSearchParams)
	if err != nil {
		return nil, fmt.Errorf("URL Parse Error: %v", err)
	}

	body, err := getXML(eSearchURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to get XML Body: %v", err)
	}

	var eSearchResult eSearchResult
	err = xml.NewDecoder(body).Decode(&eSearchResult)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body: %v", err)
	}

	body.Close()
	return eSearchResult.Pmids, nil
}

func setupURL(urlString string, params map[string]string) (string, error) {
	var result string

	urlParse, err := url.Parse(urlString)
	if err != nil {
		return result, fmt.Errorf("URL Parse Error: %v", err)
	}

	urlQuery := urlParse.Query()
	for key, val := range params {
		urlQuery.Set(key, val)
	}
	urlParse.RawQuery = urlQuery.Encode()
	result = urlParse.String()

	return result, nil
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
