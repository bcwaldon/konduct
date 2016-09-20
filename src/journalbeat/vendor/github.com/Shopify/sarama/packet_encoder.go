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

// PacketEncoder is the interface providing helpers for writing with Kafka's encoding rules.
// Types implementing Encoder only need to worry about calling methods like PutString,
// not about how a string is represented in Kafka.
type packetEncoder interface {
	// Primitives
	putInt8(in int8)
	putInt16(in int16)
	putInt32(in int32)
	putInt64(in int64)
	putArrayLength(in int) error

	// Collections
	putBytes(in []byte) error
	putRawBytes(in []byte) error
	putString(in string) error
	putStringArray(in []string) error
	putInt32Array(in []int32) error
	putInt64Array(in []int64) error

	// Stacks, see PushEncoder
	push(in pushEncoder)
	pop() error
}

// PushEncoder is the interface for encoding fields like CRCs and lengths where the value
// of the field depends on what is encoded after it in the packet. Start them with PacketEncoder.Push() where
// the actual value is located in the packet, then PacketEncoder.Pop() them when all the bytes they
// depend upon have been written.
type pushEncoder interface {
	// Saves the offset into the input buffer as the location to actually write the calculated value when able.
	saveOffset(in int)

	// Returns the length of data to reserve for the output of this encoder (eg 4 bytes for a CRC32).
	reserveLength() int

	// Indicates that all required data is now available to calculate and write the field.
	// SaveOffset is guaranteed to have been called first. The implementation should write ReserveLength() bytes
	// of data to the saved offset, based on the data between the saved offset and curOffset.
	run(curOffset int, buf []byte) error
}
