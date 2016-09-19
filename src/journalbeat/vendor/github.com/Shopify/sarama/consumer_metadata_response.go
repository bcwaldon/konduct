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

import (
	"net"
	"strconv"
)

type ConsumerMetadataResponse struct {
	Err             KError
	Coordinator     *Broker
	CoordinatorID   int32  // deprecated: use Coordinator.ID()
	CoordinatorHost string // deprecated: use Coordinator.Addr()
	CoordinatorPort int32  // deprecated: use Coordinator.Addr()
}

func (r *ConsumerMetadataResponse) decode(pd packetDecoder, version int16) (err error) {
	tmp, err := pd.getInt16()
	if err != nil {
		return err
	}
	r.Err = KError(tmp)

	coordinator := new(Broker)
	if err := coordinator.decode(pd); err != nil {
		return err
	}
	if coordinator.addr == ":0" {
		return nil
	}
	r.Coordinator = coordinator

	// this can all go away in 2.0, but we have to fill in deprecated fields to maintain
	// backwards compatibility
	host, portstr, err := net.SplitHostPort(r.Coordinator.Addr())
	if err != nil {
		return err
	}
	port, err := strconv.ParseInt(portstr, 10, 32)
	if err != nil {
		return err
	}
	r.CoordinatorID = r.Coordinator.ID()
	r.CoordinatorHost = host
	r.CoordinatorPort = int32(port)

	return nil
}

func (r *ConsumerMetadataResponse) encode(pe packetEncoder) error {
	pe.putInt16(int16(r.Err))
	if r.Coordinator != nil {
		host, portstr, err := net.SplitHostPort(r.Coordinator.Addr())
		if err != nil {
			return err
		}
		port, err := strconv.ParseInt(portstr, 10, 32)
		if err != nil {
			return err
		}
		pe.putInt32(r.Coordinator.ID())
		if err := pe.putString(host); err != nil {
			return err
		}
		pe.putInt32(int32(port))
		return nil
	}
	pe.putInt32(r.CoordinatorID)
	if err := pe.putString(r.CoordinatorHost); err != nil {
		return err
	}
	pe.putInt32(r.CoordinatorPort)
	return nil
}

func (r *ConsumerMetadataResponse) key() int16 {
	return 10
}

func (r *ConsumerMetadataResponse) version() int16 {
	return 0
}
