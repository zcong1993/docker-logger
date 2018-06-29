package logger

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/zcong1993/docker-logger/event"
	"log"
	"strings"
	"time"
)

// Manager is docker logger manager
type Manager struct {
	client  *docker.Client
	ch      chan event.LogEvent
	Ignores []string
}

// NewManager is constructor for docker logger manager
func NewManager(client *docker.Client, ignores []string) *Manager {
	ch := make(chan event.LogEvent, 200)
	return &Manager{
		client:  client,
		Ignores: ignores,
		ch:      ch,
	}
}

// Start start a docker logger watcher and return event channel
func (m *Manager) Start() chan event.LogEvent {
	containers, err := m.client.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		log.Fatalf("[ERROR] get container error, %v", err)
	}
	for _, c := range containers {
		containerName := getContainerName(c)
		if contains(containerName, m.Ignores) {
			log.Printf("[INFO] container %s excluded", containerName)
			continue
		}
		go m.startWatch(c.ID, containerName)
		if err != nil {
			log.Fatalf("[ERROR] watch log error, %v", err)
		}
	}
	dockerEventsCh := make(chan *docker.APIEvents)
	if err := m.client.AddEventListener(dockerEventsCh); err != nil {
		log.Fatalf("[ERROR] can't add even listener, %v", err)
	}

	upStatuses := []string{"start", "restart"}
	downStatuses := []string{"die", "destroy", "stop", "pause"}

	go func() {
		for dockerEvent := range dockerEventsCh {
			if dockerEvent.Type == "container" {
				if !contains(dockerEvent.Status, upStatuses) && !contains(dockerEvent.Status, downStatuses) {
					continue
				}
				containerName := strings.TrimPrefix(dockerEvent.Actor.Attributes["name"], "/")
				if contains(containerName, m.Ignores) {
					log.Printf("[INFO] container %s excluded", containerName)
					continue
				}
				go m.startWatch(dockerEvent.Actor.ID, containerName)
			}
		}
	}()
	return m.ch
}

func (m *Manager) startWatch(id, name string) {
	out := event.NewEventStream(id, name, "Log", m.ch)
	errEs := event.NewEventStream(id, name, "Err", m.ch)
	logOpts := docker.LogsOptions{
		Context:           context.Background(),
		Container:         id,
		OutputStream:      out,
		ErrorStream:       errEs,
		Tail:              "10",
		Follow:            true,
		Stdout:            true,
		Stderr:            true,
		InactivityTimeout: time.Hour * 10000,
	}
	err := m.client.Logs(logOpts)
	if strings.HasPrefix(err.Error(), "error from daemon in stream: Error grabbing logs: EOF") {
		logOpts.Tail = ""
		err = m.client.Logs(logOpts)
	}
	if err != nil {
		log.Fatalf("watch container log error: %v", err)
	}
}

func getContainerName(container docker.APIContainers) string {
	return strings.TrimPrefix(container.Names[0], "/")
}

func contains(e string, s []string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
