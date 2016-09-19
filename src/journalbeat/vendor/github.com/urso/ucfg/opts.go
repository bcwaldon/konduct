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
package ucfg

import "os"

type Option func(*options)

type options struct {
	tag          string
	validatorTag string
	pathSep      string
	meta         *Meta
	env          []*Config
	resolvers    []func(name string) (string, error)
	varexp       bool
}

func StructTag(tag string) Option {
	return func(o *options) {
		o.tag = tag
	}
}

func ValidatorTag(tag string) Option {
	return func(o *options) {
		o.validatorTag = tag
	}
}

func PathSep(sep string) Option {
	return func(o *options) {
		o.pathSep = sep
	}
}

func MetaData(meta Meta) Option {
	return func(o *options) {
		o.meta = &meta
	}
}

func Env(e *Config) Option {
	return func(o *options) {
		o.env = append(o.env, e)
	}
}

func Resolve(fn func(name string) (string, error)) Option {
	return func(o *options) {
		o.resolvers = append(o.resolvers, fn)
	}
}

var ResolveEnv = Resolve(func(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", ErrMissing
	}
	return value, nil
})

var VarExp Option = func(o *options) {
	o.varexp = true
}

func makeOptions(opts []Option) *options {
	o := options{
		tag:          "config",
		validatorTag: "validate",
		pathSep:      "", // no separator by default
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &o
}
