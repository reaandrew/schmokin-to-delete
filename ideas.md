Tickets

- Assertions
- Error vs. Failure
- Use of Environment Variables
- Use of URL to determine tech e.g. https, amqp
- Specify headers
- Specify a body
- Specify a Cookie File for persistence and retreival
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

ERRORS:

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

## When no server is listening

```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x40 pc=0x7650dc]

goroutine 19 [running]:
github.com/reaandrew/surge/client.HttpCommand.Execute.func1(0xc000122280, 0xc0000601c0, 0x1, 0x1, 0x0, 0x0)
        /home/parallels/go/src/github.com/reaandrew/surge/client/HttpCommand.go:29 +0xfc
github.com/reaandrew/surge/vendor/github.com/spf13/cobra.(*Command).execute(0xc000122280, 0xc000060160, 0x1, 0x1, 0xc000122280, 0xc000060160)
        /home/parallels/go/src/github.com/reaandrew/surge/vendor/github.com/spf13/cobra/command.go:826 +0x460
github.com/reaandrew/surge/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc000122280, 0xc000060170, 0x90e6d3, 0x4)
        /home/parallels/go/src/github.com/reaandrew/surge/vendor/github.com/spf13/cobra/command.go:914 +0x2fb
github.com/reaandrew/surge/vendor/github.com/spf13/cobra.(*Command).Execute(...)
        /home/parallels/go/src/github.com/reaandrew/surge/vendor/github.com/spf13/cobra/command.go:864
github.com/reaandrew/surge/client.HttpCommand.Execute(0x9b5300, 0xccf700, 0xc000060160, 0x1, 0x1, 0x0, 0x0)
        /home/parallels/go/src/github.com/reaandrew/surge/client/HttpCommand.go:42 +0x1bc
github.com/reaandrew/surge/client.(*surge).worker(0xc00011acb0, 0xc0000992c0, 0x1, 0x1)
        /home/parallels/go/src/github.com/reaandrew/surge/client/SurgeClient.go:36 +0xbe
created by github.com/reaandrew/surge/client.(*surge).execute
        /home/parallels/go/src/github.com/reaandrew/surge/client/SurgeClient.go:54 +0x97
```
