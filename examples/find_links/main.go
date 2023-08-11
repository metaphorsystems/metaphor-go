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

	fmt.Println("Search results: ")
	formatResults(searchResults)

	searchLink := searchResults.Results[0].Url
	linkSearchResults, err := client.FindSimilar(
		ctx, 
		searchLink, 
		metaphor.WithNumResults(5),
	)

	fmt.Println("Find similar links fo: ", searchLink)
	fmt.Println("Similar links found: ")
	formatResults(linkSearchResults)

}

func formatResults(response *metaphor.SearchResponse) {
	formattedResults := ""

	for _, result := range response.Results {
		formattedResults += fmt.Sprintf("Title: %s\nURL: %s\nID: %s\n\n", result.Title, result.Url, result.Id)
	}

	println(formattedResults)
}