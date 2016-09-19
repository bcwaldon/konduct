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
package sarama

import "sync"

// SyncProducer publishes Kafka messages, blocking until they have been acknowledged. It routes messages to the correct
// broker, refreshing metadata as appropriate, and parses responses for errors. You must call Close() on a producer
// to avoid leaks, it may not be garbage-collected automatically when it passes out of scope.
//
// The SyncProducer comes with two caveats: it will generally be less efficient than the AsyncProducer, and the actual
// durability guarantee provided when a message is acknowledged depend on the configured value of `Producer.RequiredAcks`.
// There are configurations where a message acknowledged by the SyncProducer can still sometimes be lost.
type SyncProducer interface {

	// SendMessage produces a given message, and returns only when it either has
	// succeeded or failed to produce. It will return the partition and the offset
	// of the produced message, or an error if the message failed to produce.
	SendMessage(msg *ProducerMessage) (partition int32, offset int64, err error)

	// Close shuts down the producer and flushes any messages it may have buffered.
	// You must call this function before a producer object passes out of scope, as
	// it may otherwise leak memory. You must call this before calling Close on the
	// underlying client.
	Close() error
}

type syncProducer struct {
	producer *asyncProducer
	wg       sync.WaitGroup
}

// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
func NewSyncProducer(addrs []string, config *Config) (SyncProducer, error) {
	p, err := NewAsyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}
	return newSyncProducerFromAsyncProducer(p.(*asyncProducer)), nil
}

// NewSyncProducerFromClient creates a new SyncProducer using the given client. It is still
// necessary to call Close() on the underlying client when shutting down this producer.
func NewSyncProducerFromClient(client Client) (SyncProducer, error) {
	p, err := NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return newSyncProducerFromAsyncProducer(p.(*asyncProducer)), nil
}

func newSyncProducerFromAsyncProducer(p *asyncProducer) *syncProducer {
	p.conf.Producer.Return.Successes = true
	p.conf.Producer.Return.Errors = true
	sp := &syncProducer{producer: p}

	sp.wg.Add(2)
	go withRecover(sp.handleSuccesses)
	go withRecover(sp.handleErrors)

	return sp
}

func (sp *syncProducer) SendMessage(msg *ProducerMessage) (partition int32, offset int64, err error) {
	oldMetadata := msg.Metadata
	defer func() {
		msg.Metadata = oldMetadata
	}()

	expectation := make(chan error, 1)
	msg.Metadata = expectation
	sp.producer.Input() <- msg

	if err := <-expectation; err != nil {
		return -1, -1, err
	}

	return msg.Partition, msg.Offset, nil
}

func (sp *syncProducer) handleSuccesses() {
	defer sp.wg.Done()
	for msg := range sp.producer.Successes() {
		expectation := msg.Metadata.(chan error)
		expectation <- nil
	}
}

func (sp *syncProducer) handleErrors() {
	defer sp.wg.Done()
	for err := range sp.producer.Errors() {
		expectation := err.Msg.Metadata.(chan error)
		expectation <- err.Err
	}
}

func (sp *syncProducer) Close() error {
	sp.producer.AsyncClose()
	sp.wg.Wait()
	return nil
}
