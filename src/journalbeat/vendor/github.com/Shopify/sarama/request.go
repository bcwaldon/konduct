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
	"fmt"
	"io"
)

type protocolBody interface {
	encoder
	versionedDecoder
	key() int16
	version() int16
}

type request struct {
	correlationID int32
	clientID      string
	body          protocolBody
}

func (r *request) encode(pe packetEncoder) (err error) {
	pe.push(&lengthField{})
	pe.putInt16(r.body.key())
	pe.putInt16(r.body.version())
	pe.putInt32(r.correlationID)
	err = pe.putString(r.clientID)
	if err != nil {
		return err
	}
	err = r.body.encode(pe)
	if err != nil {
		return err
	}
	return pe.pop()
}

func (r *request) decode(pd packetDecoder) (err error) {
	var key int16
	if key, err = pd.getInt16(); err != nil {
		return err
	}
	var version int16
	if version, err = pd.getInt16(); err != nil {
		return err
	}
	if r.correlationID, err = pd.getInt32(); err != nil {
		return err
	}
	r.clientID, err = pd.getString()

	r.body = allocateBody(key, version)
	if r.body == nil {
		return PacketDecodingError{fmt.Sprintf("unknown request key (%d)", key)}
	}
	return r.body.decode(pd, version)
}

func decodeRequest(r io.Reader) (req *request, err error) {
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBytes); err != nil {
		return nil, err
	}

	length := int32(binary.BigEndian.Uint32(lengthBytes))
	if length <= 4 || length > MaxRequestSize {
		return nil, PacketDecodingError{fmt.Sprintf("message of length %d too large or too small", length)}
	}

	encodedReq := make([]byte, length)
	if _, err := io.ReadFull(r, encodedReq); err != nil {
		return nil, err
	}

	req = &request{}
	if err := decode(encodedReq, req); err != nil {
		return nil, err
	}
	return req, nil
}

func allocateBody(key, version int16) protocolBody {
	switch key {
	case 0:
		return &ProduceRequest{}
	case 1:
		return &FetchRequest{}
	case 2:
		return &OffsetRequest{}
	case 3:
		return &MetadataRequest{}
	case 8:
		return &OffsetCommitRequest{Version: version}
	case 9:
		return &OffsetFetchRequest{}
	case 10:
		return &ConsumerMetadataRequest{}
	case 11:
		return &JoinGroupRequest{}
	case 12:
		return &HeartbeatRequest{}
	case 13:
		return &LeaveGroupRequest{}
	case 14:
		return &SyncGroupRequest{}
	case 15:
		return &DescribeGroupsRequest{}
	case 16:
		return &ListGroupsRequest{}
	case 17:
		return &SaslHandshakeRequest{}
	case 18:
		return &ApiVersionsRequest{}
	}
	return nil
}
