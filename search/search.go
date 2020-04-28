package search

import (
	"encoding/json"
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
	XMLName        xml.Name        `xml:"PubmedArticleSet" json:"-"`
	PubmedArticles []pubmedArticle `xml:"PubmedArticle" json:"pubmedArticles"`
}

type pubmedArticle struct {
	Pmid       int      `xml:"MedlineCitation>PMID" json:"_id"`
	Journal    string   `xml:"MedlineCitation>Article>Journal>Title" json:"journal"`
	Abstract   string   `xml:"MedlineCitation>Article>Abstract>AbstractText" json:"abstract"`
	PubDate    pubDate  `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate" json:"pubDate"`
	AuthorList []author `xml:"MedlineCitation>Article>AuthorList>Author" json:"authors"`
}

type pubDate struct {
	Year  int `xml:"Year" json:"year"`
	Month int `xml:"Month,omitempty" json:"month,omitempty"`
	Day   int `xml:"Day,omitempty" json:"day,omitempty"`
}

type author struct {
	LastName    string `xml:"LastName" json:"lastName"`
	ForeName    string `xml:"ForeName,omitempty" json:"foreName"`
	Affiliation string `xml:"AffiliationInfo>Affiliation,omitempty" json:"affiliation"`
}

type eSearchResult struct {
	XMLName xml.Name `xml:"eSearchResult"`
	Count   int      `xml:"Count"`
	Pmids   []string `xml:"IdList>Id"`
}

// Retrieve TODO comment
func Retrieve(pubmedQuery string) ([]byte, error) {
	pmids, err := GetPmids(pubmedQuery)
	if err != nil {
		return nil, fmt.Errorf("Failed to get PMIDs: %v", err)
	}

	eFetchParams := map[string]string{
		"db":      "pubmed",
		"id":      strings.Join(pmids, ","),
		"retmode": "xml",
	}

	eFetchURL, err := setupURL(eFetch, eFetchParams)
	if err != nil {
		return nil, fmt.Errorf("URL Parse Error: %v", err)
	}

	body, err := getXML(eFetchURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body: %v", err)
	}

	var pubmedArticleSet PubmedArticleSet
	err = xml.NewDecoder(body).Decode(&pubmedArticleSet)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert from XML: %v", err)
	}

	jsonResult, err := json.Marshal(pubmedArticleSet)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert to JSON: %v", err)
	}

	body.Close()
	return jsonResult, nil
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
