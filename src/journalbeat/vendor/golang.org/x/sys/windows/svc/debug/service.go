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
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package debug provides facilities to execute svc.Handler on console.
//
package debug

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/windows/svc"
)

// Run executes service name by calling appropriate handler function.
// The process is running on console, unlike real service. Use Ctrl+C to
// send "Stop" command to your service.
func Run(name string, handler svc.Handler) error {
	cmds := make(chan svc.ChangeRequest)
	changes := make(chan svc.Status)

	sig := make(chan os.Signal)
	signal.Notify(sig)

	go func() {
		status := svc.Status{State: svc.Stopped}
		for {
			select {
			case <-sig:
				cmds <- svc.ChangeRequest{svc.Stop, status}
			case status = <-changes:
			}
		}
	}()

	_, errno := handler.Execute([]string{name}, cmds, changes)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}
