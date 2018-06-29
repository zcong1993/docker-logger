# docker-logger [![Go Report Card](https://goreportcard.com/badge/github.com/zcong1993/docker-logger)](https://goreportcard.com/report/github.com/zcong1993/docker-logger) [![CircleCI branch](https://img.shields.io/circleci/project/github/zcong1993/docker-logger/master.svg)](https://circleci.com/gh/zcong1993/docker-logger/tree/master)

> docker container logger collector

## Usage

### cli

```bash
$ docker-logger
# custom endpoint
$ docker-logger -endpoint "your endpoint"
# ignore container by name
$ docker-logger -ignore foo
```

### lib

```go
package main

import (
	"flag"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/zcong1993/docker-logger/logger"
)

func main() {
	var (
		endpoint string
		ignore   string
	)
	flag.StringVar(&endpoint, "endpoint", "unix:///var/run/docker.sock", "docker endpoint")
	flag.StringVar(&ignore, "ignore", "", "ignore container")
	flag.Parse()

	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	lm := logger.NewManager(client, []string{ignore})
	ch := lm.Start()

	for ev := range ch {
		fmt.Printf("container: %s - level: %s - %s\n", ev.ContainerName, ev.LogLevel, ev.Log)
	}
}
```


## License

MIT &copy; zcong1993
