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

import (
	"hash"
	"hash/fnv"
	"math/rand"
	"time"
)

// Partitioner is anything that, given a Kafka message and a number of partitions indexed [0...numPartitions-1],
// decides to which partition to send the message. RandomPartitioner, RoundRobinPartitioner and HashPartitioner are provided
// as simple default implementations.
type Partitioner interface {
	// Partition takes a message and partition count and chooses a partition
	Partition(message *ProducerMessage, numPartitions int32) (int32, error)

	// RequiresConsistency indicates to the user of the partitioner whether the
	// mapping of key->partition is consistent or not. Specifically, if a
	// partitioner requires consistency then it must be allowed to choose from all
	// partitions (even ones known to be unavailable), and its choice must be
	// respected by the caller. The obvious example is the HashPartitioner.
	RequiresConsistency() bool
}

// PartitionerConstructor is the type for a function capable of constructing new Partitioners.
type PartitionerConstructor func(topic string) Partitioner

type manualPartitioner struct{}

// NewManualPartitioner returns a Partitioner which uses the partition manually set in the provided
// ProducerMessage's Partition field as the partition to produce to.
func NewManualPartitioner(topic string) Partitioner {
	return new(manualPartitioner)
}

func (p *manualPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) {
	return message.Partition, nil
}

func (p *manualPartitioner) RequiresConsistency() bool {
	return true
}

type randomPartitioner struct {
	generator *rand.Rand
}

// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
func NewRandomPartitioner(topic string) Partitioner {
	p := new(randomPartitioner)
	p.generator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return p
}

func (p *randomPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) {
	return int32(p.generator.Intn(int(numPartitions))), nil
}

func (p *randomPartitioner) RequiresConsistency() bool {
	return false
}

type roundRobinPartitioner struct {
	partition int32
}

// NewRoundRobinPartitioner returns a Partitioner which walks through the available partitions one at a time.
func NewRoundRobinPartitioner(topic string) Partitioner {
	return &roundRobinPartitioner{}
}

func (p *roundRobinPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) {
	if p.partition >= numPartitions {
		p.partition = 0
	}
	ret := p.partition
	p.partition++
	return ret, nil
}

func (p *roundRobinPartitioner) RequiresConsistency() bool {
	return false
}

type hashPartitioner struct {
	random Partitioner
	hasher hash.Hash32
}

// NewHashPartitioner returns a Partitioner which behaves as follows. If the message's key is nil, or fails to
// encode, then a random partition is chosen. Otherwise the FNV-1a hash of the encoded bytes of the message key
// is used, modulus the number of partitions. This ensures that messages with the same key always end up on the
// same partition.
func NewHashPartitioner(topic string) Partitioner {
	p := new(hashPartitioner)
	p.random = NewRandomPartitioner(topic)
	p.hasher = fnv.New32a()
	return p
}

func (p *hashPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) {
	if message.Key == nil {
		return p.random.Partition(message, numPartitions)
	}
	bytes, err := message.Key.Encode()
	if err != nil {
		return -1, err
	}
	p.hasher.Reset()
	_, err = p.hasher.Write(bytes)
	if err != nil {
		return -1, err
	}
	hash := int32(p.hasher.Sum32())
	if hash < 0 {
		hash = -hash
	}
	return hash % numPartitions, nil
}

func (p *hashPartitioner) RequiresConsistency() bool {
	return true
}
