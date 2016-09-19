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

import (
	"reflect"
	"strings"
)

func (prop *ModelProperty) setDescription(field reflect.StructField) {
	if tag := field.Tag.Get("description"); tag != "" {
		prop.Description = tag
	}
}

func (prop *ModelProperty) setDefaultValue(field reflect.StructField) {
	if tag := field.Tag.Get("default"); tag != "" {
		prop.DefaultValue = Special(tag)
	}
}

func (prop *ModelProperty) setEnumValues(field reflect.StructField) {
	// We use | to separate the enum values.  This value is chosen
	// since its unlikely to be useful in actual enumeration values.
	if tag := field.Tag.Get("enum"); tag != "" {
		prop.Enum = strings.Split(tag, "|")
	}
}

func (prop *ModelProperty) setMaximum(field reflect.StructField) {
	if tag := field.Tag.Get("maximum"); tag != "" {
		prop.Maximum = tag
	}
}

func (prop *ModelProperty) setType(field reflect.StructField) {
	if tag := field.Tag.Get("type"); tag != "" {
		prop.Type = &tag
	}
}

func (prop *ModelProperty) setMinimum(field reflect.StructField) {
	if tag := field.Tag.Get("minimum"); tag != "" {
		prop.Minimum = tag
	}
}

func (prop *ModelProperty) setUniqueItems(field reflect.StructField) {
	tag := field.Tag.Get("unique")
	switch tag {
	case "true":
		v := true
		prop.UniqueItems = &v
	case "false":
		v := false
		prop.UniqueItems = &v
	}
}

func (prop *ModelProperty) setPropertyMetadata(field reflect.StructField) {
	prop.setDescription(field)
	prop.setEnumValues(field)
	prop.setMinimum(field)
	prop.setMaximum(field)
	prop.setUniqueItems(field)
	prop.setDefaultValue(field)
	prop.setType(field)
}
