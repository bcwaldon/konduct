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
package console

import (
	"encoding/json"
	"os"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
)

func init() {
	outputs.RegisterOutputPlugin("console", New)
}

type console struct {
	config config
}

func New(config *common.Config, _ int) (outputs.Outputer, error) {
	c := &console{config: defaultConfig}
	err := config.Unpack(&c.config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newConsole(pretty bool) *console {
	return &console{config{pretty}}
}

func writeBuffer(buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := os.Stdout.Write(buf[written:])
		if err != nil {
			return err
		}

		written += n
	}
	return nil
}

// Implement Outputer
func (c *console) Close() error {
	return nil
}

func (c *console) PublishEvent(
	s op.Signaler,
	opts outputs.Options,
	event common.MapStr,
) error {
	var jsonEvent []byte
	var err error

	if c.config.Pretty {
		jsonEvent, err = json.MarshalIndent(event, "", "  ")
	} else {
		jsonEvent, err = json.Marshal(event)
	}
	if err != nil {
		logp.Err("Fail to convert the event to JSON (%v): %#v", err, event)
		op.SigCompleted(s)
		return err
	}

	if err = writeBuffer(jsonEvent); err != nil {
		goto fail
	}
	if err = writeBuffer([]byte{'\n'}); err != nil {
		goto fail
	}

	op.SigCompleted(s)
	return nil
fail:
	if opts.Guaranteed {
		logp.Critical("Unable to publish events to console: %v", err)
	}
	op.SigFailed(s, err)
	return err
}
