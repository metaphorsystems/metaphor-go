package main

import (
	"context"
	"fmt"
	"os"

	"github.com/metaphorsystems/metaphor-go"
)

func main() {

	apiKey := os.Getenv("METAPHOR_API_KEY")

	client, err := metaphor.NewClient(apiKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	searchQuery := "Who is RDJ?"
	ctx := context.Background()

	// 
	// Here we can use RequestOptions to customize the search.
	//
	reqOptions := metaphor.RequestOptions {
		StartCrawlDate: "2023-01-01",
		EndCrawlDate: "2023-12-32",
		ExcludeDomains: []string{"www.wikipedia.com"},
	}
	
	//
	// reqOption can be combined with ClientOptions to build the final search request.
	// If the same filed is present in both ClientOptions and RequestOptions, 
	// the last one set will be used.
	// in this example, we are setting UseAutoprompt to true.
	// but if we set in reqOptions, it will override the value set in ClientOptions.
	searchResults, err := client.Search(
		ctx, 
		searchQuery, 
		metaphor.WithAutoprompt(true),
		metaphor.WithRequestOptions(&reqOptions),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	resultContents, err := searchResults.GetContents(ctx, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	formatResults(resultContents)
}

func formatResults(response *metaphor.ContentsResponse) {
	formattedResults := ""

	for _, result := range response.Contents {
		formattedResults += fmt.Sprintf("Title: %s\nURL: %s\nID: %s\n Content: %s\n\n", result.Title, result.URL, result.ID, result.Extract)
	}

	println(formattedResults)
}