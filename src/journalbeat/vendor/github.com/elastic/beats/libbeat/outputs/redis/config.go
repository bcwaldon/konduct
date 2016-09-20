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
package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/transport"
)

type redisConfig struct {
	Password    string                `config:"password"`
	Index       string                `config:"index"`
	Port        int                   `config:"port"`
	LoadBalance bool                  `config:"loadbalance"`
	Timeout     time.Duration         `config:"timeout"`
	MaxRetries  int                   `config:"max_retries"`
	TLS         *outputs.TLSConfig    `config:"tls"`
	Proxy       transport.ProxyConfig `config:",inline"`

	Db       int    `config:"db"`
	DataType string `config:"datatype"`

	HostTopology     string `config:"host_topology"`
	PasswordTopology string `config:"password_topology"`
	DbTopology       int    `config:"db_topology"`
}

var (
	defaultConfig = redisConfig{
		Port:             6379,
		LoadBalance:      true,
		Timeout:          5 * time.Second,
		MaxRetries:       3,
		TLS:              nil,
		Db:               0,
		DataType:         "list",
		HostTopology:     "",
		PasswordTopology: "",
		DbTopology:       1,
	}
)

func (c *redisConfig) Validate() error {
	switch c.DataType {
	case "", "list", "channel":
	default:
		return fmt.Errorf("redis data type %v not supported", c.DataType)
	}

	if c.Index == "" {
		return errors.New("index required")
	}

	return nil
}
