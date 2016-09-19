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

type JoinGroupRequest struct {
	GroupId        string
	SessionTimeout int32
	MemberId       string
	ProtocolType   string
	GroupProtocols map[string][]byte
}

func (r *JoinGroupRequest) encode(pe packetEncoder) error {
	if err := pe.putString(r.GroupId); err != nil {
		return err
	}
	pe.putInt32(r.SessionTimeout)
	if err := pe.putString(r.MemberId); err != nil {
		return err
	}
	if err := pe.putString(r.ProtocolType); err != nil {
		return err
	}

	if err := pe.putArrayLength(len(r.GroupProtocols)); err != nil {
		return err
	}
	for name, metadata := range r.GroupProtocols {
		if err := pe.putString(name); err != nil {
			return err
		}
		if err := pe.putBytes(metadata); err != nil {
			return err
		}
	}

	return nil
}

func (r *JoinGroupRequest) decode(pd packetDecoder, version int16) (err error) {
	if r.GroupId, err = pd.getString(); err != nil {
		return
	}

	if r.SessionTimeout, err = pd.getInt32(); err != nil {
		return
	}

	if r.MemberId, err = pd.getString(); err != nil {
		return
	}

	if r.ProtocolType, err = pd.getString(); err != nil {
		return
	}

	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if n == 0 {
		return nil
	}

	r.GroupProtocols = make(map[string][]byte)
	for i := 0; i < n; i++ {
		name, err := pd.getString()
		if err != nil {
			return err
		}
		metadata, err := pd.getBytes()
		if err != nil {
			return err
		}

		r.GroupProtocols[name] = metadata
	}

	return nil
}

func (r *JoinGroupRequest) key() int16 {
	return 11
}

func (r *JoinGroupRequest) version() int16 {
	return 0
}

func (r *JoinGroupRequest) AddGroupProtocol(name string, metadata []byte) {
	if r.GroupProtocols == nil {
		r.GroupProtocols = make(map[string][]byte)
	}

	r.GroupProtocols[name] = metadata
}

func (r *JoinGroupRequest) AddGroupProtocolMetadata(name string, metadata *ConsumerGroupMemberMetadata) error {
	bin, err := encode(metadata)
	if err != nil {
		return err
	}

	r.AddGroupProtocol(name, bin)
	return nil
}
