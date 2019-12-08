package service

import (
	"sync"

	"github.com/rcrowley/go-metrics"
	schmokinHTTP "github.com/reaandrew/schmokin/infrastructure/http"
	"github.com/reaandrew/schmokin/utils"
)

type SchmokinServiceBuilder struct {
	service *SchmokinService
}

func NewSchmokinServiceBuilder() *SchmokinServiceBuilder {
	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	m := metrics.NewMeter()
	sc := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	c := metrics.NewHistogram(sc)
	co := metrics.NewCounter()
	sendRate := metrics.NewMeter()
	receiveRate := metrics.NewMeter()

	return &SchmokinServiceBuilder{
		service: &SchmokinService{
			workerCount:        1,
			iterations:         1,
			httpClient:         schmokinHTTP.NewDefaultClient(),
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

func (builder *SchmokinServiceBuilder) SetWorkers(count int) *SchmokinServiceBuilder {
	builder.service.workerCount = count
	return builder
}

func (builder *SchmokinServiceBuilder) SetIterations(count int) *SchmokinServiceBuilder {
	builder.service.iterations = count
	return builder
}

func (builder *SchmokinServiceBuilder) SetRandom(value bool) *SchmokinServiceBuilder {
	builder.service.random = value
	return builder
}

func (builder *SchmokinServiceBuilder) SetClient(client schmokinHTTP.Client) *SchmokinServiceBuilder {
	builder.service.httpClient = client
	return builder
}

func (builder *SchmokinServiceBuilder) SetTimer(timer utils.Timer) *SchmokinServiceBuilder {
	builder.service.timer = timer
	return builder
}

func (builder *SchmokinServiceBuilder) Build() *SchmokinService {
	return builder.service
}
