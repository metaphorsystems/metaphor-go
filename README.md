# metaphor-go

Offical metpahor Go client. 

# Installation

Install `metaphor-go` with `go get`:

```bash
go get github.com/metaphorsystems/metaphor-go
```

# Quick Start

Import the package and initialize the client:

```go
import (
  "github.com/metaphorsystems/metaphor-go"
)

func main() {

  client, err := metaphor.NewClient(os.Getenv("METAPHOR_API_KEY"))
  if err != nil {
    fmt.Println(err)
    return
  }
}
``` 

Performs a search on the Metaphor system with the given parameters.

```go

func main() {
  client, err := metaphor.NewClient(os.Getenv("METAPHOR_API_KEY"))
  if err != nil {
    fmt.Println(err)
    return
  }
  
  searchQuery := "Who is RDJ?"
	
  searchResults, err := client.Search(
    context.Background(), 
    searchQuery,
    metaphor.WithNumResults(5),
    metaphor.WithIncludeDomains([]string{"nytimes.com", "wsj.com"}),
    metaphor.WithStartPublishedDate("2023-06-12"),
    metaphor.WithAutoprompt(true),
  )

  if err != nil {
    fmt.Println(err)
    return
  }

  formattedResults := ""  
  
  for _, result := range response.Contents {
    formattedResults += fmt.Sprintf("
      Title: %s\nURL: %s\n\n", 
      result.Title, result.URL,
    )
  } 

  fmt.Println(formattedResults)
}
```

Finds content similar to the specified URL.

```go
  url := "https://waitbutwhy.com/2014/05/fermi-paradox.html"
  
  response, err := client.FindSimilar(
    context.Background(), 
    url,
    metaphor.WithNumResults(5),
  )

  if err != nil {
    fmt.Println(err)
    return
  }

  formattedResults := ""
  
  for _, result := range response.Contents {
    formattedResults += fmt.Sprintf("
      Title: %s\nURL: %s\n\n", 
      result.Title, result.URL,
    )
  } 

  fmt.Println(formattedResults)
```

Retrieves document contents: 

```go
  ids := []string{"8U71IlQ5DUTdsZFherhhYA", "X3wd0PbJmAvhu_DQjDKA7A"}
  
  response, err := client.GetContents(
    context.Background(), 
    url,
  )

  if err != nil {
    fmt.Println(err)
    return
  }

  formattedResults := ""

  for _, result := range response.Contents {
    formattedResults += fmt.Sprintf(
      "Title: %s\nURL: %s\nID: %s\n Content: %s\n\n", 
      result.Title, result.URL, result.ID, result.Extract,
    )
  } 

  fmt.Println(formattedResults)
```


> Detailed examples with full implementations can be found in the [examples](./examples) directory.

# Contributions

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.


# License

MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.



