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

import "fmt"

type responseHeader struct {
	length        int32
	correlationID int32
}

func (r *responseHeader) decode(pd packetDecoder) (err error) {
	r.length, err = pd.getInt32()
	if err != nil {
		return err
	}
	if r.length <= 4 || r.length > MaxResponseSize {
		return PacketDecodingError{fmt.Sprintf("message of length %d too large or too small", r.length)}
	}

	r.correlationID, err = pd.getInt32()
	return err
}
