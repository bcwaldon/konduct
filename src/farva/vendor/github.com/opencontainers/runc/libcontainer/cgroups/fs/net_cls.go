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
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
)

type NetClsGroup struct {
}

func (s *NetClsGroup) Name() string {
	return "net_cls"
}

func (s *NetClsGroup) Apply(d *cgroupData) error {
	_, err := d.join("net_cls")
	if err != nil && !cgroups.IsNotFound(err) {
		return err
	}
	return nil
}

func (s *NetClsGroup) Set(path string, cgroup *configs.Cgroup) error {
	if cgroup.Resources.NetClsClassid != "" {
		if err := writeFile(path, "net_cls.classid", cgroup.Resources.NetClsClassid); err != nil {
			return err
		}
	}

	return nil
}

func (s *NetClsGroup) Remove(d *cgroupData) error {
	return removePath(d.path("net_cls"))
}

func (s *NetClsGroup) GetStats(path string, stats *cgroups.Stats) error {
	return nil
}
