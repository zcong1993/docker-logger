package event

// LogEvent is struct of log event stream event
type LogEvent struct {
	ContainerId   string
	ContainerName string
	Log           string
	LogLevel      string
}

// EventStream is struct of event stream
type EventStream struct {
	Ch            chan LogEvent
	ContainerId   string
	ContainerName string
	LogLevel      string
}

// NewEventStream is constructor of event stream
func NewEventStream(ContainerId string, ContainerName string, LogLevel string, Ch chan LogEvent) *EventStream {
	return &EventStream{
		ContainerId:   ContainerId,
		ContainerName: ContainerName,
		LogLevel:      LogLevel,
		Ch:            Ch,
	}
}

func (es *EventStream) newLogEvent(b []byte) LogEvent {
	return LogEvent{
		ContainerId:   es.ContainerId,
		ContainerName: es.ContainerName,
		LogLevel:      es.LogLevel,
		Log:           string(b),
	}
}

// Write implement io writer
func (es *EventStream) Write(b []byte) (int, error) {
	es.Ch <- es.newLogEvent(b)
	return len(b), nil
}
