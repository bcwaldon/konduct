//
//Copyright 2016 Planet Labs 
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flate

// forwardCopy is like the built-in copy function except that it always goes
// forward from the start, even if the dst and src overlap.
// It is equivalent to:
//   for i := 0; i < n; i++ {
//     mem[dst+i] = mem[src+i]
//   }
func forwardCopy(mem []byte, dst, src, n int) {
	if dst <= src {
		copy(mem[dst:dst+n], mem[src:src+n])
		return
	}
	for {
		if dst >= src+n {
			copy(mem[dst:dst+n], mem[src:src+n])
			return
		}
		// There is some forward overlap.  The destination
		// will be filled with a repeated pattern of mem[src:src+k].
		// We copy one instance of the pattern here, then repeat.
		// Each time around this loop k will double.
		k := dst - src
		copy(mem[dst:dst+k], mem[src:src+k])
		n -= k
		dst += k
	}
}
