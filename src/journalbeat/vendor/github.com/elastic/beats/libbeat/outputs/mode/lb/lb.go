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
package lb

import (
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/mode"
)

// LB balances the sending of events between multiple connections.
//
// The balancing algorithm is mostly pull-based, with multiple workers trying to pull
// some amount of work from a shared queue. Workers will try to get a new work item
// only if they have a working/active connection. Workers without active connection
// do not participate until a connection has been re-established.
// Due to the pull based nature the algorithm will load-balance events by random
// with workers having less latencies/turn-around times potentially getting more
// work items then other workers with higher latencies. Thusly the algorithm
// dynamically adapts to resource availability of server events are forwarded to.
//
// Workers not participating in the load-balancing will continuously try to reconnect
// to their configured endpoints. Once a new connection has been established,
// these workers will participate in in load-balancing again.
//
// If a connection becomes unavailable, the events are rescheduled for another
// connection to pick up. Rescheduling events is limited to a maximum number of
// send attempts. If events have not been send after maximum number of allowed
// attemps has been passed, they will be dropped.
//
// Like network connections, distributing events to workers is subject to
// timeout. If no worker is available to pickup a message for sending, the message
// will be dropped internally after max_retries. If mode or message requires
// guaranteed send, message is retried infinitely.
type LB struct {
	ctx context

	// waitGroup + signaling channel for handling shutdown
	wg sync.WaitGroup
}

var (
	debugf = logp.MakeDebug("output")
)

func NewSync(
	clients []mode.ProtocolClient,
	maxAttempts int,
	waitRetry, timeout, maxWaitRetry time.Duration,
) (*LB, error) {
	return New(SyncClients(clients, waitRetry, maxWaitRetry),
		maxAttempts, timeout)
}

func NewAsync(
	clients []mode.AsyncProtocolClient,
	maxAttempts int,
	waitRetry, timeout, maxWaitRetry time.Duration,
) (*LB, error) {
	return New(AsyncClients(clients, waitRetry, maxWaitRetry),
		maxAttempts, timeout)
}

// New create a new load balancer connection mode.
func New(
	makeWorkers WorkerFactory,
	maxAttempts int,
	timeout time.Duration,
) (*LB, error) {
	debugf("configure maxattempts: %v", maxAttempts)

	// maxAttempts signals infinite retry. Convert to -1, so attempts left and
	// and infinite retry can be more easily distinguished by load balancer
	if maxAttempts == 0 {
		maxAttempts = -1
	}

	m := &LB{
		ctx: makeContext(makeWorkers.count(), maxAttempts, timeout),
	}

	if err := m.start(makeWorkers); err != nil {
		return nil, err
	}
	return m, nil
}

// Close stops all workers and closes all open connections. In flight events
// are signaled as failed.
func (m *LB) Close() error {
	m.ctx.Close()
	m.wg.Wait()
	return nil
}

func (m *LB) start(makeWorkers WorkerFactory) error {
	var waitStart sync.WaitGroup
	run := func(w worker) {
		defer m.wg.Done()
		waitStart.Done()
		w.run()
	}

	workers, err := makeWorkers.mk(m.ctx)
	if err != nil {
		return err
	}

	for _, w := range workers {
		m.wg.Add(1)
		waitStart.Add(1)
		go run(w)
	}
	waitStart.Wait()
	return nil
}

// PublishEvent forwards the event to some load balancing worker.
func (m *LB) PublishEvent(
	signaler op.Signaler,
	opts outputs.Options,
	event common.MapStr,
) error {
	return m.publishEventsMessage(opts, eventsMessage{
		worker:   -1,
		signaler: signaler,
		event:    event,
	})
}

// PublishEvents forwards events to some load balancing worker.
func (m *LB) PublishEvents(
	signaler op.Signaler,
	opts outputs.Options,
	events []common.MapStr,
) error {
	return m.publishEventsMessage(opts, eventsMessage{
		worker:   -1,
		signaler: signaler,
		events:   events,
	})
}

func (m *LB) publishEventsMessage(opts outputs.Options, msg eventsMessage) error {
	m.ctx.pushEvents(msg, opts.Guaranteed)
	return nil
}
