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

type LeaveGroupRequest struct {
	GroupId  string
	MemberId string
}

func (r *LeaveGroupRequest) encode(pe packetEncoder) error {
	if err := pe.putString(r.GroupId); err != nil {
		return err
	}
	if err := pe.putString(r.MemberId); err != nil {
		return err
	}

	return nil
}

func (r *LeaveGroupRequest) decode(pd packetDecoder, version int16) (err error) {
	if r.GroupId, err = pd.getString(); err != nil {
		return
	}
	if r.MemberId, err = pd.getString(); err != nil {
		return
	}

	return nil
}

func (r *LeaveGroupRequest) key() int16 {
	return 13
}

func (r *LeaveGroupRequest) version() int16 {
	return 0
}
