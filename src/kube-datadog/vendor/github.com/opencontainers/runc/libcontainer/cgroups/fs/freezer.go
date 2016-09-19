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
// +build linux

package fs

import (
	"fmt"
	"strings"
	"time"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
)

type FreezerGroup struct {
}

func (s *FreezerGroup) Name() string {
	return "freezer"
}

func (s *FreezerGroup) Apply(d *cgroupData) error {
	_, err := d.join("freezer")
	if err != nil && !cgroups.IsNotFound(err) {
		return err
	}
	return nil
}

func (s *FreezerGroup) Set(path string, cgroup *configs.Cgroup) error {
	switch cgroup.Resources.Freezer {
	case configs.Frozen, configs.Thawed:
		if err := writeFile(path, "freezer.state", string(cgroup.Resources.Freezer)); err != nil {
			return err
		}

		for {
			state, err := readFile(path, "freezer.state")
			if err != nil {
				return err
			}
			if strings.TrimSpace(state) == string(cgroup.Resources.Freezer) {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}
	case configs.Undefined:
		return nil
	default:
		return fmt.Errorf("Invalid argument '%s' to freezer.state", string(cgroup.Resources.Freezer))
	}

	return nil
}

func (s *FreezerGroup) Remove(d *cgroupData) error {
	return removePath(d.path("freezer"))
}

func (s *FreezerGroup) GetStats(path string, stats *cgroups.Stats) error {
	return nil
}
