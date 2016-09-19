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
//+build !amd64 noasm appengine

// Copyright 2015, Klaus Post, see LICENSE for details.

package flate

func init() {
	useSSE42 = false
}

// crc32sse should never be called.
func crc32sse(a []byte) hash {
	panic("no assembler")
}

// crc32sseAll should never be called.
func crc32sseAll(a []byte, dst []hash) {
	panic("no assembler")
}

// matchLenSSE4 should never be called.
func matchLenSSE4(a, b []byte, max int) int {
	panic("no assembler")
	return 0
}

// histogram accumulates a histogram of b in h.
// h must be at least 256 entries in length,
// and must be cleared before calling this function.
func histogram(b []byte, h []int32) {
	for _, t := range b {
		h[t]++
	}
}
