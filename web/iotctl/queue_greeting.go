package iotctl

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// greetingEngine will be a single instance of our greeting queue engine.
// It should keep a map of Node list to which a greeting is need to be sent to,
// a queue channel to communicate with a goroutine and a mutex to synchronize
// access to our map.
type greetingEngine struct {
	queueItems map[string]bool
	queue      chan greetingQueue
	mutex      *sync.Mutex
}

// greetingQueue will be used as a framework struct for passing information
// using a channel.
type greetingQueue struct {
	Network string
	Node    string
}

// GreetingQueueInit will initialize the greeting queue with required
// parameters.
func (app *Iotctl) GreetingQueueInit(size int) {
	app.greetingQueue = &greetingEngine{
		queue:      make(chan greetingQueue, size),
		queueItems: make(map[string]bool),
		mutex:      &sync.Mutex{},
	}
}

func (app *Iotctl) greetingQueueAdd(net, nodeId string, greeting bool) {
	app.greetingQueue.queue <- greetingQueue{
		Network: net,
		Node:    nodeId,
	}
	logrus.Debugf("queued greeting for %s", net)

	app.greetingQueue.mutex.Lock()
	defer app.greetingQueue.mutex.Unlock()
	app.greetingQueue.queueItems[nodeId] = true

	return
}

// greetingQueueLoop will be ran as a seperate goroutine and will
// process the greeting queue - greetings will be passed into the queue
// and sent to the device.
func (app *Iotctl) greetingQueueLoop(die <-chan struct{}) {
	logrus.Debug("starting greetingQueueLoop")
	app.wg.Add(1)

	for {
		select {
		case greeting := <-app.greetingQueue.queue:
			// public a greeting to the device
			// go app.PublishGreeting()

			app.greetingQueue.mutex.Lock()
			app.greetingQueue.queueItems[greeting.Node] = false
			logrus.Debugf("greeting queue published for %s", greeting.Node)
			app.greetingQueue.mutex.Unlock()
		case <-die:
			app.wg.Done()
			return
		}
	}
}
