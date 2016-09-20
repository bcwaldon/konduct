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

package fs

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
)

type PidsGroup struct {
}

func (s *PidsGroup) Name() string {
	return "pids"
}

func (s *PidsGroup) Apply(d *cgroupData) error {
	_, err := d.join("pids")
	if err != nil && !cgroups.IsNotFound(err) {
		return err
	}
	return nil
}

func (s *PidsGroup) Set(path string, cgroup *configs.Cgroup) error {
	if cgroup.Resources.PidsLimit != 0 {
		// "max" is the fallback value.
		limit := "max"

		if cgroup.Resources.PidsLimit > 0 {
			limit = strconv.FormatInt(cgroup.Resources.PidsLimit, 10)
		}

		if err := writeFile(path, "pids.max", limit); err != nil {
			return err
		}
	}

	return nil
}

func (s *PidsGroup) Remove(d *cgroupData) error {
	return removePath(d.path("pids"))
}

func (s *PidsGroup) GetStats(path string, stats *cgroups.Stats) error {
	current, err := getCgroupParamUint(path, "pids.current")
	if err != nil {
		return fmt.Errorf("failed to parse pids.current - %s", err)
	}

	maxString, err := getCgroupParamString(path, "pids.max")
	if err != nil {
		return fmt.Errorf("failed to parse pids.max - %s", err)
	}

	// Default if pids.max == "max" is 0 -- which represents "no limit".
	var max uint64
	if maxString != "max" {
		max, err = parseUint(maxString, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse pids.max - unable to parse %q as a uint from Cgroup file %q", maxString, filepath.Join(path, "pids.max"))
		}
	}

	stats.PidsStats.Current = current
	stats.PidsStats.Limit = max
	return nil
}
