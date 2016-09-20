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
package logstash

import (
	"time"

	"github.com/urso/go-lumber/client/v2"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs/transport"
)

const (
	minWindowSize             int = 1
	defaultStartMaxWindowSize int = 10
)

type client struct {
	*transport.Client
	client *v2.SyncClient
	win    window
}

func newLumberjackClient(
	conn *transport.Client,
	compressLevel int,
	maxWindowSize int,
	timeout time.Duration,
	beat string,
) (*client, error) {
	c := &client{}
	c.Client = conn
	c.win.init(defaultStartMaxWindowSize, maxWindowSize)

	enc, err := makeLogstashEventEncoder(beat)
	if err != nil {
		return nil, err
	}

	cl, err := v2.NewSyncClientWithConn(conn,
		v2.JSONEncoder(enc),
		v2.Timeout(timeout),
		v2.CompressionLevel(compressLevel))
	if err != nil {
		return nil, err
	}

	c.client = cl
	return c, nil
}

func (c *client) Connect(timeout time.Duration) error {
	logp.Debug("logstash", "connect")
	return c.Client.Connect()
}

func (c *client) Close() error {
	logp.Debug("logstash", "close connection")
	return c.Client.Close()
}

func (c *client) PublishEvent(event common.MapStr) error {
	_, err := c.PublishEvents([]common.MapStr{event})
	return err
}

// PublishEvents sends all events to logstash. On error a slice with all events
// not published or confirmed to be processed by logstash will be returned.
func (l *client) PublishEvents(
	events []common.MapStr,
) ([]common.MapStr, error) {
	publishEventsCallCount.Add(1)
	totalNumberOfEvents := len(events)
	for len(events) > 0 {
		n, err := l.publishWindowed(events)

		debug("%v events out of %v events sent to logstash. Continue sending",
			n, len(events))

		events = events[n:]
		if err != nil {
			l.win.shrinkWindow()
			_ = l.Close()

			logp.Err("Failed to publish events caused by: %v", err)

			eventsNotAcked.Add(int64(len(events)))
			ackedEvents.Add(int64(totalNumberOfEvents - len(events)))
			return events, err
		}
	}
	ackedEvents.Add(int64(totalNumberOfEvents))
	return nil, nil
}

// publishWindowed published events with current maximum window size to logstash
// returning the total number of events sent (due to window size, or acks until
// failure).
func (c *client) publishWindowed(events []common.MapStr) (int, error) {
	if len(events) == 0 {
		return 0, nil
	}

	batchSize := len(events)
	windowSize := c.win.get()
	debug("Try to publish %v events to logstash with window size %v",
		batchSize, windowSize)

	// prepare message payload
	if batchSize > windowSize {
		events = events[:windowSize]
	}

	n, err := c.sendEvents(events)
	if err != nil {
		return n, err
	}

	c.win.tryGrowWindow(batchSize)
	return len(events), nil
}

func (c *client) sendEvents(events []common.MapStr) (int, error) {
	if len(events) == 0 {
		return 0, nil
	}

	window := make([]interface{}, len(events))
	for i, event := range events {
		window[i] = event
	}
	return c.client.Send(window)
}
