/*
	Asynch log writing helper
*/

package golog

import (
	"container/list"
	"sync"
)

// queueManager is responsible for handling asynch logging
type queueManager struct {
	queue          *list.List // queue of messages to process for logging
	                          // prefer list to array since array memory is never returned
	isInitialized  bool       // if true, an instance of this structure has been initialized and is ready for use
	isStarted      bool       // if true, the queue manager is already running
	mux            sync.Mutex // used to lock the queue to prevent double reads
	shouldShutDown bool       // if true, stop the queueManager since logger is shutting down
}

// enqueue adds a new log message to the message queue
func (mgr *queueManager) enqueue(loggingMessage logMessage) {
	if !mgr.isInitialized {
		panic("Queue manager is uninitalized. Initalize before use.")
	}

	if mgr.shouldShutDown {
		return
	}

	mgr.mux.Lock()
	mgr.queue.PushBack(loggingMessage)
	mgr.mux.Unlock()
}

func (mgr *queueManager) start() {
	if !mgr.isInitialized {
		panic("Queue manager is uninitalized. Initalize before use.")
	}

	if mgr.shouldShutDown || mgr.isStarted {
		return
	}
	mgr.isStarted = true

	go mgr.processMessages()
}

func (mgr *queueManager) stop() {
	if !mgr.isInitialized {
		panic("Queue manager is uninitalized. Initalize before use.")
	}

	if mgr.shouldShutDown {
		return
	}
	mgr.shouldShutDown = true

	mgr.mux.Lock()
	for mgr.queue.Len() > 0 {
		node := mgr.queue.Front()
		mgr.queue.Remove(node)

		nodeValue := node.Value
		loggingMessage := nodeValue.(logMessage)
		writeLog(loggingMessage)
	}
	mgr.mux.Unlock()

	mgr.isInitialized = false
}

// processMessages takes messages off the queue and outputs them
func (mgr *queueManager) processMessages() {
	if !mgr.isInitialized {
		panic("Queue manager is uninitalized. Initalize before use.")
	}

	for {
		if mgr.shouldShutDown {
			return
		}

		mgr.mux.Lock()
		for mgr.queue.Len() > 0 {
			node := mgr.queue.Front()
			mgr.queue.Remove(node)

			nodeValue := node.Value
			loggingMessage := nodeValue.(logMessage)
			writeLog(loggingMessage)
		}
		mgr.mux.Unlock()
	}
}

func createQueueMgr() queueManager {
	return queueManager{list.New(), true, false, sync.Mutex{}, false}
}
