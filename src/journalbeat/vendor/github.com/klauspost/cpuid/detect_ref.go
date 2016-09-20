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

// Copyright (c) 2015 Klaus Post, released under MIT License. See LICENSE file.

// +build !amd64,!386 gccgo

package cpuid

func initCPU() {
	cpuid = func(op uint32) (eax, ebx, ecx, edx uint32) {
		return 0, 0, 0, 0
	}

	cpuidex = func(op, op2 uint32) (eax, ebx, ecx, edx uint32) {
		return 0, 0, 0, 0
	}

	xgetbv = func(index uint32) (eax, edx uint32) {
		return 0, 0
	}

	rdtscpAsm = func() (eax, ebx, ecx, edx uint32) {
		return 0, 0, 0, 0
	}
}
