Milestones (Choose five features from each)

- Siege
  - Specify headers
  - Specify a body
  - Verbose Mode
  - Timeout
  - Supply top level headers which apply to every run
- Locust
  - Connect to existing set of endpoints
  - Config Mode e.g. Node Discovery
  - Suport not validating certs
- K6 - External tool integrations (Kibana, Grafana, Prrometheus) - Metrics
- Jmeter
  - Assertions
  - Extractors (JavaScript, Lua)
  - Ramp Up
  - Specify a Cookie File for persistence and retreival
  - Logic Controllers 
- Gattling
  - Recorder
  - Generators (JavaScript, Lua)
  - Reporting (Generated Artefact)
- Load Runner
  - UI interactive
  - Support specifying a desired transactions per second
- AB
- Wrk

Notes

  - Move to fasthttp
  - Use the reuseport package
  - https://github.com/dop251/goja

Tickets
  - Error vs. Failure
  - Use of Environment Variables
  - Use of URL to determine tech e.g. https, amqp
  - Support AMQP Executor
  - Support MongoDB Executor
  - Support Redis Executor
  - Support MySQL Executor
  - Support PSQL Executor
  - Control Granularity
  - Reporting Mode

Validation
   - URL FILE is at least required OR something else when the features allow it.
