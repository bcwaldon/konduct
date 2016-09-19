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
// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mergo merges same-type structs and maps by setting default values in zero-value fields.

Mergo won't merge unexported (private) fields but will do recursively any exported one. It also won't merge structs inside maps (because they are not addressable using Go reflection).

Usage

From my own work-in-progress project:

	type networkConfig struct {
		Protocol string
		Address string
		ServerType string `json: "server_type"`
		Port uint16
	}

	type FssnConfig struct {
		Network networkConfig
	}

	var fssnDefault = FssnConfig {
		networkConfig {
			"tcp",
			"127.0.0.1",
			"http",
			31560,
		},
	}

	// Inside a function [...]

	if err := mergo.Merge(&config, fssnDefault); err != nil {
		log.Fatal(err)
	}

	// More code [...]

*/
package mergo
