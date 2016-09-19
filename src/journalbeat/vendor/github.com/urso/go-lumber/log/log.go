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
package log

import "log"

type Logging interface {
	Printf(string, ...interface{})
	Println(...interface{})
	Print(...interface{})
}

type defaultLogger struct{}

// The logger use by go-lumber
var Logger Logging = defaultLogger{}

func (defaultLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (defaultLogger) Println(args ...interface{}) {
	log.Println(args...)
}

func (defaultLogger) Print(args ...interface{}) {
	log.Print(args...)
}

func Printf(format string, args ...interface{}) {
	Logger.Printf(format, args...)
}

func Println(args ...interface{}) {
	Logger.Println(args...)
}

func Print(args ...interface{}) {
	Logger.Print(args...)
}
