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
	
	searchResults, err := client.Search(
		ctx, 
		searchQuery, 
		metaphor.WithAutoprompt(true),
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