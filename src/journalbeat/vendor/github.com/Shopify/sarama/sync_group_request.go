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

type SyncGroupRequest struct {
	GroupId          string
	GenerationId     int32
	MemberId         string
	GroupAssignments map[string][]byte
}

func (r *SyncGroupRequest) encode(pe packetEncoder) error {
	if err := pe.putString(r.GroupId); err != nil {
		return err
	}

	pe.putInt32(r.GenerationId)

	if err := pe.putString(r.MemberId); err != nil {
		return err
	}

	if err := pe.putArrayLength(len(r.GroupAssignments)); err != nil {
		return err
	}
	for memberId, memberAssignment := range r.GroupAssignments {
		if err := pe.putString(memberId); err != nil {
			return err
		}
		if err := pe.putBytes(memberAssignment); err != nil {
			return err
		}
	}

	return nil
}

func (r *SyncGroupRequest) decode(pd packetDecoder, version int16) (err error) {
	if r.GroupId, err = pd.getString(); err != nil {
		return
	}
	if r.GenerationId, err = pd.getInt32(); err != nil {
		return
	}
	if r.MemberId, err = pd.getString(); err != nil {
		return
	}

	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if n == 0 {
		return nil
	}

	r.GroupAssignments = make(map[string][]byte)
	for i := 0; i < n; i++ {
		memberId, err := pd.getString()
		if err != nil {
			return err
		}
		memberAssignment, err := pd.getBytes()
		if err != nil {
			return err
		}

		r.GroupAssignments[memberId] = memberAssignment
	}

	return nil
}

func (r *SyncGroupRequest) key() int16 {
	return 14
}

func (r *SyncGroupRequest) version() int16 {
	return 0
}

func (r *SyncGroupRequest) AddGroupAssignment(memberId string, memberAssignment []byte) {
	if r.GroupAssignments == nil {
		r.GroupAssignments = make(map[string][]byte)
	}

	r.GroupAssignments[memberId] = memberAssignment
}

func (r *SyncGroupRequest) AddGroupAssignmentMember(memberId string, memberAssignment *ConsumerGroupMemberAssignment) error {
	bin, err := encode(memberAssignment)
	if err != nil {
		return err
	}

	r.AddGroupAssignment(memberId, bin)
	return nil
}
