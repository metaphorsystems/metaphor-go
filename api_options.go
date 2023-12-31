package metaphor

type RequestOptions struct {
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

type ClientOptions func(*Client)

// WithNumResults sets the number of expected search results.
//
// Parameters:
//   - numResults: The desired number of results.
//
// Returns: a ClientOptions function that updates the numResults field of the RequestBody struct.
func WithNumResults(numResults int) ClientOptions {
	return func(client *Client) {
		client.RequestBody.NumResults = numResults
	}
}

// WithIncludeDomains sets the includeDomains field of the RequestBody.
// List of domains to include in the search. If specified, results will
// only come from these domains. Only one of includeDomains and excludeDomains
// should be specified.
//
// Parameters:
// - includeDomains: a slice of strings representing the domains to include.
//
// Returns: a ClientOptions function that updates the includeDomains field of the RequestBody struct.
func WithIncludeDomains(includeDomains []string) ClientOptions {
	return func(client *Client) {
		client.RequestBody.IncludeDomains = includeDomains
	}
}

// WithExcludeDomains sets the ExcludeDomains field of the client's RequestBody.
// List of domains to exclude in the search. If specified, results will only come
// from these domains. Only one of includeDomains and excludeDomains should be specified.
//
// Parameters:
// - excludeDomains: an array of strings representing the domains to be excluded.
//
// Returns: a ClientOptions function that updates the excludeDomains field of the RequestBody struct.
func WithExcludeDomains(excludeDomains []string) ClientOptions {
	return func(client *Client) {
		client.RequestBody.ExcludeDomains = excludeDomains
	}
}

// WithStartCrawlDate sets the start crawl date for the client options.
// If startCrawlDate is specified, results will only include links that
// were crawled after startCrawlDate.
// Must be specified in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)
//
// Parameters:
// - startCrawlDate: the start date for the crawl
//
// Returns: a ClientOptions function that updates the startCrawlDate field of the RequestBody struct.
func WithStartCrawlDate(startCrawlDate string) ClientOptions {
	return func(client *Client) {
		client.RequestBody.StartCrawlDate = startCrawlDate
	}
}

// WithEndCrawlDate sets the end crawl date for the client options.
// If endCrawlDate is specified, results will only include links that
// were crawled before endCrawlDate.
// Must be specified in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)
//
// Parameters:
// - endCrawlDate: the end crawl date to be set.
//
// Returns: a ClientOptions function that updates the endCrawlDate field of the RequestBody struct.
func WithEndCrawlDate(endCrawlDate string) ClientOptions {
	return func(client *Client) {
		client.RequestBody.EndCrawlDate = endCrawlDate
	}
}

// WithStartPublishedDate sets the start published date for the client options.
// If specified, only links with a published date after startPublishedDate will
// be returned.
// Must be specified in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ).
//
// Parameters:
// - startPublishedDate: a string representing the start published date.
//
// Returns: a ClientOptions function that updates the startPublishedDate field of the RequestBody struct.
func WithStartPublishedDate(startPublishedDate string) ClientOptions {
	return func(client *Client) {
		client.RequestBody.StartPublishedDate = startPublishedDate
	}
}

// WithEndPublishedDate sets the end published date for the client options.
// If specified, only links with a published date before endPublishedDate will
// be returned.
// Must be specified in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ).
//
// Parameters:
// - endPublishedDate: the end published date to be set.
//
// Returns: a ClientOptions function that updates the endPublishedDate field of the RequestBody struct.
func WithEndPublishedDate(endPublishedDate string) ClientOptions {
	return func(client *Client) {
		client.RequestBody.EndPublishedDate = endPublishedDate
	}
}


// If ExcludeSourceDomain is true, links from the base domain of the input will be 
// automatically excluded from the results. 
// Default: true
//
// Parameters:
// - excludeSourceDomain: a boolean value indicating whether to exclude the source domain.
//
// Returns: a ClientOptions function that updates the ExcludeSourceDomain field in the RequestBody struct.
func WithExcludeSourceDomain(excludeSourceDomain bool) ClientOptions {
	return func(client *Client) {
		client.RequestBody.ExcludeSourceDomain = excludeSourceDomain
	}
}

// WithAutoprompt sets the value of the UseAutoprompt field in the RequestBody.
// If true, your query will be converted to a Metaphor query. Latency will be much higher.
// Default: false
//
// Parameters:
// - useAutoprompt: a boolean value indicating whether to use autoprompt or not.
//
// Returns: a ClientOptions function that updates the useAutoprompt field of the RequestBody struct.
func WithAutoprompt(useAutoprompt bool) ClientOptions {
	return func(client *Client) {
		client.RequestBody.UseAutoprompt = useAutoprompt
	}
}

// WithType sets the search type for the client.
// Type of search, 'keyword' or 'neural'.
// Default: neural
//
// Parameters:
// - searchType: the type of search to be performed.
//
// Returns: a ClientOptions function that updates the type field of the RequestBody struct.
func WithType(searchType string) ClientOptions {
	return func(client *Client) {
		client.RequestBody.Type = searchType
	}
}

// WithBaseURL sets the base api URK type for the client.
// Default: "https://api.metaphor.systems"
//
// Parameters:
// - baseURL: the metaphor api url string .
//
// Returns: a ClientOptions function that updates the baseURL field of the Client struct.
func WithBaseURL(baseURL string) ClientOptions {
	return func(client *Client) {
		client.BaseURL = baseURL
	}
}

// WithRequestOptions sets the request options for the client.
//
// Parameters:
// - reqOptions: The request options to be set of RequestOptions type.
// Returns: a ClientOptions function that updates the RequestBody with additional options.
func WithRequestOptions(reqOptions *RequestOptions) ClientOptions {
	return func(client *Client) {
		if reqOptions.EndCrawlDate != "" {
			client.RequestBody.EndCrawlDate = reqOptions.EndCrawlDate
		}

		if reqOptions.EndPublishedDate != "" {
			client.RequestBody.EndPublishedDate = reqOptions.EndPublishedDate
		}

		if reqOptions.StartCrawlDate != "" {
			client.RequestBody.StartCrawlDate = reqOptions.StartCrawlDate
		}

		if reqOptions.StartPublishedDate != "" {
			client.RequestBody.StartPublishedDate = reqOptions.StartPublishedDate
		}

		if reqOptions.ExcludeSourceDomain {
			client.RequestBody.ExcludeSourceDomain = reqOptions.ExcludeSourceDomain
		}

		if reqOptions.UseAutoprompt {
			client.RequestBody.UseAutoprompt = reqOptions.UseAutoprompt
		}

		if reqOptions.Type != "" {
			client.RequestBody.Type = reqOptions.Type
		}

		if reqOptions.NumResults != 0 {
			client.RequestBody.NumResults = reqOptions.NumResults
		}

		if reqOptions.ExcludeDomains != nil {
			client.RequestBody.ExcludeDomains = reqOptions.ExcludeDomains
		}

		if reqOptions.IncludeDomains != nil {
			client.RequestBody.IncludeDomains = reqOptions.IncludeDomains
		}
	}
}