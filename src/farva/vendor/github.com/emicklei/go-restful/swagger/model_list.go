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
package swagger

// Copyright 2015 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"bytes"
	"encoding/json"
)

// NamedModel associates a name with a Model (not using its Id)
type NamedModel struct {
	Name  string
	Model Model
}

// ModelList encapsulates a list of NamedModel (association)
type ModelList struct {
	List []NamedModel
}

// Put adds or replaces a Model by its name
func (l *ModelList) Put(name string, model Model) {
	for i, each := range l.List {
		if each.Name == name {
			// replace
			l.List[i] = NamedModel{name, model}
			return
		}
	}
	// add
	l.List = append(l.List, NamedModel{name, model})
}

// At returns a Model by its name, ok is false if absent
func (l *ModelList) At(name string) (m Model, ok bool) {
	for _, each := range l.List {
		if each.Name == name {
			return each.Model, true
		}
	}
	return m, false
}

// Do enumerates all the models, each with its assigned name
func (l *ModelList) Do(block func(name string, value Model)) {
	for _, each := range l.List {
		block(each.Name, each.Model)
	}
}

// MarshalJSON writes the ModelList as if it was a map[string]Model
func (l ModelList) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	buf.WriteString("{\n")
	for i, each := range l.List {
		buf.WriteString("\"")
		buf.WriteString(each.Name)
		buf.WriteString("\": ")
		encoder.Encode(each.Model)
		if i < len(l.List)-1 {
			buf.WriteString(",\n")
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// UnmarshalJSON reads back a ModelList. This is an expensive operation.
func (l *ModelList) UnmarshalJSON(data []byte) error {
	raw := map[string]interface{}{}
	json.NewDecoder(bytes.NewReader(data)).Decode(&raw)
	for k, v := range raw {
		// produces JSON bytes for each value
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		var m Model
		json.NewDecoder(bytes.NewReader(data)).Decode(&m)
		l.Put(k, m)
	}
	return nil
}
