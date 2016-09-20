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

type MetadataRequest struct {
	Topics []string
}

func (mr *MetadataRequest) encode(pe packetEncoder) error {
	err := pe.putArrayLength(len(mr.Topics))
	if err != nil {
		return err
	}

	for i := range mr.Topics {
		err = pe.putString(mr.Topics[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr *MetadataRequest) decode(pd packetDecoder, version int16) error {
	topicCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if topicCount == 0 {
		return nil
	}

	mr.Topics = make([]string, topicCount)
	for i := range mr.Topics {
		topic, err := pd.getString()
		if err != nil {
			return err
		}
		mr.Topics[i] = topic
	}
	return nil
}

func (mr *MetadataRequest) key() int16 {
	return 3
}

func (mr *MetadataRequest) version() int16 {
	return 0
}
