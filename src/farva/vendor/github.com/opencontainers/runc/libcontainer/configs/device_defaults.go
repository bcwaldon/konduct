// Copyright 2016 Planet Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// +build linux freebsd

package configs

var (
	// DefaultSimpleDevices are devices that are to be both allowed and created.
	DefaultSimpleDevices = []*Device{
		// /dev/null and zero
		{
			Path:        "/dev/null",
			Type:        'c',
			Major:       1,
			Minor:       3,
			Permissions: "rwm",
			FileMode:    0666,
		},
		{
			Path:        "/dev/zero",
			Type:        'c',
			Major:       1,
			Minor:       5,
			Permissions: "rwm",
			FileMode:    0666,
		},

		{
			Path:        "/dev/full",
			Type:        'c',
			Major:       1,
			Minor:       7,
			Permissions: "rwm",
			FileMode:    0666,
		},

		// consoles and ttys
		{
			Path:        "/dev/tty",
			Type:        'c',
			Major:       5,
			Minor:       0,
			Permissions: "rwm",
			FileMode:    0666,
		},

		// /dev/urandom,/dev/random
		{
			Path:        "/dev/urandom",
			Type:        'c',
			Major:       1,
			Minor:       9,
			Permissions: "rwm",
			FileMode:    0666,
		},
		{
			Path:        "/dev/random",
			Type:        'c',
			Major:       1,
			Minor:       8,
			Permissions: "rwm",
			FileMode:    0666,
		},
	}
	DefaultAllowedDevices = append([]*Device{
		// allow mknod for any device
		{
			Type:        'c',
			Major:       Wildcard,
			Minor:       Wildcard,
			Permissions: "m",
		},
		{
			Type:        'b',
			Major:       Wildcard,
			Minor:       Wildcard,
			Permissions: "m",
		},

		{
			Path:        "/dev/console",
			Type:        'c',
			Major:       5,
			Minor:       1,
			Permissions: "rwm",
		},
		// /dev/pts/ - pts namespaces are "coming soon"
		{
			Path:        "",
			Type:        'c',
			Major:       136,
			Minor:       Wildcard,
			Permissions: "rwm",
		},
		{
			Path:        "",
			Type:        'c',
			Major:       5,
			Minor:       2,
			Permissions: "rwm",
		},

		// tuntap
		{
			Path:        "",
			Type:        'c',
			Major:       10,
			Minor:       200,
			Permissions: "rwm",
		},
	}, DefaultSimpleDevices...)
	DefaultAutoCreatedDevices = append([]*Device{
		{
			// /dev/fuse is created but not allowed.
			// This is to allow java to work.  Because java
			// Insists on there being a /dev/fuse
			// https://github.com/docker/docker/issues/514
			// https://github.com/docker/docker/issues/2393
			//
			Path:        "/dev/fuse",
			Type:        'c',
			Major:       10,
			Minor:       229,
			Permissions: "rwm",
		},
	}, DefaultSimpleDevices...)
)
