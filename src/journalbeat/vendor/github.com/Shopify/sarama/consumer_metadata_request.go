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

type ConsumerMetadataRequest struct {
	ConsumerGroup string
}

func (r *ConsumerMetadataRequest) encode(pe packetEncoder) error {
	return pe.putString(r.ConsumerGroup)
}

func (r *ConsumerMetadataRequest) decode(pd packetDecoder, version int16) (err error) {
	r.ConsumerGroup, err = pd.getString()
	return err
}

func (r *ConsumerMetadataRequest) key() int16 {
	return 10
}

func (r *ConsumerMetadataRequest) version() int16 {
	return 0
}
