package watcher

// Queue is the interface to specify a queue service
type Queue interface {
	ReadNextMessage(url string) (int, error)
	Setup()
}

// QueueMonitor is the base system for creating unique queue monitors
type QueueMonitor struct {
	Queue Queue
}
