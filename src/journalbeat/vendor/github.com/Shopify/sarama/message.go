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
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/eapache/go-xerial-snappy"
)

// CompressionCodec represents the various compression codecs recognized by Kafka in messages.
type CompressionCodec int8

// only the last two bits are really used
const compressionCodecMask int8 = 0x03

const (
	CompressionNone   CompressionCodec = 0
	CompressionGZIP   CompressionCodec = 1
	CompressionSnappy CompressionCodec = 2
)

// The spec just says: "This is a version id used to allow backwards compatible evolution of the message
// binary format." but it doesn't say what the current value is, so presumably 0...
const messageFormat int8 = 0

type Message struct {
	Codec CompressionCodec // codec used to compress the message contents
	Key   []byte           // the message key, may be nil
	Value []byte           // the message contents
	Set   *MessageSet      // the message set a message might wrap

	compressedCache []byte
}

func (m *Message) encode(pe packetEncoder) error {
	pe.push(&crc32Field{})

	pe.putInt8(messageFormat)

	attributes := int8(m.Codec) & compressionCodecMask
	pe.putInt8(attributes)

	err := pe.putBytes(m.Key)
	if err != nil {
		return err
	}

	var payload []byte

	if m.compressedCache != nil {
		payload = m.compressedCache
		m.compressedCache = nil
	} else if m.Value != nil {
		switch m.Codec {
		case CompressionNone:
			payload = m.Value
		case CompressionGZIP:
			var buf bytes.Buffer
			writer := gzip.NewWriter(&buf)
			if _, err = writer.Write(m.Value); err != nil {
				return err
			}
			if err = writer.Close(); err != nil {
				return err
			}
			m.compressedCache = buf.Bytes()
			payload = m.compressedCache
		case CompressionSnappy:
			tmp := snappy.Encode(m.Value)
			m.compressedCache = tmp
			payload = m.compressedCache
		default:
			return PacketEncodingError{fmt.Sprintf("unsupported compression codec (%d)", m.Codec)}
		}
	}

	if err = pe.putBytes(payload); err != nil {
		return err
	}

	return pe.pop()
}

func (m *Message) decode(pd packetDecoder) (err error) {
	err = pd.push(&crc32Field{})
	if err != nil {
		return err
	}

	format, err := pd.getInt8()
	if err != nil {
		return err
	}
	if format != messageFormat {
		return PacketDecodingError{"unexpected messageFormat"}
	}

	attribute, err := pd.getInt8()
	if err != nil {
		return err
	}
	m.Codec = CompressionCodec(attribute & compressionCodecMask)

	m.Key, err = pd.getBytes()
	if err != nil {
		return err
	}

	m.Value, err = pd.getBytes()
	if err != nil {
		return err
	}

	switch m.Codec {
	case CompressionNone:
		// nothing to do
	case CompressionGZIP:
		if m.Value == nil {
			break
		}
		reader, err := gzip.NewReader(bytes.NewReader(m.Value))
		if err != nil {
			return err
		}
		if m.Value, err = ioutil.ReadAll(reader); err != nil {
			return err
		}
		if err := m.decodeSet(); err != nil {
			return err
		}
	case CompressionSnappy:
		if m.Value == nil {
			break
		}
		if m.Value, err = snappy.Decode(m.Value); err != nil {
			return err
		}
		if err := m.decodeSet(); err != nil {
			return err
		}
	default:
		return PacketDecodingError{fmt.Sprintf("invalid compression specified (%d)", m.Codec)}
	}

	return pd.pop()
}

// decodes a message set from a previousy encoded bulk-message
func (m *Message) decodeSet() (err error) {
	pd := realDecoder{raw: m.Value}
	m.Set = &MessageSet{}
	return m.Set.decode(&pd)
}
