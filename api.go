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

const (
	//// DEFAULT REQUEST VALUES.

	// DefaultNumResults is the default number of expected results.
	DefaultNumResults = 10

	// DefaultAutoprompt if true, your query will be converted to a Metaphor query.
	DefaultAutoprompt = false

	// DefaultSearchType is a string defining what type of search will be performed, "neural" or by "keyword".
	DefaultSearchType = "neural"

	//// DEFAULT API ENDPOINT URL's

	// DefaultSearchPath is the default url for metaphor systems api.
	DefaultBaseURL = "https://api.metaphor.systems"

	// DefaultSearchPath is the default search endpoint.
	DefaultSearchPath = "/search"

	// DefaultContentsPath is the default contents endpoint.
	DefaultContentsPath = "/contents"

	// DefaultFindSimilarPath is the default find links endpoint.
	DefaultFindSimilarPath = "/findSimilar"
)

var (
	ErrMissingApiKey = errors.New("missing the Metaphor API key, set it as the METAPHOR_API_KEY environment variable")
	ErrRequestFailed = errors.New("request failed with error")
	ErrSearchFailed = errors.New("search failed with error")
	ErrFindSimilarLinkdFailed = errors.New("find similar links failed with error")
	ErrGetContentsFailed = errors.New("get contents failed with error")
	ErrNoSearchResults = errors.New("no search results were found")
	ErrNoLinksFound = errors.New("no links were found")
	ErrNoContentExtracted = errors.New("no content was extracted")
)

type RequestBody struct {
	Query               string   `json:"query,omitempty"`
	URL                 string   `json:"url,omitempty"`
	NumResults          int      `json:"numResults,omitempty"`
	IncludeDomains      []string `json:"includeDomains,omitempty"`
	ExcludeDomains      []string `json:"excludeDomains,omitempty"`
	StartCrawlDate      string   `json:"startCrawlDate,omitempty"`
	EndCrawlDate        string   `json:"endCrawlDate,omitempty"`
	StartPublishedDate  string   `json:"startPublishedDate,omitempty"`
	EndPublishedDate    string   `json:"endPublishedDate,omitempty"`
	ExcludeSourceDomain bool     `json:"excludeSourceDomain,omitempty"`
	UseAutoprompt       bool     `json:"useAutoprompt,omitempty"`
	Type                string   `json:"type,omitempty"`
}

type Client struct {
	apiKey      string
	options     []ClientOptions
	BaseURL     string
	RequestBody *RequestBody
}

// NewClient creates a new MetaphorClient with the provided API key and options.
//
// Parameters:
// - apiKey: The API key used for authentication.
// - options: Optional client options that can be passed to customize the client.
//
// Returns:
// - *MetaphorClient: A new instance of the MetaphorClient.
// - error: An error if the client creation fails.
func NewClient(apiKey string, options ...ClientOptions) (*Client, error) {
	if apiKey == "" {
		return nil, ErrMissingApiKey
	}

	client := &Client{
		apiKey:      apiKey,
		options:     options,
		BaseURL:     DefaultBaseURL,
		RequestBody: &RequestBody{},
	}

	return client, nil
}

// Search searches for a given query using the Metaphor client.
//
// Parameters:
// - ctx: The context.Context for the request.
// - query: The search query.
// - options: The optional client options.
//
// Returns:
// - *SearchResponse: The search response object.
// - error: An error if the search fails.
func (client *Client) Search(ctx context.Context, query string, options ...ClientOptions) (*SearchResponse, error) {
	searchResults := &SearchResponse{}
	client.RequestBody = &RequestBody{
		Query:         query,
		NumResults:    DefaultNumResults,
		UseAutoprompt: DefaultAutoprompt,
		Type:          DefaultSearchType,
	}

	client.loadOptions(options...)

	reqBytes, err := json.Marshal(client.RequestBody)
	if err != nil {
		return searchResults, fmt.Errorf("%w: %w", ErrSearchFailed, err)
	}

	reqURL := client.BaseURL + DefaultSearchPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return searchResults, fmt.Errorf("%w: %w", ErrSearchFailed, err)
	}

	responseBody, err := client.runRequest(req)
	if err != nil {
		return searchResults, fmt.Errorf("%w: %w", ErrSearchFailed, err)
	}

	err = json.Unmarshal(responseBody, &searchResults)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSearchFailed, err)
	}

	if len(searchResults.Results) == 0 {
		return searchResults, ErrNoSearchResults
	}

	return searchResults, nil
}

