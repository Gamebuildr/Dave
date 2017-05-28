package watcher

// Queue is the interface to specify a queue service
type Queue interface {
	ReadNextMessage(url string) (*MessageInfo, error)
	DeleteMessage(messageID string, url string) error
}

// QueueMonitor is the base system for creating unique queue monitors
type QueueMonitor struct {
	Queue Queue
}

// MessageInfo is the representation of a message from a queueing system
type MessageInfo struct {
	MessageID string
	Message   string
}
