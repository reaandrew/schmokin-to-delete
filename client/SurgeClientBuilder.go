package client

import (
	"sync"

	"github.com/rcrowley/go-metrics"
	"github.com/reaandrew/surge/utils"
)

type SurgeClientBuilder struct {
	client *surge
}

func NewSurgeClientBuilder() *SurgeClientBuilder {
	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	m := metrics.NewMeter()
	sc := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	c := metrics.NewHistogram(sc)
	co := metrics.NewCounter()
	sendRate := metrics.NewMeter()
	receiveRate := metrics.NewMeter()

	return &SurgeClientBuilder{
		client: &surge{
			workerCount:        1,
			iterations:         1,
			processes:          1,
			httpClient:         NewDefaultHttpClient(),
			timer:              &utils.DefaultTimer{},
			lock:               sync.Mutex{},
			waitGroup:          sync.WaitGroup{},
			responseTime:       h,
			transactionRate:    m,
			concurrencyCounter: co,
			concurrencyRate:    c,
			dataSendRate:       sendRate,
			dataReceiveRate:    receiveRate,
			server:             false,
			serverPort:         54321,
			serverHost:         "localhost",
		},
	}
}
func (builder *SurgeClientBuilder) SetProcesses(value int) *SurgeClientBuilder {
	builder.client.processes = value
	return builder
}

func (builder *SurgeClientBuilder) SetServerHost(value string) *SurgeClientBuilder {
	builder.client.serverHost = value
	return builder
}

func (builder *SurgeClientBuilder) SetServerPort(value int) *SurgeClientBuilder {
	builder.client.serverPort = value
	return builder
}

func (builder *SurgeClientBuilder) SetServer(value bool) *SurgeClientBuilder {
	builder.client.server = value
	return builder
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
