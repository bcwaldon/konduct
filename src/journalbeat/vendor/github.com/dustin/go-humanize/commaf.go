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
// +build go1.6

package humanize

import (
	"bytes"
	"math/big"
	"strings"
)

// BigCommaf produces a string form of the given big.Float in base 10
// with commas after every three orders of magnitude.
func BigCommaf(v *big.Float) string {
	buf := &bytes.Buffer{}
	if v.Sign() < 0 {
		buf.Write([]byte{'-'})
		v.Abs(v)
	}

	comma := []byte{','}

	parts := strings.Split(v.Text('f', -1), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}
