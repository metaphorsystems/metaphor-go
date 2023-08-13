package metaphor

import "context"

type SearchResponse struct {
	Results []struct {
		Id            string  `json:"id"`
		Url           string  `json:"url"`
		Title         string  `json:"title"`
		PublishedDate string  `json:"publishedDate"`
		Author        string  `json:"author"`
		Score         float64 `json:"score"`
		Extract       string
	} `json:"results"`
}

type ContentsResponse struct {
	Contents []struct {
		Url     string `json:"url"`
		Title   string `json:"title"`
		Id      string `json:"id"`
		Extract string `json:"extract"`
	} `json:"contents"`
}

type ErrorResponse struct {
	Text string `json:"error"`
}

// GetContents retrieves contents for the latest search results.
//
// Parameters:
// - ctx: the context.Context for the request.
// - client: The Metaphor client used for the request.
//
// Returns:
// - *ContentsResponse: The contents response object.
// - error: An error if the contents retrieval fails.
func (response SearchResponse) GetContents(ctx context.Context, client *Client) (*ContentsResponse, error) {
	ids := []string{}
	for _, result := range response.Results {
		ids = append(ids, result.Id)
	}
	return client.GetContents(ctx, ids)
}