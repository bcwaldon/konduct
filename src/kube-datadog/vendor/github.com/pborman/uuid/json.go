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
// Copyright 2014 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import "errors"

func (u UUID) MarshalJSON() ([]byte, error) {
	if len(u) != 16 {
		return []byte(`""`), nil
	}
	var js [38]byte
	js[0] = '"'
	encodeHex(js[1:], u)
	js[37] = '"'
	return js[:], nil
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	if string(data) == `""` {
		return nil
	}
	if data[0] != '"' {
		return errors.New("invalid UUID format")
	}
	data = data[1 : len(data)-1]
	uu := Parse(string(data))
	if uu == nil {
		return errors.New("invalid UUID format")
	}
	*u = uu
	return nil
}
