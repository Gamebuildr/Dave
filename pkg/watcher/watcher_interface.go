package watcher

// Queue is the interface to specify a queue service
type Queue interface {
	ReadQueueMessagesCount(url string) (int, error)
}

// QueueMonitor is the base system for creating unique queue monitors
type QueueMonitor struct {
	Queue Queue
}
