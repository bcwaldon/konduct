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

// +build linux

package cgroups

import (
	"fmt"

	"github.com/opencontainers/runc/libcontainer/configs"
)

type Manager interface {
	// Applies cgroup configuration to the process with the specified pid
	Apply(pid int) error

	// Returns the PIDs inside the cgroup set
	GetPids() ([]int, error)

	// Returns the PIDs inside the cgroup set & all sub-cgroups
	GetAllPids() ([]int, error)

	// Returns statistics for the cgroup set
	GetStats() (*Stats, error)

	// Toggles the freezer cgroup according with specified state
	Freeze(state configs.FreezerState) error

	// Destroys the cgroup set
	Destroy() error

	// NewCgroupManager() and LoadCgroupManager() require following attributes:
	// 	Paths   map[string]string
	// 	Cgroups *cgroups.Cgroup
	// Paths maps cgroup subsystem to path at which it is mounted.
	// Cgroups specifies specific cgroup settings for the various subsystems

	// Returns cgroup paths to save in a state file and to be able to
	// restore the object later.
	GetPaths() map[string]string

	// Set the cgroup as configured.
	Set(container *configs.Config) error
}

type NotFoundError struct {
	Subsystem string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("mountpoint for %s not found", e.Subsystem)
}

func NewNotFoundError(sub string) error {
	return &NotFoundError{
		Subsystem: sub,
	}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*NotFoundError)
	return ok
}
