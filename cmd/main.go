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
