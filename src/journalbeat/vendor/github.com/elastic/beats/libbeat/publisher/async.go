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
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
)

type asyncPipeline struct {
	outputs []worker
	pub     *Publisher
}

const (
	defaultBulkSize = 2048
)

func newAsyncPipeline(
	pub *Publisher,
	hwm, bulkHWM int,
	ws *workerSignal,
) *asyncPipeline {
	p := &asyncPipeline{pub: pub}

	var outputs []worker
	for _, out := range pub.Output {
		outputs = append(outputs, makeAsyncOutput(ws, hwm, bulkHWM, out))
	}

	p.outputs = outputs
	return p
}

func (p *asyncPipeline) publish(m message) bool {
	if p.pub.disabled {
		debug("publisher disabled")
		op.SigCompleted(m.context.Signal)
		return true
	}

	if m.context.Signal != nil {
		s := op.CancelableSignaler(m.client.canceler, m.context.Signal)
		if len(p.outputs) > 1 {
			s = op.SplitSignaler(s, len(p.outputs))
		}
		m.context.Signal = s
	}

	for _, o := range p.outputs {
		o.send(m)
	}
	return true
}

func makeAsyncOutput(
	ws *workerSignal,
	hwm, bulkHWM int,
	worker *outputWorker,
) worker {
	config := worker.config

	flushInterval := config.FlushInterval
	maxBulkSize := config.BulkMaxSize
	logp.Info("Flush Interval set to: %v", flushInterval)
	logp.Info("Max Bulk Size set to: %v", maxBulkSize)

	// batching disabled
	if flushInterval <= 0 || maxBulkSize <= 0 {
		return worker
	}

	debug("create bulk processing worker (interval=%v, bulk size=%v)",
		flushInterval, maxBulkSize)
	return newBulkWorker(ws, hwm, bulkHWM, worker, flushInterval, maxBulkSize)
}
