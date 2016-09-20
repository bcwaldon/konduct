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
package sarama

import (
	"encoding/binary"

	"github.com/klauspost/crc32"
)

// crc32Field implements the pushEncoder and pushDecoder interfaces for calculating CRC32s.
type crc32Field struct {
	startOffset int
}

func (c *crc32Field) saveOffset(in int) {
	c.startOffset = in
}

func (c *crc32Field) reserveLength() int {
	return 4
}

func (c *crc32Field) run(curOffset int, buf []byte) error {
	crc := crc32.ChecksumIEEE(buf[c.startOffset+4 : curOffset])
	binary.BigEndian.PutUint32(buf[c.startOffset:], crc)
	return nil
}

func (c *crc32Field) check(curOffset int, buf []byte) error {
	crc := crc32.ChecksumIEEE(buf[c.startOffset+4 : curOffset])

	if crc != binary.BigEndian.Uint32(buf[c.startOffset:]) {
		return PacketDecodingError{"CRC didn't match"}
	}

	return nil
}
