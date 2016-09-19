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

type SaslHandshakeResponse struct {
	Err               KError
	EnabledMechanisms []string
}

func (r *SaslHandshakeResponse) encode(pe packetEncoder) error {
	pe.putInt16(int16(r.Err))
	return pe.putStringArray(r.EnabledMechanisms)
}

func (r *SaslHandshakeResponse) decode(pd packetDecoder, version int16) error {
	if kerr, err := pd.getInt16(); err != nil {
		return err
	} else {
		r.Err = KError(kerr)
	}

	var err error
	if r.EnabledMechanisms, err = pd.getStringArray(); err != nil {
		return err
	}

	return nil
}

func (r *SaslHandshakeResponse) key() int16 {
	return 17
}

func (r *SaslHandshakeResponse) version() int16 {
	return 0
}
