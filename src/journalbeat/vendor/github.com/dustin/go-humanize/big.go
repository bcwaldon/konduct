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
package humanize

import (
	"math/big"
)

// order of magnitude (to a max order)
func oomm(n, b *big.Int, maxmag int) (float64, int) {
	mag := 0
	m := &big.Int{}
	for n.Cmp(b) >= 0 {
		n.DivMod(n, b, m)
		mag++
		if mag == maxmag && maxmag >= 0 {
			break
		}
	}
	return float64(n.Int64()) + (float64(m.Int64()) / float64(b.Int64())), mag
}

// total order of magnitude
// (same as above, but with no upper limit)
func oom(n, b *big.Int) (float64, int) {
	mag := 0
	m := &big.Int{}
	for n.Cmp(b) >= 0 {
		n.DivMod(n, b, m)
		mag++
	}
	return float64(n.Int64()) + (float64(m.Int64()) / float64(b.Int64())), mag
}
