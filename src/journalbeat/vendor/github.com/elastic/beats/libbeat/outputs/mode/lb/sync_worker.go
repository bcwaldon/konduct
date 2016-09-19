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
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs/mode"
)

type syncWorkerFactory struct {
	clients                 []mode.ProtocolClient
	waitRetry, maxWaitRetry time.Duration
}

// worker instances handle one load-balanced output per instance. Workers receive
// messages from context and return failed send attempts back to the context.
// Client connection state is fully handled by the worker.
type syncWorker struct {
	id      int
	client  mode.ProtocolClient
	backoff *common.Backoff
	ctx     context
}

func SyncClients(
	clients []mode.ProtocolClient,
	waitRetry, maxWaitRetry time.Duration,
) WorkerFactory {
	return &syncWorkerFactory{
		clients:      clients,
		waitRetry:    waitRetry,
		maxWaitRetry: maxWaitRetry,
	}
}

func (s *syncWorkerFactory) count() int { return len(s.clients) }

func (s *syncWorkerFactory) mk(ctx context) ([]worker, error) {
	workers := make([]worker, len(s.clients))
	for i, client := range s.clients {
		workers[i] = newSyncWorker(i, client, ctx, s.waitRetry, s.maxWaitRetry)
	}
	return workers, nil
}

func newSyncWorker(
	id int,
	client mode.ProtocolClient,
	ctx context,
	waitRetry, maxWaitRetry time.Duration,
) *syncWorker {
	return &syncWorker{
		id:      id,
		client:  client,
		backoff: common.NewBackoff(ctx.done, waitRetry, maxWaitRetry),
		ctx:     ctx,
	}
}

func (w *syncWorker) run() {
	client := w.client
	defer func() {
		if client.IsConnected() {
			_ = client.Close()
		}
	}()

	debugf("load balancer: start client loop")
	defer debugf("load balancer: stop client loop")

	done := false
	for !done {
		if done = w.connect(); !done {
			done = w.sendLoop()
		}
		debugf("close client (done=%v)", done)
		client.Close()
	}
}

func (w *syncWorker) connect() bool {
	for {
		debugf("try to (re-)connect client")
		err := w.client.Connect(w.ctx.timeout)
		if !w.backoff.WaitOnError(err) {
			return true
		}

		if err == nil {
			return false
		}

		debugf("connect failed with: %v", err)
	}
}

func (w *syncWorker) sendLoop() (done bool) {
	for {
		msg, ok := w.ctx.receive()
		if !ok {
			return true
		}

		msg.worker = w.id
		err := w.onMessage(msg)
		done = !w.backoff.WaitOnError(err)
		if done || err != nil {
			return done
		}
	}
}

func (w *syncWorker) onMessage(msg eventsMessage) error {
	client := w.client

	if msg.event != nil {
		err := client.PublishEvent(msg.event)
		if err != nil {
			if msg.attemptsLeft > 0 {
				msg.attemptsLeft--
			}
			w.onFail(msg, err)
			return err
		}
	} else {
		events := msg.events
		total := len(events)

		for len(events) > 0 {
			var err error

			events, err = client.PublishEvents(events)
			if err != nil {
				if msg.attemptsLeft > 0 {
					msg.attemptsLeft--
				}

				// reset attempt count if subset of messages has been processed
				if len(events) < total && msg.attemptsLeft >= 0 {
					debugf("reset fails")
					msg.attemptsLeft = w.ctx.maxAttempts
				}

				if err != mode.ErrTempBulkFailure {
					// retry non-published subset of events in batch
					msg.events = events
					w.onFail(msg, err)
					return err
				}

				if w.ctx.maxAttempts > 0 && msg.attemptsLeft == 0 {
					// no more attempts left => drop
					dropping(msg)
					return err
				}

				// reset total count for temporary failure loop
				total = len(events)
			}
		}
	}

	op.SigCompleted(msg.signaler)
	return nil
}

func (w *syncWorker) onFail(msg eventsMessage, err error) {
	logp.Info("Error publishing events (retrying): %s", err)
	w.ctx.pushFailed(msg)
}
