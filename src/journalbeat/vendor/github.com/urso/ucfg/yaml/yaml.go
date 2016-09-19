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
package yaml

import (
	"io/ioutil"

	"github.com/urso/ucfg"
	"gopkg.in/yaml.v2"
)

func NewConfig(in []byte, opts ...ucfg.Option) (*ucfg.Config, error) {
	var m map[string]interface{}
	if err := yaml.Unmarshal(in, &m); err != nil {
		return nil, err
	}
	return ucfg.NewFrom(m, opts...)
}

func NewConfigWithFile(name string, opts ...ucfg.Option) (*ucfg.Config, error) {
	input, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	opts = append([]ucfg.Option{ucfg.MetaData(ucfg.Meta{name})}, opts...)
	return NewConfig(input, opts...)
}
