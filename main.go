package main

import (
	"context"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/zcong1993/docker-logger/event"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type LogStreamer struct {
	DockerClient  *docker.Client
	ContainerID   string
	ContainerName string
	LogWriter     io.WriteCloser
	ErrWriter     io.WriteCloser
	Context       context.Context
	CancelFn      context.CancelFunc
}

type LogStruct struct {
	Name string
	Log  string
}

type EventWC struct {
	ch   chan *LogStruct
	name string
}

func newEWC(name string) *EventWC {
	ch := make(chan *LogStruct)
	return &EventWC{
		ch:   ch,
		name: name,
	}
}

func (e *EventWC) GetChannal() chan *LogStruct {
	return e.ch
}

func (e *EventWC) Write(b []byte) (int, error) {
	e.ch <- &LogStruct{
		Name: e.name,
		Log:  string(b),
	}
	return len(b), nil
}

func (e *EventWC) Close() error {
	close(e.ch)
	return nil
}

func main() {
	mapping := map[string]docker.APIContainers{}
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	containers, err := client.ListContainers(docker.ListContainersOptions{All: false})
	for _, c := range containers {
		containerName := strings.TrimPrefix(c.Names[0], "/")
		fmt.Printf("%+v\n", containerName)
		mapping[containerName] = c
	}
	containerName := "ticker"
	ct := mapping[containerName]
	ctx, cancel := context.WithCancel(context.Background())
	es := event.NewEventStream(ct.ID, strings.TrimPrefix(ct.Names[0], "/"), ct.Labels, "Info")
	l := &LogStreamer{
		Context:       ctx,
		CancelFn:      cancel,
		ContainerID:   ct.ID,
		DockerClient:  client,
		LogWriter:     es,
		ErrWriter:     os.Stderr,
		ContainerName: containerName,
	}

	logOpts := docker.LogsOptions{
		Context:           l.Context,
		Container:         l.ContainerID,
		OutputStream:      l.LogWriter, // logs writer for stdout
		ErrorStream:       l.ErrWriter, // err writer for stderr
		Tail:              "10",
		Follow:            true,
		Stdout:            true,
		Stderr:            true,
		InactivityTimeout: time.Hour * 10000,
	}
	go func() {
		err = l.DockerClient.Logs(logOpts)
		if strings.HasPrefix(err.Error(), "error from daemon in stream: Error grabbing logs: EOF") {
			logOpts.Tail = ""
			err = l.DockerClient.Logs(logOpts)
		}
	}()
	ch := es.GetChannel()
	for msg := range ch {
		fmt.Printf("out %s %s", msg.ContainerName, msg.Log)
	}
	log.Printf("[INFO] stream from %s terminated, %v", l.ContainerID, err)
}
