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
package pflag

import (
	"encoding/csv"
	"fmt"
	"strings"
)

var _ = fmt.Fprint

// -- stringSlice Value
type stringSliceValue struct {
	value   *[]string
	changed bool
}

func newStringSliceValue(val []string, p *[]string) *stringSliceValue {
	ssv := new(stringSliceValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringSliceValue) Set(val string) error {
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	v, err := csvReader.Read()
	if err != nil {
		return err
	}
	if !s.changed {
		*s.value = v
	} else {
		*s.value = append(*s.value, v...)
	}
	s.changed = true
	return nil
}

func (s *stringSliceValue) Type() string {
	return "stringSlice"
}

func (s *stringSliceValue) String() string { return "[" + strings.Join(*s.value, ",") + "]" }

func stringSliceConv(sval string) (interface{}, error) {
	sval = strings.Trim(sval, "[]")
	// An empty string would cause a slice with one (empty) string
	if len(sval) == 0 {
		return []string{}, nil
	}
	v := strings.Split(sval, ",")
	return v, nil
}

// GetStringSlice return the []string value of a flag with the given name
func (f *FlagSet) GetStringSlice(name string) ([]string, error) {
	val, err := f.getFlagType(name, "stringSlice", stringSliceConv)
	if err != nil {
		return []string{}, err
	}
	return val.([]string), nil
}

// StringSliceVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
func (f *FlagSet) StringSliceVar(p *[]string, name string, value []string, usage string) {
	f.VarP(newStringSliceValue(value, p), name, "", usage)
}

// StringSliceVarP is like StringSliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringSliceVarP(p *[]string, name, shorthand string, value []string, usage string) {
	f.VarP(newStringSliceValue(value, p), name, shorthand, usage)
}

// StringSliceVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
func StringSliceVar(p *[]string, name string, value []string, usage string) {
	CommandLine.VarP(newStringSliceValue(value, p), name, "", usage)
}

// StringSliceVarP is like StringSliceVar, but accepts a shorthand letter that can be used after a single dash.
func StringSliceVarP(p *[]string, name, shorthand string, value []string, usage string) {
	CommandLine.VarP(newStringSliceValue(value, p), name, shorthand, usage)
}

// StringSlice defines a string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
func (f *FlagSet) StringSlice(name string, value []string, usage string) *[]string {
	p := []string{}
	f.StringSliceVarP(&p, name, "", value, usage)
	return &p
}

// StringSliceP is like StringSlice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringSliceP(name, shorthand string, value []string, usage string) *[]string {
	p := []string{}
	f.StringSliceVarP(&p, name, shorthand, value, usage)
	return &p
}

// StringSlice defines a string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
func StringSlice(name string, value []string, usage string) *[]string {
	return CommandLine.StringSliceP(name, "", value, usage)
}

// StringSliceP is like StringSlice, but accepts a shorthand letter that can be used after a single dash.
func StringSliceP(name, shorthand string, value []string, usage string) *[]string {
	return CommandLine.StringSliceP(name, shorthand, value, usage)
}
