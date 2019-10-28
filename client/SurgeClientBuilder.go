package client

import (
	"sync"

	"github.com/reaandrew/surge/utils"
)

type SurgeClientBuilder struct {
	client *surge
}

func NewSurgeClientBuilder() *SurgeClientBuilder {
	return &SurgeClientBuilder{
		client: &surge{
			workerCount: 1,
			iterations:  1,
			httpClient:  NewDefaultHttpClient(),
			timer:       &utils.DefaultTimer{},
			lock:        sync.Mutex{},
			waitGroup:   sync.WaitGroup{},
		},
	}
}

func (builder *SurgeClientBuilder) SetWorkers(count int) *SurgeClientBuilder {
	builder.client.workerCount = count
	return builder
}

func (builder *SurgeClientBuilder) SetIterations(count int) *SurgeClientBuilder {
	builder.client.iterations = count
	return builder
}

func (builder *SurgeClientBuilder) SetRandom(value bool) *SurgeClientBuilder {
	builder.client.random = value
	return builder
}

func (builder *SurgeClientBuilder) SetURLFilePath(value string) *SurgeClientBuilder {
	builder.client.urlFilePath = value
	return builder
}

func (builder *SurgeClientBuilder) SetHTTPClient(client HttpClient) *SurgeClientBuilder {
	builder.client.httpClient = client
	return builder
}

func (builder *SurgeClientBuilder) SetTimer(timer utils.Timer) *SurgeClientBuilder {
	builder.client.timer = timer
	return builder
}

func (builder *SurgeClientBuilder) Build() *surge {
	return builder.client
}