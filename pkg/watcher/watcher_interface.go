package watcher

import "github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"

// Queue is the interface to specify a queue service
type Queue interface {
	ReadQueueMessages(url string) (int, error)
}

// Service is the base system for creating unique queue services
type Service struct {
	Queue Queue
	Log   logger.Log
}
