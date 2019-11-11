package cli

type SurgeCLIBuilder struct {
	cli *SurgeCLI
}

func NewSurgeCLIBuilder() *SurgeCLIBuilder {
	return &SurgeCLIBuilder{}
}

func (builder *SurgeCLIBuilder) SetWorkers(count int) *SurgeCLIBuilder {
	builder.cli.workerCount = count
	return builder
}

func (builder *SurgeCLIBuilder) SetIterations(count int) *SurgeCLIBuilder {
	builder.cli.iterations = count
	return builder
}

func (builder *SurgeCLIBuilder) SetRandom(value bool) *SurgeCLIBuilder {
	builder.cli.random = value
	return builder
}

func (builder *SurgeCLIBuilder) SetURLFilePath(value string) *SurgeCLIBuilder {
	builder.cli.urlFilePath = value
	return builder
}

func (builder *SurgeCLIBuilder) SetProcesses(value int) *SurgeCLIBuilder {
	builder.cli.processes = value
	return builder
}

func (builder *SurgeCLIBuilder) SetServerHost(value string) *SurgeCLIBuilder {
	builder.cli.serverHost = value
	return builder
}

func (builder *SurgeCLIBuilder) SetServerPort(value int) *SurgeCLIBuilder {
	builder.cli.serverPort = value
	return builder
}

func (builder *SurgeCLIBuilder) SetServer(value bool) *SurgeCLIBuilder {
	builder.cli.server = value
	return builder
}

func (builder *SurgeCLIBuilder) Build() *SurgeCLI {
	return builder.cli
}
