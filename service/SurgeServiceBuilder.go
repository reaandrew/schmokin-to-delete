package service

import (
	"sync"

	"github.com/rcrowley/go-metrics"
	surgeHTTP "github.com/reaandrew/surge/infrastructure/http"
	"github.com/reaandrew/surge/utils"
)

type SurgeServiceBuilder struct {
	service *SurgeService
}

func NewSurgeServiceBuilder() *SurgeServiceBuilder {
	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	m := metrics.NewMeter()
	sc := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	c := metrics.NewHistogram(sc)
	co := metrics.NewCounter()
	sendRate := metrics.NewMeter()
	receiveRate := metrics.NewMeter()

	return &SurgeServiceBuilder{
		service: &SurgeService{
			workerCount:        1,
			iterations:         1,
			httpClient:         surgeHTTP.NewDefaultHttpClient(),
			timer:              &utils.DefaultTimer{},
			lock:               sync.Mutex{},
			waitGroup:          sync.WaitGroup{},
			responseTime:       h,
			transactionRate:    m,
			concurrencyCounter: co,
			concurrencyRate:    c,
			dataSendRate:       sendRate,
			dataReceiveRate:    receiveRate,
		},
	}
}
func (builder *SurgeServiceBuilder) SetProcesses(value int) *SurgeServiceBuilder {
	builder.service.processes = value
	return builder
}

func (builder *SurgeServiceBuilder) SetServerHost(value string) *SurgeServiceBuilder {
	builder.service.serverHost = value
	return builder
}

func (builder *SurgeServiceBuilder) SetServerPort(value int) *SurgeServiceBuilder {
	builder.service.serverPort = value
	return builder
}

func (builder *SurgeServiceBuilder) SetServer(value bool) *SurgeServiceBuilder {
	builder.service.server = value
	return builder
}

func (builder *SurgeServiceBuilder) SetWorkers(count int) *SurgeServiceBuilder {
	builder.service.workerCount = count
	return builder
}

func (builder *SurgeServiceBuilder) SetIterations(count int) *SurgeServiceBuilder {
	builder.service.iterations = count
	return builder
}

func (builder *SurgeServiceBuilder) SetRandom(value bool) *SurgeServiceBuilder {
	builder.service.random = value
	return builder
}

func (builder *SurgeServiceBuilder) SetHTTPClient(client surgeHTTP.HttpClient) *SurgeServiceBuilder {
	builder.service.httpClient = client
	return builder
}

func (builder *SurgeServiceBuilder) SetTimer(timer utils.Timer) *SurgeServiceBuilder {
	builder.service.timer = timer
	return builder
}

func (builder *SurgeServiceBuilder) Build() *SurgeService {
	return builder.service
}
