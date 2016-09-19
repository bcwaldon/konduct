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

type OffsetCommitResponse struct {
	Errors map[string]map[int32]KError
}

func (r *OffsetCommitResponse) AddError(topic string, partition int32, kerror KError) {
	if r.Errors == nil {
		r.Errors = make(map[string]map[int32]KError)
	}
	partitions := r.Errors[topic]
	if partitions == nil {
		partitions = make(map[int32]KError)
		r.Errors[topic] = partitions
	}
	partitions[partition] = kerror
}

func (r *OffsetCommitResponse) encode(pe packetEncoder) error {
	if err := pe.putArrayLength(len(r.Errors)); err != nil {
		return err
	}
	for topic, partitions := range r.Errors {
		if err := pe.putString(topic); err != nil {
			return err
		}
		if err := pe.putArrayLength(len(partitions)); err != nil {
			return err
		}
		for partition, kerror := range partitions {
			pe.putInt32(partition)
			pe.putInt16(int16(kerror))
		}
	}
	return nil
}

func (r *OffsetCommitResponse) decode(pd packetDecoder, version int16) (err error) {
	numTopics, err := pd.getArrayLength()
	if err != nil || numTopics == 0 {
		return err
	}

	r.Errors = make(map[string]map[int32]KError, numTopics)
	for i := 0; i < numTopics; i++ {
		name, err := pd.getString()
		if err != nil {
			return err
		}

		numErrors, err := pd.getArrayLength()
		if err != nil {
			return err
		}

		r.Errors[name] = make(map[int32]KError, numErrors)

		for j := 0; j < numErrors; j++ {
			id, err := pd.getInt32()
			if err != nil {
				return err
			}

			tmp, err := pd.getInt16()
			if err != nil {
				return err
			}
			r.Errors[name][id] = KError(tmp)
		}
	}

	return nil
}

func (r *OffsetCommitResponse) key() int16 {
	return 8
}

func (r *OffsetCommitResponse) version() int16 {
	return 0
}
