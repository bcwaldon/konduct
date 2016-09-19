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
// +build !windows

package utils

import (
	"io/ioutil"
	"strconv"
	"syscall"
)

func CloseExecFrom(minFd int) error {
	fdList, err := ioutil.ReadDir("/proc/self/fd")
	if err != nil {
		return err
	}
	for _, fi := range fdList {
		fd, err := strconv.Atoi(fi.Name())
		if err != nil {
			// ignore non-numeric file names
			continue
		}

		if fd < minFd {
			// ignore descriptors lower than our specified minimum
			continue
		}

		// intentionally ignore errors from syscall.CloseOnExec
		syscall.CloseOnExec(fd)
		// the cases where this might fail are basically file descriptors that have already been closed (including and especially the one that was created when ioutil.ReadDir did the "opendir" syscall)
	}
	return nil
}
