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
// Copyright 2011 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"net"
	"sync"
)

var (
	nodeMu     sync.Mutex
	interfaces []net.Interface // cached list of interfaces
	ifname     string          // name of interface being used
	nodeID     []byte          // hardware for version 1 UUIDs
)

// NodeInterface returns the name of the interface from which the NodeID was
// derived.  The interface "user" is returned if the NodeID was set by
// SetNodeID.
func NodeInterface() string {
	defer nodeMu.Unlock()
	nodeMu.Lock()
	return ifname
}

// SetNodeInterface selects the hardware address to be used for Version 1 UUIDs.
// If name is "" then the first usable interface found will be used or a random
// Node ID will be generated.  If a named interface cannot be found then false
// is returned.
//
// SetNodeInterface never fails when name is "".
func SetNodeInterface(name string) bool {
	defer nodeMu.Unlock()
	nodeMu.Lock()
	return setNodeInterface(name)
}

func setNodeInterface(name string) bool {
	if interfaces == nil {
		var err error
		interfaces, err = net.Interfaces()
		if err != nil && name != "" {
			return false
		}
	}

	for _, ifs := range interfaces {
		if len(ifs.HardwareAddr) >= 6 && (name == "" || name == ifs.Name) {
			if setNodeID(ifs.HardwareAddr) {
				ifname = ifs.Name
				return true
			}
		}
	}

	// We found no interfaces with a valid hardware address.  If name
	// does not specify a specific interface generate a random Node ID
	// (section 4.1.6)
	if name == "" {
		if nodeID == nil {
			nodeID = make([]byte, 6)
		}
		randomBits(nodeID)
		return true
	}
	return false
}

// NodeID returns a slice of a copy of the current Node ID, setting the Node ID
// if not already set.
func NodeID() []byte {
	defer nodeMu.Unlock()
	nodeMu.Lock()
	if nodeID == nil {
		setNodeInterface("")
	}
	nid := make([]byte, 6)
	copy(nid, nodeID)
	return nid
}

// SetNodeID sets the Node ID to be used for Version 1 UUIDs.  The first 6 bytes
// of id are used.  If id is less than 6 bytes then false is returned and the
// Node ID is not set.
func SetNodeID(id []byte) bool {
	defer nodeMu.Unlock()
	nodeMu.Lock()
	if setNodeID(id) {
		ifname = "user"
		return true
	}
	return false
}

func setNodeID(id []byte) bool {
	if len(id) < 6 {
		return false
	}
	if nodeID == nil {
		nodeID = make([]byte, 6)
	}
	copy(nodeID, id)
	return true
}

// NodeID returns the 6 byte node id encoded in uuid.  It returns nil if uuid is
// not valid.  The NodeID is only well defined for version 1 and 2 UUIDs.
func (uuid UUID) NodeID() []byte {
	if len(uuid) != 16 {
		return nil
	}
	node := make([]byte, 6)
	copy(node, uuid[10:])
	return node
}
