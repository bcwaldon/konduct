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

// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows
// +build go1.3

package svc

import "unsafe"

const ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

// Should be a built-in for unsafe.Pointer?
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// funcPC returns the entry PC of the function f.
// It assumes that f is a func value. Otherwise the behavior is undefined.
func funcPC(f interface{}) uintptr {
	return **(**uintptr)(add(unsafe.Pointer(&f), ptrSize))
}

// from sys_386.s and sys_amd64.s
func servicectlhandler(ctl uint32) uintptr
func servicemain(argc uint32, argv **uint16)

func getServiceMain(r *uintptr) {
	*r = funcPC(servicemain)
}
