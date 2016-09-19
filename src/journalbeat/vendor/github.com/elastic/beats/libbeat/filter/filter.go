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
	"strings"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type Filters struct {
	list []FilterRule
}

func New(config FilterPluginConfig) (*Filters, error) {

	filters := Filters{}

	for _, filter := range config {

		if len(filter) != 1 {
			return nil, fmt.Errorf("each filtering rule needs to have exactly one action, but found %d actions.", len(filter))
		}

		for filterName, cfg := range filter {

			constructor, exists := filterConstructors[filterName]
			if !exists {
				return nil, fmt.Errorf("the filtering rule %s doesn't exist", filterName)
			}

			plugin, err := constructor(cfg)
			if err != nil {
				return nil, err
			}

			filters.addRule(plugin)
		}
	}

	logp.Debug("filter", "filters: %v", filters)
	return &filters, nil
}

func (filters *Filters) addRule(filter FilterRule) {

	if filters.list == nil {
		filters.list = []FilterRule{}
	}
	filters.list = append(filters.list, filter)
}

// Applies a sequence of filtering rules and returns the filtered event
func (filters *Filters) Filter(event common.MapStr) common.MapStr {

	// Check if filters are set, just return event if not
	if len(filters.list) == 0 {
		return event
	}

	// clone the event at first, before starting filtering
	filtered := event.Clone()
	var err error

	for _, filter := range filters.list {
		filtered, err = filter.Filter(filtered)
		if err != nil {
			logp.Debug("filter", "fail to apply filtering rule %s: %s", filter, err)
		}
		if filtered == nil {
			// drop event
			return nil
		}
	}

	return filtered
}

func (filters Filters) String() string {
	s := []string{}

	for _, filter := range filters.list {

		s = append(s, filter.String())
	}
	return strings.Join(s, ", ")
}
