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
package snappy

import (
	"bytes"
	"encoding/binary"

	master "github.com/golang/snappy"
)

var xerialHeader = []byte{130, 83, 78, 65, 80, 80, 89, 0}

// Encode encodes data as snappy with no framing header.
func Encode(src []byte) []byte {
	return master.Encode(nil, src)
}

// Decode decodes snappy data whether it is traditional unframed
// or includes the xerial framing format.
func Decode(src []byte) ([]byte, error) {
	if !bytes.Equal(src[:8], xerialHeader) {
		return master.Decode(nil, src)
	}

	var (
		pos   = uint32(16)
		max   = uint32(len(src))
		dst   = make([]byte, 0, len(src))
		chunk []byte
		err   error
	)
	for pos < max {
		size := binary.BigEndian.Uint32(src[pos : pos+4])
		pos += 4

		chunk, err = master.Decode(chunk, src[pos:pos+size])
		if err != nil {
			return nil, err
		}
		pos += size
		dst = append(dst, chunk...)
	}
	return dst, nil
}
