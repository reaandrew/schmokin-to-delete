Notes

- Move to fasthttp
- Use the reuseport package

Tickets

- Assertions
- Error vs. Failure
- Use of Environment Variables
- Use of URL to determine tech e.g. https, amqp
- Specify headers
- Specify a body
- Specify a Cookie File for persistence and retreival
- Supply top level headers which apply to every run
- Suport not validating certs
- Support AMQP
- Support PSQL
- JavaScript
- Lua (Might be easiest to start with)
- Recorder
- Reporting (Generated Artefact)
- Swarm
- Server Mode
- UI interactive
- Ramp Up
- Timeout
- Control Granularity
- Logic Controllers
- Generators (JavaScript, Lua)
- Extractors (JavaScript, Lua)
- Metrics
- Support specifying a desired transactions per second
- Verbose Mode
- Reporting Mode
- External tool integrations (Kibana, Grafana, Prrometheus)

Validation
 - URL FILE is at least required OR something else when the features allow it.

# ERRORS:

## TCP Connection Re-use

The following was fixed by 1) reading the entire body of the response and 2) closing the body after use.
```shell
Error: Get http://localhost:8000: dial tcp: lookup localhost: device or resource busy                                                                                             
Usage:                                                                                                                                                                            
   [flags]

Flags:
  -h, --help          help for this command
  -X, --verb string    (default "GET")

Error: Get http://localhost:8000: dial tcp: lookup localhost: device or resource busy                                                                                             
Usage:                                                                                                                                                                            
   [flags]

Flags:
  -h, --help          help for this command
  -X, --verb string    (default "GET")

Transactions: 5000
Availability: 100%
Elapsed Time: 3.290351086s
```




