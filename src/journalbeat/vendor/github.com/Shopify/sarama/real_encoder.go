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

type realEncoder struct {
	raw   []byte
	off   int
	stack []pushEncoder
}

// primitives

func (re *realEncoder) putInt8(in int8) {
	re.raw[re.off] = byte(in)
	re.off++
}

func (re *realEncoder) putInt16(in int16) {
	binary.BigEndian.PutUint16(re.raw[re.off:], uint16(in))
	re.off += 2
}

func (re *realEncoder) putInt32(in int32) {
	binary.BigEndian.PutUint32(re.raw[re.off:], uint32(in))
	re.off += 4
}

func (re *realEncoder) putInt64(in int64) {
	binary.BigEndian.PutUint64(re.raw[re.off:], uint64(in))
	re.off += 8
}

func (re *realEncoder) putArrayLength(in int) error {
	re.putInt32(int32(in))
	return nil
}

// collection

func (re *realEncoder) putRawBytes(in []byte) error {
	copy(re.raw[re.off:], in)
	re.off += len(in)
	return nil
}

func (re *realEncoder) putBytes(in []byte) error {
	if in == nil {
		re.putInt32(-1)
		return nil
	}
	re.putInt32(int32(len(in)))
	copy(re.raw[re.off:], in)
	re.off += len(in)
	return nil
}

func (re *realEncoder) putString(in string) error {
	re.putInt16(int16(len(in)))
	copy(re.raw[re.off:], in)
	re.off += len(in)
	return nil
}

func (re *realEncoder) putStringArray(in []string) error {
	err := re.putArrayLength(len(in))
	if err != nil {
		return err
	}

	for _, val := range in {
		if err := re.putString(val); err != nil {
			return err
		}
	}

	return nil
}

func (re *realEncoder) putInt32Array(in []int32) error {
	err := re.putArrayLength(len(in))
	if err != nil {
		return err
	}
	for _, val := range in {
		re.putInt32(val)
	}
	return nil
}

func (re *realEncoder) putInt64Array(in []int64) error {
	err := re.putArrayLength(len(in))
	if err != nil {
		return err
	}
	for _, val := range in {
		re.putInt64(val)
	}
	return nil
}

// stacks

func (re *realEncoder) push(in pushEncoder) {
	in.saveOffset(re.off)
	re.off += in.reserveLength()
	re.stack = append(re.stack, in)
}

func (re *realEncoder) pop() error {
	// this is go's ugly pop pattern (the inverse of append)
	in := re.stack[len(re.stack)-1]
	re.stack = re.stack[:len(re.stack)-1]

	return in.run(re.off, re.raw)
}
