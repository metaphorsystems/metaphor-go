package metaphor

type ClientOptions func(*MetaphorClient)

func WithNumResults(numResults int) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.NumResults = numResults
	}
}

func WithIncludeDomains(includeDomains []string) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.IncludeDomains = includeDomains
	}
}

func WithExcludeDomains(excludeDomains []string) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.ExcludeDomains = excludeDomains
	}
}

func WithStartCrawlDate(startCrawlDate string) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.StartCrawlDate = startCrawlDate
	}
}

func WithEndCrawlDate(endCrawlDate string) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.EndCrawlDate = endCrawlDate
	}
}

func WithStartPublishedDate(startPublishedDate string) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.StartPublishedDate = startPublishedDate
	}
}

func WithEndPublishedDate(endPublishedDate string) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.EndPublishedDate = endPublishedDate
	}
}

func WithAutoprompt(useAutoprompt bool) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.UseAutoprompt = useAutoprompt
	}
}

func WithType(type_ string) ClientOptions {
	return func(client *MetaphorClient) {
		client.RequestBody.Type = type_
	}
}