// FindSimilar searches for similar urls using the provided URL.
//
// Parameters:
// - ctx: The context.Context for the function.
// - url: The URL to search for similar items.
// - options: Optional client options.
//
// Returns:
// - *SearchResponse: The search response object.
// - error: An error if the search fails.
func (client *Client) FindSimilar(ctx context.Context, url string, options ...ClientOptions) (*SearchResponse, error) {
	searchResults := &SearchResponse{}
	client.RequestBody = &RequestBody{
		URL:           url,
		NumResults:    DefaultNumResults,
		UseAutoprompt: DefaultAutoprompt,
	}

	client.loadOptions(options...)

	reqBytes, err := json.Marshal(client.RequestBody)
	if err != nil {
		return searchResults, fmt.Errorf("%w: %w", ErrFindSimilarLinkdFailed, err)
	}

	reqURL := client.BaseURL + DefaultFindSimilarPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return searchResults, fmt.Errorf("%w: %w", ErrFindSimilarLinkdFailed, err)
	}

	responseBody, err := client.runRequest(req)
	if err != nil {
		return searchResults, fmt.Errorf("%w: %w", ErrFindSimilarLinkdFailed, err)
	}

	err = json.Unmarshal(responseBody, &searchResults)
	if err != nil {
		return searchResults, fmt.Errorf("%w: %w", ErrFindSimilarLinkdFailed, err)
	}

	if len(searchResults.Results) == 0 {
		return searchResults, ErrNoLinksFound
	}

	return searchResults, nil
}

// GetContents retrieves the contents of urls for the given set of IDs.
//
// Parameters:
// - ctx: the context.Context for the request.
// - ids: a slice of strings containing the IDs to retrieve the contents for.
//
// Returns:
// - *ContentsResponse: The contents response object.
// - error: An error if the contents retrieval fails.
func (client *Client) GetContents(ctx context.Context, ids []string) (*ContentsResponse, error) {
	contentsResults := &ContentsResponse{}
	
	client.loadOptions()

	joinedIds := strings.Join(ids, "\",\"")

	reqURL := client.BaseURL + DefaultContentsPath

	URL := fmt.Sprintf("%s?ids=\"%s\"", reqURL, joinedIds)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return contentsResults, fmt.Errorf("%w: %w", ErrGetContentsFailed, err)
	}

	responseBody, err := client.runRequest(req)
	if err != nil {
		return contentsResults, fmt.Errorf("%w: %w", ErrGetContentsFailed, err)
	}

	err = json.Unmarshal(responseBody, &contentsResults)
	if err != nil {
		return contentsResults, fmt.Errorf("%w: %w", ErrGetContentsFailed, err)
	}

	if len(contentsResults.Contents) == 0 {
		return contentsResults, ErrNoSearchResults
	}

	return contentsResults, nil
}

// runRequest sends an HTTP request and returns the response body as a byte array.
//
// Parameters:
// - ctx: the context.Context for the request.
// - req: the HTTP request to send
//
// Returns:
// - []byte: the response body as a byte array
// - error: an error if the request fails
func (client *Client) runRequest(req *http.Request) ([]byte, error) {
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

func (client *Client) loadOptions(options ...ClientOptions) {
	if len(options) > 0 {
		client.options = options
	}

	for _, option := range client.options {
		option(client)
	}
}