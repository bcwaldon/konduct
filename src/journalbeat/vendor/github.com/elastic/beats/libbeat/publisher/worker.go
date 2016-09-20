/*
Copyright 2016 Planet Labs 

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package publisher

import (
	"expvar"
	"sync"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
)

// Metrics that can retrieved through the expvar web interface.
var (
	messagesInWorkerQueues = expvar.NewInt("libbeatMessagesInWorkerQueues")
)

type worker interface {
	send(m message)
}

type messageWorker struct {
	queue     chan message
	bulkQueue chan message
	ws        *workerSignal
	handler   messageHandler
}

type workerSignal struct {
	done chan struct{}
	wg   sync.WaitGroup
}

type message struct {
	client  *client
	context Context
	event   common.MapStr
	events  []common.MapStr
}

type messageHandler interface {
	onMessage(m message)
	onStop()
}

func newMessageWorker(ws *workerSignal, hwm, bulkHWM int, h messageHandler) *messageWorker {
	p := &messageWorker{}
	p.init(ws, hwm, bulkHWM, h)
	return p
}

func (p *messageWorker) init(ws *workerSignal, hwm, bulkHWM int, h messageHandler) {
	p.queue = make(chan message, hwm)
	p.bulkQueue = make(chan message, bulkHWM)
	p.ws = ws
	p.handler = h

	ws.wg.Add(1)
	go p.run()
}

func (p *messageWorker) run() {
	defer p.shutdown()
	for {
		select {
		case <-p.ws.done:
			return
		case m := <-p.queue:
			p.onEvent(m)
		case m := <-p.bulkQueue:
			p.onEvent(m)
		}
	}
}

func (p *messageWorker) shutdown() {
	p.handler.onStop()
	stopQueue(p.queue)
	stopQueue(p.bulkQueue)
	p.ws.wg.Done()
}

func (p *messageWorker) onEvent(m message) {
	messagesInWorkerQueues.Add(-1)
	p.handler.onMessage(m)
}

func (p *messageWorker) send(m message) {
	send(p.queue, p.bulkQueue, m)
}

func (ws *workerSignal) stop() {
	close(ws.done)
	ws.wg.Wait()
}

func newWorkerSignal() *workerSignal {
	w := &workerSignal{}
	w.Init()
	return w
}

func (ws *workerSignal) Init() {
	ws.done = make(chan struct{})
}

func stopQueue(qu chan message) {
	close(qu)
	for msg := range qu { // clear queue and send fail signal
		op.SigFailed(msg.context.Signal, nil)
	}

}

func send(qu, bulkQu chan message, m message) {
	var ch chan message
	if m.event != nil {
		ch = qu
	} else {
		ch = bulkQu
	}

	var done <-chan struct{}
	if m.client != nil {
		done = m.client.canceler.Done()
	}

	select {
	case <-done: // blocks if nil
		// client closed -> signal drop
		// XXX: send Cancel or Fail signal?
		op.SigFailed(m.context.Signal, ErrClientClosed)

	case ch <- m:
		messagesInWorkerQueues.Add(1)
	}
}
