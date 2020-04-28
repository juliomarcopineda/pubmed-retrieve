package search

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "Failed to get PMIDs")
	}

	eFetchParams := map[string]string{
		"db":      "pubmed",
		"id":      strings.Join(pmids, ","),
		"retmode": "xml",
	}

	eFetchURL, err := setupURL(eFetch, eFetchParams)
	if err != nil {
		return nil, errors.Wrap(err, "Failed parsing URL")
	}

	body, err := getXML(eFetchURL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get XML")
	}

	var pubmedArticleSet PubmedArticleSet
	err = xml.NewDecoder(body).Decode(&pubmedArticleSet)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert XML to PubmedArticleSet")
	}

	jsonResult, err := json.Marshal(pubmedArticleSet)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert PubmedArticleSet to JSON")
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
		return nil, errors.Wrap(err, "Failed parsing URL")
	}

	body, err := getXML(eSearchURL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get XML")
	}

	var eSearchResult eSearchResult
	err = xml.NewDecoder(body).Decode(&eSearchResult)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert XML to eSearchResult")
	}

	body.Close()
	return eSearchResult.Pmids, nil
}

func setupURL(urlString string, params map[string]string) (string, error) {
	var result string

	urlParse, err := url.Parse(urlString)
	if err != nil {
		return result, errors.Wrap(err, "Failed parsing URL")
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
		return nil, errors.Wrap(err, "GET error with"+url)
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.Wrap(err, "HTTP status code not OK")
	}

	return response.Body, nil
}
