package client

import "sync"

type SurgeClientBuilder struct {
	client *Surge
}

func NewSurgeClientBuilder() *SurgeClientBuilder {
	return &SurgeClientBuilder{
		client: &Surge{
			WorkerCount: 1,
			Iterations:  1,
			HttpClient:  NewDefaultHttpClient(),
			lock:        sync.Mutex{},
			waitGroup:   sync.WaitGroup{},
		},
	}
}

func (builder *SurgeClientBuilder) SetWorkers(count int) *SurgeClientBuilder {
	builder.client.WorkerCount = count
	return builder
}

func (builder *SurgeClientBuilder) SetIterations(count int) *SurgeClientBuilder {
	builder.client.Iterations = count
	return builder
}

func (builder *SurgeClientBuilder) SetRandom(value bool) *SurgeClientBuilder {
	builder.client.Random = value
	return builder
}

func (builder *SurgeClientBuilder) SetURLFilePath(value string) *SurgeClientBuilder {
	builder.client.UrlFilePath = value
	return builder
}

func (builder *SurgeClientBuilder) SetHTTPClient(client HttpClient) *SurgeClientBuilder {
	builder.client.HttpClient = client
	return builder
}

func (builder *SurgeClientBuilder) Build() *Surge {
	return builder.client
}
