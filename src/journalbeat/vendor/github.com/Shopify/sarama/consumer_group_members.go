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

type ConsumerGroupMemberMetadata struct {
	Version  int16
	Topics   []string
	UserData []byte
}

func (m *ConsumerGroupMemberMetadata) encode(pe packetEncoder) error {
	pe.putInt16(m.Version)

	if err := pe.putStringArray(m.Topics); err != nil {
		return err
	}

	if err := pe.putBytes(m.UserData); err != nil {
		return err
	}

	return nil
}

func (m *ConsumerGroupMemberMetadata) decode(pd packetDecoder) (err error) {
	if m.Version, err = pd.getInt16(); err != nil {
		return
	}

	if m.Topics, err = pd.getStringArray(); err != nil {
		return
	}

	if m.UserData, err = pd.getBytes(); err != nil {
		return
	}

	return nil
}

type ConsumerGroupMemberAssignment struct {
	Version  int16
	Topics   map[string][]int32
	UserData []byte
}

func (m *ConsumerGroupMemberAssignment) encode(pe packetEncoder) error {
	pe.putInt16(m.Version)

	if err := pe.putArrayLength(len(m.Topics)); err != nil {
		return err
	}

	for topic, partitions := range m.Topics {
		if err := pe.putString(topic); err != nil {
			return err
		}
		if err := pe.putInt32Array(partitions); err != nil {
			return err
		}
	}

	if err := pe.putBytes(m.UserData); err != nil {
		return err
	}

	return nil
}

func (m *ConsumerGroupMemberAssignment) decode(pd packetDecoder) (err error) {
	if m.Version, err = pd.getInt16(); err != nil {
		return
	}

	var topicLen int
	if topicLen, err = pd.getArrayLength(); err != nil {
		return
	}

	m.Topics = make(map[string][]int32, topicLen)
	for i := 0; i < topicLen; i++ {
		var topic string
		if topic, err = pd.getString(); err != nil {
			return
		}
		if m.Topics[topic], err = pd.getInt32Array(); err != nil {
			return
		}
	}

	if m.UserData, err = pd.getBytes(); err != nil {
		return
	}

	return nil
}
