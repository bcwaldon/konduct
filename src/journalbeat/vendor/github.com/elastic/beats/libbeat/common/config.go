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
package common

import (
	"github.com/urso/ucfg"
	"github.com/urso/ucfg/yaml"
)

type Config ucfg.Config

func NewConfig() *Config {
	return fromConfig(ucfg.New())
}

func NewConfigFrom(from interface{}) (*Config, error) {
	c, err := ucfg.NewFrom(from, ucfg.PathSep("."))
	return fromConfig(c), err
}

func NewConfigWithYAML(in []byte, source string) (*Config, error) {
	c, err := yaml.NewConfig(in, ucfg.PathSep("."), ucfg.MetaData(ucfg.Meta{source}))
	return fromConfig(c), err
}

func LoadFile(path string) (*Config, error) {
	c, err := yaml.NewConfigWithFile(path, ucfg.PathSep("."))
	return fromConfig(c), err
}

func (c *Config) Merge(from interface{}) error {
	return c.access().Merge(from, ucfg.PathSep("."))
}

func (c *Config) Unpack(to interface{}) error {
	return c.access().Unpack(to, ucfg.PathSep("."))
}

func (c *Config) Path() string {
	return c.access().Path(".")
}

func (c *Config) PathOf(field string) string {
	return c.access().PathOf(field, ".")
}

func (c *Config) HasField(name string) bool {
	return c.access().HasField(name)
}

func (c *Config) CountField(name string) (int, error) {
	return c.access().CountField(name)
}

func (c *Config) Bool(name string, idx int) (bool, error) {
	return c.access().Bool(name, idx, ucfg.PathSep("."))
}

func (c *Config) String(name string, idx int) (string, error) {
	return c.access().String(name, idx, ucfg.PathSep("."))
}

func (c *Config) Int(name string, idx int) (int64, error) {
	return c.access().Int(name, idx, ucfg.PathSep("."))
}

func (c *Config) Float(name string, idx int) (float64, error) {
	return c.access().Float(name, idx, ucfg.PathSep("."))
}

func (c *Config) Child(name string, idx int) (*Config, error) {
	sub, err := c.access().Child(name, idx, ucfg.PathSep("."))
	return fromConfig(sub), err
}

func (c *Config) SetBool(name string, idx int, value bool) error {
	return c.access().SetBool(name, idx, value, ucfg.PathSep("."))
}

func (c *Config) SetInt(name string, idx int, value int64) error {
	return c.access().SetInt(name, idx, value, ucfg.PathSep("."))
}

func (c *Config) SetFloat(name string, idx int, value float64) error {
	return c.access().SetFloat(name, idx, value, ucfg.PathSep("."))
}

func (c *Config) SetString(name string, idx int, value string) error {
	return c.access().SetString(name, idx, value, ucfg.PathSep("."))
}

func (c *Config) SetChild(name string, idx int, value *Config) error {
	return c.access().SetChild(name, idx, value.access(), ucfg.PathSep("."))
}

func fromConfig(in *ucfg.Config) *Config {
	return (*Config)(in)
}

func (c *Config) access() *ucfg.Config {
	return (*ucfg.Config)(c)
}
func (c *Config) GetFields() []string {
	return c.access().GetFields()
}
