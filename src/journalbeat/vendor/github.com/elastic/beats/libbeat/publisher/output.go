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
	"errors"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
)

type outputWorker struct {
	messageWorker
	out         outputs.BulkOutputer
	config      outputConfig
	maxBulkSize int
}

type outputConfig struct {
	BulkMaxSize   int           `config:"bulk_max_size"`
	FlushInterval time.Duration `config:"flush_interval"`
}

var (
	defaultConfig = outputConfig{
		FlushInterval: 1 * time.Second,
		BulkMaxSize:   2048,
	}
)

var (
	errSendFailed = errors.New("failed send attempt")
)

func newOutputWorker(
	cfg *common.Config,
	out outputs.Outputer,
	ws *workerSignal,
	hwm int,
	bulkHWM int,
) *outputWorker {
	config := defaultConfig
	err := cfg.Unpack(&config)
	if err != nil {
		logp.Err("Failed to read output worker config: %v", err)
		return nil
	}

	o := &outputWorker{
		out:         outputs.CastBulkOutputer(out),
		config:      config,
		maxBulkSize: config.BulkMaxSize,
	}
	o.messageWorker.init(ws, hwm, bulkHWM, o)
	return o
}

func (o *outputWorker) onStop() {
	err := o.out.Close()
	if err != nil {
		logp.Info("Failed to close outputer: %s", err)
	}
}

func (o *outputWorker) onMessage(m message) {
	if m.event != nil {
		o.onEvent(&m.context, m.event)
	} else {
		o.onBulk(&m.context, m.events)
	}
}

func (o *outputWorker) onEvent(ctx *Context, event common.MapStr) {
	debug("output worker: publish single event")
	o.out.PublishEvent(ctx.Signal, outputs.Options{Guaranteed: ctx.Guaranteed}, event)
}

func (o *outputWorker) onBulk(ctx *Context, events []common.MapStr) {
	if len(events) == 0 {
		debug("output worker: no events to publish")
		op.SigCompleted(ctx.Signal)
		return
	}

	if o.maxBulkSize < 0 || len(events) <= o.maxBulkSize {
		o.sendBulk(ctx, events)
		return
	}

	// start splitting bulk request
	splits := (len(events) + (o.maxBulkSize - 1)) / o.maxBulkSize
	ctx.Signal = op.SplitSignaler(ctx.Signal, splits)
	for len(events) > 0 {
		sz := o.maxBulkSize
		if sz > len(events) {
			sz = len(events)
		}
		o.sendBulk(ctx, events[:sz])
		events = events[sz:]
	}
}

func (o *outputWorker) sendBulk(
	ctx *Context,
	events []common.MapStr,
) {
	debug("output worker: publish %v events", len(events))

	opts := outputs.Options{Guaranteed: ctx.Guaranteed}
	err := o.out.BulkPublish(ctx.Signal, opts, events)
	if err != nil {
		logp.Info("Error bulk publishing events: %s", err)
	}
}
