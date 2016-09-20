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
package filter

import (
	"fmt"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type FilterRule interface {
	Filter(event common.MapStr) (common.MapStr, error)
	String() string
}

type FilterConstructor func(config common.Config) (FilterRule, error)

var filterConstructors = map[string]FilterConstructor{}

func RegisterPlugin(name string, constructor FilterConstructor) error {

	logp.Debug("filter", "Register plugin %s", name)

	if _, exists := filterConstructors[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}
	filterConstructors[name] = constructor
	return nil
}
