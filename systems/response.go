package systems

import (
	"crypto/tls"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type TotalData struct {
	TotalCount int `json:"total_data"`
}

type Result struct {
	Data []map[string]interface{} `json: "data"`
}

type Links struct {
	Self  string  `json:"self,omitempty"`
	First string  `json:"first,omitempty"`
	Prev  *string `json:"prev,omitempty"`
	Next  *string `json:"next,omitempty"`
	Last  *string `json:"last,omitempty"`
}

type PaginationResponse struct {
	URI          string
	Path         string
	QueryStrings url.Values
	TotalData    int
	TLS          *tls.ConnectionState
	PageNumber   int
	PageLimit    int
}

func (pr *PaginationResponse) BuildPaginationLinks(request *http.Request, totalData int) *Links {
	// Example: localhost:8080
	pr.URI = request.Host

	// Example: /v1/shopping_lists/items
	pr.Path = request.URL.Path

	// Example: map[page[limit]:[10] last_update:[2016-09-29T17:14:51Z] page[offset]:[10]]
	pr.QueryStrings = request.URL.Query()

	// Example: 12537
	pr.TotalData = totalData

	pr.TLS = request.TLS

	// Set Page Number equal to the page number in the request URI
	pageNumber, _ := strconv.Atoi(pr.QueryStrings.Get("page_number"))
	pr.PageNumber = pageNumber

	// Set Page Number default as 1 if query string page[number] doesn't exist in the request URI
	if pr.QueryStrings.Get("page_number") == "" || pr.QueryStrings.Get("page_number") == "0" {
		pr.PageNumber = 1
	}

	// Convert page limit value from string to integer
	pageLimit, _ := strconv.Atoi(pr.QueryStrings.Get("page_limit"))
	pr.PageLimit = pageLimit

	paginationLinks := Links{}

	// Build all pagination links
	paginationLinks.Self = pr.buildSelfLink()
	paginationLinks.First = pr.buildFirstLink()
	paginationLinks.Last = pr.buildLastLink()
	paginationLinks.Next = pr.buildNextLink()
	paginationLinks.Prev = pr.buildPreviousLink()

	lastLinkPageNumber := math.Ceil(float64(pr.TotalData) / float64(pr.PageLimit))

	if pr.PageNumber == 1 || (paginationLinks.Next == nil && float64(pr.PageNumber) > lastLinkPageNumber) {
		paginationLinks.Prev = nil
	}

	return &paginationLinks
}

func (pr *PaginationResponse) buildSelfLink() string {
	selfLinkQueryString := pr.QueryStrings

	selfLinkPageNumber := strconv.Itoa(pr.PageNumber)
	selfLinkPagelimit := strconv.Itoa(pr.PageLimit)

	selfLinkQueryString.Set("page_number", selfLinkPageNumber)
	selfLinkQueryString.Set("page_limit", selfLinkPagelimit)

	return pr.GetURIScheme() + pr.URI + pr.Path + "?" + selfLinkQueryString.Encode()
}

func (pr *PaginationResponse) buildFirstLink() string {
	firstLinkQueryString := pr.QueryStrings

	firstLinkPageNumber := strconv.Itoa(1)
	firstLinkPagelimit := strconv.Itoa(pr.PageLimit)

	firstLinkQueryString.Set("page_number", firstLinkPageNumber)
	firstLinkQueryString.Set("page_limit", firstLinkPagelimit)

	return pr.GetURIScheme() + pr.URI + pr.Path + "?" + firstLinkQueryString.Encode()
}

func (pr *PaginationResponse) buildNextLink() *string {
	nextLinkQueryString := pr.QueryStrings

	lastLinkPageNumber := math.Ceil(float64(pr.TotalData) / float64(pr.PageLimit))
	nextLinkPageNumber := pr.PageNumber + 1

	if float64(nextLinkPageNumber) > lastLinkPageNumber {
		return nil
	}

	nextLinkPageNumberStr := strconv.Itoa(pr.PageNumber + 1)
	nextLinkPagelimitStr := strconv.Itoa(pr.PageLimit)

	nextLinkQueryString.Set("page_number", nextLinkPageNumberStr)
	nextLinkQueryString.Set("page_limit", nextLinkPagelimitStr)

	nextLinkURI := pr.GetURIScheme() + pr.URI + pr.Path + "?" + nextLinkQueryString.Encode()

	return &nextLinkURI
}

func (pr *PaginationResponse) buildPreviousLink() *string {
	previousLinkQueryString := pr.QueryStrings

	previousLinkPageNumber := strconv.Itoa(pr.PageNumber - 1)
	previousLinkPagelimit := strconv.Itoa(pr.PageLimit)

	previousLinkQueryString.Set("page_number", previousLinkPageNumber)
	previousLinkQueryString.Set("page_limit", previousLinkPagelimit)

	previousLinkURI := pr.GetURIScheme() + pr.URI + pr.Path + "?" + previousLinkQueryString.Encode()

	return &previousLinkURI
}

func (pr *PaginationResponse) buildLastLink() *string {
	lastLinkQueryString := pr.QueryStrings

	lastLinkPageNumber := math.Ceil(float64(pr.TotalData) / float64(pr.PageLimit))

	lastLinkQueryString.Set("page_number", strconv.FormatFloat(lastLinkPageNumber, 'f', 0, 64))

	if lastLinkPageNumber < 1 {
		lastLinkQueryString.Set("page_number", "1")
	}

	lastLinkPagelimit := strconv.Itoa(pr.PageLimit)

	lastLinkQueryString.Set("page_limit", lastLinkPagelimit)

	lastLinkURI := pr.GetURIScheme() + pr.URI + pr.Path + "?" + lastLinkQueryString.Encode()

	return &lastLinkURI
}

// GetURIScheme function used to retrieve URI Scheme http or https
func (pr *PaginationResponse) GetURIScheme() string {
	if pr.TLS != nil {
		return "https://"
	}

	return "http://"
}
