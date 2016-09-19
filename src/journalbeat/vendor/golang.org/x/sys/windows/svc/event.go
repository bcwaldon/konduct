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
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package svc

import (
	"errors"

	"golang.org/x/sys/windows"
)

// event represents auto-reset, initially non-signaled Windows event.
// It is used to communicate between go and asm parts of this package.
type event struct {
	h windows.Handle
}

func newEvent() (*event, error) {
	h, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return nil, err
	}
	return &event{h: h}, nil
}

func (e *event) Close() error {
	return windows.CloseHandle(e.h)
}

func (e *event) Set() error {
	return windows.SetEvent(e.h)
}

func (e *event) Wait() error {
	s, err := windows.WaitForSingleObject(e.h, windows.INFINITE)
	switch s {
	case windows.WAIT_OBJECT_0:
		break
	case windows.WAIT_FAILED:
		return err
	default:
		return errors.New("unexpected result from WaitForSingleObject")
	}
	return nil
}
