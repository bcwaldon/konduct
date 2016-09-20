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
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// look in /proc to find the process start time so that we can verify
// that this pid has started after ourself
func GetProcessStartTime(pid int) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(data), " ")
	// the starttime is located at pos 22
	// from the man page
	//
	// starttime %llu (was %lu before Linux 2.6)
	// (22)  The  time the process started after system boot.  In kernels before Linux 2.6, this
	// value was expressed in jiffies.  Since Linux 2.6, the value is expressed in  clock  ticks
	// (divide by sysconf(_SC_CLK_TCK)).
	return parts[22-1], nil // starts at 1
}
