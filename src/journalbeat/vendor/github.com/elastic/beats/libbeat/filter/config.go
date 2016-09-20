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
	"math"
	"strconv"

	"github.com/elastic/beats/libbeat/common"
)

type ConditionConfig struct {
	Equals   *ConditionFilter `config:"equals"`
	Contains *ConditionFilter `config:"contains"`
	Regexp   *ConditionFilter `config:"regexp"`
	Range    *ConditionFilter `config:"range"`
}

type ConditionFilter struct {
	fields map[string]interface{}
}

type FilterPluginConfig []map[string]common.Config

// fields that should be always exported
var MandatoryExportedFields = []string{"@timestamp", "type"}

func (f *ConditionFilter) Unpack(to interface{}) error {
	m, ok := to.(map[string]interface{})
	if !ok {
		return fmt.Errorf("wrong type, expect map")
	}

	f.fields = map[string]interface{}{}

	var expand func(key string, value interface{})

	expand = func(key string, value interface{}) {
		switch v := value.(type) {
		case map[string]interface{}:
			for k, val := range v {
				expand(fmt.Sprintf("%v.%v", key, k), val)
			}
		case []interface{}:
			for i, _ := range v {
				expand(fmt.Sprintf("%v.%v", key, i), v[i])
			}
		default:
			f.fields[key] = value
		}
	}

	for k, val := range m {
		expand(k, val)
	}
	return nil
}

func extractFloat(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return float64(i), nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int16:
		return float64(i), nil
	case int8:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint16:
		return float64(i), nil
	case uint8:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		f, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return math.NaN(), err
		}
		return f, err
	default:
		return math.NaN(), fmt.Errorf("unknown type %T passed to extractFloat", unk)
	}
}

func extractInt(unk interface{}) (uint64, error) {
	switch i := unk.(type) {
	case int64:
		return uint64(i), nil
	case int32:
		return uint64(i), nil
	case int16:
		return uint64(i), nil
	case int8:
		return uint64(i), nil
	case uint64:
		return uint64(i), nil
	case uint32:
		return uint64(i), nil
	case uint16:
		return uint64(i), nil
	case uint8:
		return uint64(i), nil
	case int:
		return uint64(i), nil
	case uint:
		return uint64(i), nil
	default:
		return 0, fmt.Errorf("unknown type %T passed to extractInt", unk)
	}
}

func extractString(unk interface{}) (string, error) {
	switch s := unk.(type) {
	case string:
		return string(s), nil
	default:
		return "", fmt.Errorf("unkown type %T passed to extractString", unk)
	}
}
