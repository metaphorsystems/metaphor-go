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

	//
	// Extract Document contents for the current search query
	//
	resultContents, err := searchResults.GetContents(ctx, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Search results with extracted contents:")
	formatResults(resultContents)

	//
	// Extract the Document contents for an array of id's
	//
	ids := make([]string, 0)
	for _, result := range searchResults.Results {
		ids = append(ids, result.ID)
	}

	contents, err := client.GetContents(ctx, ids)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Extracted contents for the provided array of strings:")
	formatResults(contents)
}

func formatResults(response *metaphor.ContentsResponse) {
	formattedResults := ""

	for _, result := range response.Contents {
		formattedResults += fmt.Sprintf("Title: %s\nURL: %s\nID: %s\n Content: %s\n\n", result.Title, result.URL, result.ID, result.Extract)
	}

	println(formattedResults)
}