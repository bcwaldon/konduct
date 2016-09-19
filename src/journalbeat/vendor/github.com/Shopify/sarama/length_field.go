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

import "encoding/binary"

// LengthField implements the PushEncoder and PushDecoder interfaces for calculating 4-byte lengths.
type lengthField struct {
	startOffset int
}

func (l *lengthField) saveOffset(in int) {
	l.startOffset = in
}

func (l *lengthField) reserveLength() int {
	return 4
}

func (l *lengthField) run(curOffset int, buf []byte) error {
	binary.BigEndian.PutUint32(buf[l.startOffset:], uint32(curOffset-l.startOffset-4))
	return nil
}

func (l *lengthField) check(curOffset int, buf []byte) error {
	if uint32(curOffset-l.startOffset-4) != binary.BigEndian.Uint32(buf[l.startOffset:]) {
		return PacketDecodingError{"length field invalid"}
	}

	return nil
}
