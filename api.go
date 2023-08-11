package metaphor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type MetaphorClient struct {
	apiKey      string
	RequestBody *RequestBody
}

type RequestBody struct {
	Query              string   `json:"query,omitempty"`
	Url                string   `json:"url,omitempty"`
	NumResults         int      `json:"numResults,omitempty"`
	IncludeDomains     []string `json:"includeDomains,omitempty"`
	ExcludeDomains     []string `json:"excludeDomains,omitempty"`
	StartCrawlDate     string   `json:"startCrawlDate,omitempty"`
	EndCrawlDate       string   `json:"endCrawlDate,omitempty"`
	StartPublishedDate string   `json:"startPublishedDate,omitempty"`
	EndPublishedDate   string   `json:"endPublishedDate,omitempty"`
	UseAutoprompt      bool     `json:"useAutoprompt,omitempty"`
	Type               string   `json:"type,omitempty"`
}

const (
	// DEFAULT REQUEST VALUES

	// DefaultNumResults is the default number of expected results
	DefaultNumResults = 10
	// DefaultAutoprompt if true, your query will be converted to a Metaphor query.
	// If findLinks ednpoint is used needs to be nil to omit useAutoprompt field from RequestBody
	DefaultAutoprompt = false

	// DEFAULT API ENDPOINT URL's

	// DefaultSearchURL is the default search endpoint
	DefaultSearchURL = "https://api.metaphor.systems/search"
	// DefaultContentsURL is the default contents endpoint
	DefaultContentsURL = "https://api.metaphor.systems/contents"
	// DefaultFindLinksURL is the default find links endpoint
	DefaultFindLinksURL = "https://api.metaphor.systems/findSimilar"
)

var (
	ErrRequestFailed          = errors.New("request failed with error")
	ErrSearchFailed           = errors.New("search failed with error")
	ErrFindSimilarLinkdFailed = errors.New("find similar links failed with error")
	ErrGetContentsFailed      = errors.New("get contents failed with error")
	ErrNoSearchResults        = errors.New("no search results were found")
	ErrNoLinksFound           = errors.New("no links were found")
	ErrNoContentExtracted     = errors.New("no content was extracted")
)

// NewClient creates a new MetaphorClient with the provided API key and options.
//
// Parameters:
// - apiKey: The API key used for authentication.
// - options: Optional client options that can be passed to customize the client.
//
// Returns:
// - *MetaphorClient: A new instance of the MetaphorClient.
// - error: An error if the client creation fails.
func NewClient(apiKey string, options ...ClientOptions) (*MetaphorClient, error) {
	client := &MetaphorClient{
		apiKey: apiKey,
		RequestBody: &RequestBody{
			NumResults:    DefaultNumResults,
			UseAutoprompt: DefaultAutoprompt,
		},
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

// Search searches for a given query using the Metaphor client.
//
// ctx: The context.Context for the request.
// query: The search query.
// options: The optional client options.
// Return type: *SearchResponse, error.
func (client *MetaphorClient) Search(ctx context.Context, query string, options ...ClientOptions) (*SearchResponse, error) {
	for _, option := range options {
		option(client)
	}

	var searchResults = &SearchResponse{}
	client.RequestBody.Query = query

	reqBytes, err := json.Marshal(client.RequestBody)
	if err != nil {
		return searchResults, fmt.Errorf("%v: %w", ErrSearchFailed, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", DefaultSearchURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return searchResults, fmt.Errorf("%v: %w", ErrSearchFailed, err)
	}

	responseBody, err := client.runRequest(ctx, req)
	if err != nil {
		return searchResults, fmt.Errorf("%v: %w", ErrSearchFailed, err)
	}

	err = json.Unmarshal(responseBody, &searchResults)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", ErrSearchFailed, err)
	}

	if len(searchResults.Results) == 0 {
		return searchResults, ErrNoSearchResults
	}

	return searchResults, nil
}

// FindSimilar searches for similar urls using the provided URL.
//
// ctx: The context.Context for the function.
// url: The URL to search for similar items.
// options: Optional client options.
// Return type: (*SearchResponse, error).
func (client *MetaphorClient) FindSimilar(ctx context.Context, url string, options ...ClientOptions) (*SearchResponse, error) {
	for _, option := range options {
		option(client)
	}

	searchResults := &SearchResponse{}
	client.RequestBody.Url = url

	reqBytes, err := json.Marshal(client.RequestBody)
	if err != nil {
		return searchResults, fmt.Errorf("%v: %w", ErrFindSimilarLinkdFailed, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", DefaultFindLinksURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return searchResults, fmt.Errorf("%v: %w", ErrFindSimilarLinkdFailed, err)
	}

	responseBody, err := client.runRequest(ctx, req)
	if err != nil {
		return searchResults, fmt.Errorf("%v: %w", ErrFindSimilarLinkdFailed, err)
	}

	err = json.Unmarshal(responseBody, &searchResults)
	if err != nil {
		return searchResults, fmt.Errorf("%v: %w", ErrFindSimilarLinkdFailed, err)
	}

	if len(searchResults.Results) == 0 {
		return searchResults, ErrNoLinksFound
	}

	return searchResults, nil
}

// GetContents retrieves the contents of urls for the given set of IDs.
//
// It takes the following parameters:
// - ctx: the context.Context for the request.
// - ids: a slice of strings containing the IDs to retrieve the contents for.
//
// It returns a *ContentsResponse which contains the retrieved contents,
// and an error if any occurred.
func (client *MetaphorClient) GetContents(ctx context.Context, ids []string) (*ContentsResponse, error) {
	contentsResults := &ContentsResponse{}
	joinedIds := strings.Join(ids, "\",\"")
	url := fmt.Sprintf("%s?ids=\"%s\"", DefaultContentsURL, joinedIds)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return contentsResults, fmt.Errorf("%v: %w", ErrGetContentsFailed, err)
	}

	responseBody, err := client.runRequest(ctx, req)
	if err != nil {
		return &ContentsResponse{}, fmt.Errorf("%v: %w", ErrGetContentsFailed, err)
	}

	err = json.Unmarshal(responseBody, &contentsResults)
	if err != nil {
		return contentsResults, fmt.Errorf("%v: %w", ErrGetContentsFailed, err)
	}

	if len(contentsResults.Contents) == 0 {
		return contentsResults, ErrNoSearchResults
	}

	return contentsResults, nil
}

// runRequest sends an HTTP request and returns the response body as a byte array.
//
// It takes a context.Context object and an http.Request object as parameters.
// It returns the response body as a byte array and an error object.
func (client *MetaphorClient) runRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	req.Header.Add("x-api-key", client.apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	// trunk-ignore(gokart/CWE-918:-Server-Side-Request-Forgery)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		errorResponse := &ErrorResponse{}
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		errorTxt := errorResponse.Text

		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, errorTxt)
	}

	return body, nil
}



