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
package transport

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func TLSDialer(tlscfg *tls.Config, timeout time.Duration, forward Dialer) Dialer {
	if tlscfg == nil {
		return forward
	}

	return DialerFunc(func(network, address string) (net.Conn, error) {
		switch network {
		case "tcp", "tcp4", "tcp6":
		default:
			return nil, fmt.Errorf("unsupported network type %v", network)
		}

		socket, err := forward.Dial(network, address)
		if err != nil {
			return nil, err
		}

		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}

		tlscfg.ServerName = host
		conn := tls.Client(socket, tlscfg)
		if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			_ = conn.Close()
			return nil, err
		}
		if err := conn.Handshake(); err != nil {
			_ = conn.Close()
			return nil, err
		}

		return conn, nil
	})
}
