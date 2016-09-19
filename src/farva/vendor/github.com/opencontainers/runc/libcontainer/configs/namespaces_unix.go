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
// +build linux freebsd

package configs

import (
	"fmt"
	"os"
	"sync"
)

const (
	NEWNET  NamespaceType = "NEWNET"
	NEWPID  NamespaceType = "NEWPID"
	NEWNS   NamespaceType = "NEWNS"
	NEWUTS  NamespaceType = "NEWUTS"
	NEWIPC  NamespaceType = "NEWIPC"
	NEWUSER NamespaceType = "NEWUSER"
)

var (
	nsLock              sync.Mutex
	supportedNamespaces = make(map[NamespaceType]bool)
)

// nsToFile converts the namespace type to its filename
func nsToFile(ns NamespaceType) string {
	switch ns {
	case NEWNET:
		return "net"
	case NEWNS:
		return "mnt"
	case NEWPID:
		return "pid"
	case NEWIPC:
		return "ipc"
	case NEWUSER:
		return "user"
	case NEWUTS:
		return "uts"
	}
	return ""
}

// IsNamespaceSupported returns whether a namespace is available or
// not
func IsNamespaceSupported(ns NamespaceType) bool {
	nsLock.Lock()
	defer nsLock.Unlock()
	supported, ok := supportedNamespaces[ns]
	if ok {
		return supported
	}
	nsFile := nsToFile(ns)
	// if the namespace type is unknown, just return false
	if nsFile == "" {
		return false
	}
	_, err := os.Stat(fmt.Sprintf("/proc/self/ns/%s", nsFile))
	// a namespace is supported if it exists and we have permissions to read it
	supported = err == nil
	supportedNamespaces[ns] = supported
	return supported
}

func NamespaceTypes() []NamespaceType {
	return []NamespaceType{
		NEWNET,
		NEWPID,
		NEWNS,
		NEWUTS,
		NEWIPC,
		NEWUSER,
	}
}

// Namespace defines configuration for each namespace.  It specifies an
// alternate path that is able to be joined via setns.
type Namespace struct {
	Type NamespaceType `json:"type"`
	Path string        `json:"path"`
}

func (n *Namespace) GetPath(pid int) string {
	if n.Path != "" {
		return n.Path
	}
	return fmt.Sprintf("/proc/%d/ns/%s", pid, nsToFile(n.Type))
}

func (n *Namespaces) Remove(t NamespaceType) bool {
	i := n.index(t)
	if i == -1 {
		return false
	}
	*n = append((*n)[:i], (*n)[i+1:]...)
	return true
}

func (n *Namespaces) Add(t NamespaceType, path string) {
	i := n.index(t)
	if i == -1 {
		*n = append(*n, Namespace{Type: t, Path: path})
		return
	}
	(*n)[i].Path = path
}

func (n *Namespaces) index(t NamespaceType) int {
	for i, ns := range *n {
		if ns.Type == t {
			return i
		}
	}
	return -1
}

func (n *Namespaces) Contains(t NamespaceType) bool {
	return n.index(t) != -1
}

func (n *Namespaces) PathOf(t NamespaceType) string {
	i := n.index(t)
	if i == -1 {
		return ""
	}
	return (*n)[i].Path
}
