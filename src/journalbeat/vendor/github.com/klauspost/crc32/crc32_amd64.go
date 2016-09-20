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

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine,!gccgo

package crc32

// This file contains the code to call the SSE 4.2 version of the Castagnoli
// and IEEE CRC.

// haveSSE41/haveSSE42/haveCLMUL are defined in crc_amd64.s and use
// CPUID to test for SSE 4.1, 4.2 and CLMUL support.
func haveSSE41() bool
func haveSSE42() bool
func haveCLMUL() bool

// castagnoliSSE42 is defined in crc_amd64.s and uses the SSE4.2 CRC32
// instruction.
//go:noescape
func castagnoliSSE42(crc uint32, p []byte) uint32

// ieeeCLMUL is defined in crc_amd64.s and uses the PCLMULQDQ
// instruction as well as SSE 4.1.
//go:noescape
func ieeeCLMUL(crc uint32, p []byte) uint32

var sse42 = haveSSE42()
var useFastIEEE = haveCLMUL() && haveSSE41()

func updateCastagnoli(crc uint32, p []byte) uint32 {
	if sse42 {
		return castagnoliSSE42(crc, p)
	}
	// only use slicing-by-8 when input is >= 16 Bytes
	if len(p) >= 16 {
		return updateSlicingBy8(crc, castagnoliTable8, p)
	}
	return update(crc, castagnoliTable, p)
}

func updateIEEE(crc uint32, p []byte) uint32 {
	if useFastIEEE && len(p) >= 64 {
		left := len(p) & 15
		do := len(p) - left
		crc = ^ieeeCLMUL(^crc, p[:do])
		if left > 0 {
			crc = update(crc, IEEETable, p[do:])
		}
		return crc
	}

	// only use slicing-by-8 when input is >= 16 Bytes
	if len(p) >= 16 {
		ieeeTable8Once.Do(func() {
			ieeeTable8 = makeTable8(IEEE)
		})
		return updateSlicingBy8(crc, ieeeTable8, p)
	}

	return update(crc, IEEETable, p)
}
