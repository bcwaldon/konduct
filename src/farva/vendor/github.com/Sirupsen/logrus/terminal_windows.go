// Copyright 2016 Planet Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Based on ssh/terminal:
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package logrus

import (
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
)

// IsTerminal returns true if stderr's file descriptor is a terminal.
func IsTerminal() bool {
	fd := syscall.Stderr
	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}
