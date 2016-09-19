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
// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build gccgo

package unix

import "syscall"

// We can't use the gc-syntax .s files for gccgo.  On the plus side
// much of the functionality can be written directly in Go.

//extern gccgoRealSyscall
func realSyscall(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r, errno uintptr)

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	syscall.Entersyscall()
	r, errno := realSyscall(trap, a1, a2, a3, 0, 0, 0, 0, 0, 0)
	syscall.Exitsyscall()
	return r, 0, syscall.Errno(errno)
}

func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	syscall.Entersyscall()
	r, errno := realSyscall(trap, a1, a2, a3, a4, a5, a6, 0, 0, 0)
	syscall.Exitsyscall()
	return r, 0, syscall.Errno(errno)
}

func Syscall9(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	syscall.Entersyscall()
	r, errno := realSyscall(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	syscall.Exitsyscall()
	return r, 0, syscall.Errno(errno)
}

func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	r, errno := realSyscall(trap, a1, a2, a3, 0, 0, 0, 0, 0, 0)
	return r, 0, syscall.Errno(errno)
}

func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	r, errno := realSyscall(trap, a1, a2, a3, a4, a5, a6, 0, 0, 0)
	return r, 0, syscall.Errno(errno)
}
