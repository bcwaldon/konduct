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
	"expvar"
	"net"
)

type IOStats struct {
	Read, Write, ReadErrors, WriteErrors *expvar.Int
}

type statsConn struct {
	net.Conn
	stats *IOStats
}

func StatsDialer(d Dialer, s *IOStats) Dialer {
	return ConnWrapper(d, func(c net.Conn) net.Conn {
		return &statsConn{c, s}
	})
}

func (s *statsConn) Read(b []byte) (int, error) {
	n, err := s.Conn.Read(b)
	if err != nil {
		s.stats.ReadErrors.Add(1)
	}
	s.stats.Read.Add(int64(n))
	return n, err
}

func (s *statsConn) Write(b []byte) (int, error) {
	n, err := s.Conn.Write(b)
	if err != nil {
		s.stats.WriteErrors.Add(1)
	}
	s.stats.Write.Add(int64(n))
	return n, err
}
