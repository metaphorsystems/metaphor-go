package main

import (
	"context"
	"fmt"
	"os"

	"github.com/metaphorsystems/metaphor-go"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

const (
	SearchAsistantTemplate = `
	You are a helpful assistant that generates search queiries based on user questions. 
	Only generate one search query.
	`

	SummaryAssistantTemplate = `
	You are a helpful assistant that summarizes the content of a webpage. 
	Summarize the users input.
	`
)

func main() {

	apiKey := os.Getenv("METAPHOR_API_KEY")
	openaiApiKey := os.Getenv("OPENAI_API_KEY")

	input := "What's the recent news on physics today?"

	llm, err := openai.NewChat(
		openai.WithToken(openaiApiKey),
		openai.WithModel("gpt-3.5-turbo"),
	)

	llmResponse, err := llm.Call(
		context.Background(),
		[]schema.ChatMessage{
			schema.SystemChatMessage{Content: SearchAsistantTemplate},
			schema.HumanChatMessage{Content: input},
		},
	)

	if err != nil {
		fmt.Println(err)
		return
	}
		
	client, err := metaphor.NewClient(apiKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	searchQuery := llmResponse.Content
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

	llmSummary, err := llm.Call(
		ctx, 
		[]schema.ChatMessage{
			schema.SystemChatMessage{Content: SummaryAssistantTemplate},
			schema.HumanChatMessage{Content: resultContents.Contents[0].Extract},
		},
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Summary for %s: %s",resultContents.Contents[0].Title, llmSummary.Content)

}