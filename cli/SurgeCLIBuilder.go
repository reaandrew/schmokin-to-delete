package cli

type SchmokinCLIBuilder struct {
	cli *SchmokinCLI
}

func NewSchmokinCLIBuilder() *SchmokinCLIBuilder {
	return &SchmokinCLIBuilder{
		cli: &SchmokinCLI{
			workerCount: 1,
			iterations:  1,
			processes:   1,
			serverHost:  "localhost",
			serverPort:  54321,
			server:      false,
		},
	}
}

func (builder *SchmokinCLIBuilder) SetWorkers(count int) *SchmokinCLIBuilder {
	builder.cli.workerCount = count
	return builder
}

func (builder *SchmokinCLIBuilder) SetIterations(count int) *SchmokinCLIBuilder {
	builder.cli.iterations = count
	return builder
}

func (builder *SchmokinCLIBuilder) SetRandom(value bool) *SchmokinCLIBuilder {
	builder.cli.random = value
	return builder
}

func (builder *SchmokinCLIBuilder) SetURLFilePath(value string) *SchmokinCLIBuilder {
	builder.cli.urlFilePath = value
	return builder
}

func (builder *SchmokinCLIBuilder) SetProcesses(value int) *SchmokinCLIBuilder {
	builder.cli.processes = value
	return builder
}

func (builder *SchmokinCLIBuilder) SetServerHost(value string) *SchmokinCLIBuilder {
	builder.cli.serverHost = value
	return builder
}

func (builder *SchmokinCLIBuilder) SetServerPort(value int) *SchmokinCLIBuilder {
	builder.cli.serverPort = value
	return builder
}

func (builder *SchmokinCLIBuilder) SetServer(value bool) *SchmokinCLIBuilder {
	builder.cli.server = value
	return builder
}

func (builder *SchmokinCLIBuilder) Build() *SchmokinCLI {
	return builder.cli
}
