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

	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/transport"
)

type logstashConfig struct {
	Index            string                `config:"index"`
	Port             int                   `config:"port"`
	LoadBalance      bool                  `config:"loadbalance"`
	BulkMaxSize      int                   `config:"bulk_max_size"`
	Timeout          time.Duration         `config:"timeout"`
	Pipelining       int                   `config:"pipelining"        validate:"min=0"`
	CompressionLevel int                   `config:"compression_level" validate:"min=0, max=9"`
	MaxRetries       int                   `config:"max_retries"       validate:"min=-1"`
	TLS              *outputs.TLSConfig    `config:"tls"`
	Proxy            transport.ProxyConfig `config:",inline"`
}

var (
	defaultConfig = logstashConfig{
		Port:             10200,
		LoadBalance:      false,
		BulkMaxSize:      2048,
		CompressionLevel: 3,
		Timeout:          30 * time.Second,
		MaxRetries:       3,
	}
)
