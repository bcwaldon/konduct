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
// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package windows

type WSAData struct {
	Version      uint16
	HighVersion  uint16
	MaxSockets   uint16
	MaxUdpDg     uint16
	VendorInfo   *byte
	Description  [WSADESCRIPTION_LEN + 1]byte
	SystemStatus [WSASYS_STATUS_LEN + 1]byte
}

type Servent struct {
	Name    *byte
	Aliases **byte
	Proto   *byte
	Port    uint16
}
