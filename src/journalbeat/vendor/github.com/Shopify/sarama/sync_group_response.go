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

type SyncGroupResponse struct {
	Err              KError
	MemberAssignment []byte
}

func (r *SyncGroupResponse) GetMemberAssignment() (*ConsumerGroupMemberAssignment, error) {
	assignment := new(ConsumerGroupMemberAssignment)
	err := decode(r.MemberAssignment, assignment)
	return assignment, err
}

func (r *SyncGroupResponse) encode(pe packetEncoder) error {
	pe.putInt16(int16(r.Err))
	return pe.putBytes(r.MemberAssignment)
}

func (r *SyncGroupResponse) decode(pd packetDecoder, version int16) (err error) {
	if kerr, err := pd.getInt16(); err != nil {
		return err
	} else {
		r.Err = KError(kerr)
	}

	r.MemberAssignment, err = pd.getBytes()
	return
}

func (r *SyncGroupResponse) key() int16 {
	return 14
}

func (r *SyncGroupResponse) version() int16 {
	return 0
}
