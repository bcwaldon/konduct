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
package rules

import (
	"fmt"
	"strings"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/filter"
)

type DropFields struct {
	Fields []string
	// condition
	Cond *filter.Condition
}

type DropFieldsConfig struct {
	Fields                 []string `config:"fields"`
	filter.ConditionConfig `config:",inline"`
}

func init() {
	if err := filter.RegisterPlugin("drop_fields", newDropFields); err != nil {
		panic(err)
	}
}

func newDropFields(c common.Config) (filter.FilterRule, error) {

	f := DropFields{}

	if err := f.CheckConfig(c); err != nil {
		return nil, err
	}

	config := DropFieldsConfig{}

	err := c.Unpack(&config)
	if err != nil {
		return nil, fmt.Errorf("fail to unpack the drop_fields configuration: %s", err)
	}

	/* remove read only fields */
	for _, readOnly := range filter.MandatoryExportedFields {
		for i, field := range config.Fields {
			if readOnly == field {
				config.Fields = append(config.Fields[:i], config.Fields[i+1:]...)
			}
		}
	}
	f.Fields = config.Fields

	cond, err := filter.NewCondition(config.ConditionConfig)
	if err != nil {
		return nil, err
	}
	f.Cond = cond

	return &f, nil
}

func (f *DropFields) CheckConfig(c common.Config) error {

	complete := false

	for _, field := range c.GetFields() {
		if !filter.AvailableCondition(field) {
			if field != "fields" {
				return fmt.Errorf("unexpected %s option in the drop_fields configuration", field)
			}
		}
		if field == "fields" {
			complete = true
		}
	}

	if !complete {
		return fmt.Errorf("missing fields option in the drop_fields configuration")
	}
	return nil
}

func (f *DropFields) Filter(event common.MapStr) (common.MapStr, error) {

	if f.Cond != nil && !f.Cond.Check(event) {
		return event, nil
	}

	for _, field := range f.Fields {
		err := event.Delete(field)
		if err != nil {
			return event, fmt.Errorf("Fail to delete key %s: %s", field, err)
		}

	}
	return event, nil
}

func (f DropFields) String() string {

	if f.Cond != nil {
		return "drop_fields=" + strings.Join(f.Fields, ", ") + ", condition=" + f.Cond.String()
	}
	return "drop_fields=" + strings.Join(f.Fields, ", ")

}
