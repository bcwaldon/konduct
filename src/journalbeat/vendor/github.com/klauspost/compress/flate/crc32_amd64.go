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
//+build !noasm
//+build !appengine

// Copyright 2015, Klaus Post, see LICENSE for details.

package flate

import (
	"github.com/klauspost/cpuid"
)

// crc32sse returns a hash for the first 4 bytes of the slice
// len(a) must be >= 4.
//go:noescape
func crc32sse(a []byte) hash

// crc32sseAll calculates hashes for each 4-byte set in a.
// dst must be east len(a) - 4 in size.
// The size is not checked by the assembly.
//go:noescape
func crc32sseAll(a []byte, dst []hash)

// matchLenSSE4 returns the number of matching bytes in a and b
// up to length 'max'. Both slices must be at least 'max'
// bytes in size.
// It uses the PCMPESTRI SSE 4.2 instruction.
//go:noescape
func matchLenSSE4(a, b []byte, max int) int

// histogram accumulates a histogram of b in h.
// h must be at least 256 entries in length,
// and must be cleared before calling this function.
//go:noescape
func histogram(b []byte, h []int32)

// Detect SSE 4.2 feature.
func init() {
	useSSE42 = cpuid.CPU.SSE42()
}
