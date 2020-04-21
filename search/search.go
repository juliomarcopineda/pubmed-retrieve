package search

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var esearchURLString string = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?"

// GetPmids ...
func GetPmids(query string) ([]int, error) {
	pmidQuery := map[string]string{
		"db":     "pubmed",
		"term":   query,
		"retmax": "100000",
	}

	esearchURL, err := url.Parse(esearchURLString)
	if err != nil {
		return nil, fmt.Errorf("URL Parse Error: %v", err)
	}

	esearchQuery := esearchURL.Query()
	for key, val := range pmidQuery {
		esearchQuery.Set(key, val)
	}
	esearchURL.RawQuery = esearchQuery.Encode()

	xmlString, err := getXML(esearchURL.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to get XML: %v", err)
	}
	fmt.Println(xmlString)

	return nil, nil
}

// getXML returns the XML body of a given URL string
func getXML(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status error: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Read body: %v", err)
	}

	fmt.Println(string(data))

	response.Body.Close()

	return string(data), nil
}
