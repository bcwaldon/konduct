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
package outputs

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
)

type Options struct {
	Guaranteed bool
}

type Outputer interface {
	// Publish event
	PublishEvent(sig op.Signaler, opts Options, event common.MapStr) error

	Close() error
}

type TopologyOutputer interface {
	// Register the agent name and its IPs to the topology map
	PublishIPs(name string, localAddrs []string) error

	// Get the agent name with a specific IP from the topology map
	GetNameByIP(ip string) string
}

// BulkOutputer adds BulkPublish to publish batches of events without looping.
// Outputers still might loop on events or use more efficient bulk-apis if present.
type BulkOutputer interface {
	Outputer
	BulkPublish(sig op.Signaler, opts Options, event []common.MapStr) error
}

// Create and initialize the output plugin
type OutputBuilder func(config *common.Config, topologyExpire int) (Outputer, error)

// Functions to be exported by a output plugin
type OutputInterface interface {
	Outputer
	TopologyOutputer
}

type OutputPlugin struct {
	Name   string
	Config *common.Config
	Output Outputer
}

type bulkOutputAdapter struct {
	Outputer
}

var enabledOutputPlugins = make(map[string]OutputBuilder)

func RegisterOutputPlugin(name string, builder OutputBuilder) {
	enabledOutputPlugins[name] = builder
}

func FindOutputPlugin(name string) OutputBuilder {
	return enabledOutputPlugins[name]
}

func InitOutputs(
	beatName string,
	configs map[string]*common.Config,
	topologyExpire int,
) ([]OutputPlugin, error) {
	var plugins []OutputPlugin = nil
	for name, plugin := range enabledOutputPlugins {
		config, exists := configs[name]
		if !exists {
			continue
		}

		if !config.HasField("index") {
			config.SetString("index", -1, beatName)
		}

		output, err := plugin(config, topologyExpire)
		if err != nil {
			logp.Err("failed to initialize %s plugin as output: %s", name, err)
			return nil, err
		}

		plugin := OutputPlugin{Name: name, Config: config, Output: output}
		plugins = append(plugins, plugin)
		logp.Info("Activated %s as output plugin.", name)
	}
	return plugins, nil
}

// CastBulkOutputer casts out into a BulkOutputer if out implements
// the BulkOutputer interface. If out does not implement the interface an outputer
// wrapper implementing the BulkOutputer interface is returned.
func CastBulkOutputer(out Outputer) BulkOutputer {
	if bo, ok := out.(BulkOutputer); ok {
		return bo
	}
	return &bulkOutputAdapter{out}
}

func (b *bulkOutputAdapter) BulkPublish(
	signal op.Signaler,
	opts Options,
	events []common.MapStr,
) error {
	signal = op.SplitSignaler(signal, len(events))
	for _, evt := range events {
		err := b.PublishEvent(signal, opts, evt)
		if err != nil {
			return err
		}
	}
	return nil
}
