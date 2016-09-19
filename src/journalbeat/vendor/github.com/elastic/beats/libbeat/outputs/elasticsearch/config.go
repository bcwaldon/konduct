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
package elasticsearch

import (
	"time"

	"github.com/elastic/beats/libbeat/outputs"
)

type elasticsearchConfig struct {
	Protocol     string             `config:"protocol"`
	Path         string             `config:"path"`
	Params       map[string]string  `config:"parameters"`
	Username     string             `config:"username"`
	Password     string             `config:"password"`
	ProxyURL     string             `config:"proxy_url"`
	Index        string             `config:"index"`
	LoadBalance  bool               `config:"loadbalance"`
	TLS          *outputs.TLSConfig `config:"tls"`
	MaxRetries   int                `config:"max_retries"`
	Timeout      time.Duration      `config:"timeout"`
	SaveTopology bool               `config:"save_topology"`
	Template     Template           `config:"template"`
}

type Template struct {
	Name      string `config:"name"`
	Path      string `config:"path"`
	Overwrite bool   `config:"overwrite"`
}

const (
	defaultBulkSize = 50
)

var (
	defaultConfig = elasticsearchConfig{
		Protocol:    "",
		Path:        "",
		ProxyURL:    "",
		Params:      nil,
		Username:    "",
		Password:    "",
		Timeout:     90 * time.Second,
		MaxRetries:  3,
		TLS:         nil,
		LoadBalance: true,
	}
)

func (c *elasticsearchConfig) Validate() error {
	if c.ProxyURL != "" {
		if _, err := parseProxyURL(c.ProxyURL); err != nil {
			return err
		}
	}

	return nil
}
