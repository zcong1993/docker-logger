package event

type LogEvent struct {
	ContainerId   string
	ContainerName string
	Labels        map[string]string
	Log           string
	LogLevel      string
}

type EventStream struct {
	ch            chan *LogEvent
	ContainerId   string
	ContainerName string
	Labels        map[string]string
	LogLevel      string
}

func NewEventStream(ContainerId string, ContainerName string, Labels map[string]string, LogLevel string) *EventStream {
	ch := make(chan *LogEvent)
	return &EventStream{
		ContainerId:   ContainerId,
		ContainerName: ContainerName,
		Labels:        Labels,
		LogLevel:      LogLevel,
		ch:            ch,
	}
}

func (es *EventStream) GetChannel() chan *LogEvent {
	return es.ch
}

func (es *EventStream) newLogEvent(b []byte) *LogEvent {
	return &LogEvent{
		ContainerId:   es.ContainerId,
		ContainerName: es.ContainerName,
		Labels:        es.Labels,
		LogLevel:      es.LogLevel,
		Log:           string(b),
	}
}

func (es *EventStream) Write(b []byte) (int, error) {
	es.ch <- es.newLogEvent(b)
	return len(b), nil
}

func (es *EventStream) Close() error {
	close(es.ch)
	return nil
}
