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

import "github.com/elastic/beats/libbeat/common/op"

type syncPipeline struct {
	pub *Publisher
}

func newSyncPipeline(pub *Publisher, hwm, bulkHWM int) *syncPipeline {
	return &syncPipeline{pub: pub}
}

func (p *syncPipeline) publish(m message) bool {
	if p.pub.disabled {
		debug("publisher disabled")
		op.SigCompleted(m.context.Signal)
		return true
	}

	client := m.client
	signal := m.context.Signal
	sync := op.NewSignalChannel()
	if len(p.pub.Output) > 1 {
		m.context.Signal = op.SplitSignaler(sync, len(p.pub.Output))
	} else {
		m.context.Signal = sync
	}

	for _, o := range p.pub.Output {
		o.send(m)
	}

	// Await completion signal from output plugin. If client has been disconnected
	// ignore any signal and drop events no matter if send or not.
	select {
	case <-client.canceler.Done():
		return true
	case sig := <-sync.C:
		sig.Apply(signal)
		return sig == op.SignalCompleted
	}
}
