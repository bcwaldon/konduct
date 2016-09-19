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
package system

import (
	"fmt"
	"runtime"
	"syscall"
)

// Via http://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/commit/?id=7b21fddd087678a70ad64afc0f632e0f1071b092
//
// We need different setns values for the different platforms and arch
// We are declaring the macro here because the SETNS syscall does not exist in th stdlib
var setNsMap = map[string]uintptr{
	"linux/386":     346,
	"linux/arm64":   268,
	"linux/amd64":   308,
	"linux/arm":     375,
	"linux/ppc":     350,
	"linux/ppc64":   350,
	"linux/ppc64le": 350,
	"linux/s390x":   339,
}

var sysSetns = setNsMap[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]

func SysSetns() uint32 {
	return uint32(sysSetns)
}

func Setns(fd uintptr, flags uintptr) error {
	ns, exists := setNsMap[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !exists {
		return fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	_, _, err := syscall.RawSyscall(ns, fd, flags, 0)
	if err != 0 {
		return err
	}
	return nil
}
