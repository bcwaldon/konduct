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
package v2

import "net"

// SyncClient synchronously publishes events to lumberjack endpoint waiting for
// ACK before allowing another send request. The client is not thread-safe.
type SyncClient struct {
	cl *Client
}

func NewSyncClientWith(c *Client) (*SyncClient, error) {
	return &SyncClient{c}, nil
}

func NewSyncClientWithConn(c net.Conn, opts ...Option) (*SyncClient, error) {
	cl, err := NewWithConn(c, opts...)
	if err != nil {
		return nil, err
	}
	return NewSyncClientWith(cl)
}

func SyncDial(address string, opts ...Option) (*SyncClient, error) {
	cl, err := Dial(address, opts...)
	if err != nil {
		return nil, err
	}
	return NewSyncClientWith(cl)
}

func SyncDialWith(
	dial func(network, address string) (net.Conn, error),
	address string,
	opts ...Option,
) (*SyncClient, error) {
	cl, err := DialWith(dial, address, opts...)
	if err != nil {
		return nil, err
	}
	return NewSyncClientWith(cl)
}

func (c *SyncClient) Close() error {
	return c.cl.Close()
}

func (c *SyncClient) Send(data []interface{}) (int, error) {
	if err := c.cl.Send(data); err != nil {
		return 0, err
	}

	seq, err := c.cl.AwaitACK(uint32(len(data)))
	return int(seq), err
}
