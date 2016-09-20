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
package configs

import "fmt"

// blockIODevice holds major:minor format supported in blkio cgroup
type blockIODevice struct {
	// Major is the device's major number
	Major int64 `json:"major"`
	// Minor is the device's minor number
	Minor int64 `json:"minor"`
}

// WeightDevice struct holds a `major:minor weight`|`major:minor leaf_weight` pair
type WeightDevice struct {
	blockIODevice
	// Weight is the bandwidth rate for the device, range is from 10 to 1000
	Weight uint16 `json:"weight"`
	// LeafWeight is the bandwidth rate for the device while competing with the cgroup's child cgroups, range is from 10 to 1000, cfq scheduler only
	LeafWeight uint16 `json:"leafWeight"`
}

// NewWeightDevice returns a configured WeightDevice pointer
func NewWeightDevice(major, minor int64, weight, leafWeight uint16) *WeightDevice {
	wd := &WeightDevice{}
	wd.Major = major
	wd.Minor = minor
	wd.Weight = weight
	wd.LeafWeight = leafWeight
	return wd
}

// WeightString formats the struct to be writable to the cgroup specific file
func (wd *WeightDevice) WeightString() string {
	return fmt.Sprintf("%d:%d %d", wd.Major, wd.Minor, wd.Weight)
}

// LeafWeightString formats the struct to be writable to the cgroup specific file
func (wd *WeightDevice) LeafWeightString() string {
	return fmt.Sprintf("%d:%d %d", wd.Major, wd.Minor, wd.LeafWeight)
}

// ThrottleDevice struct holds a `major:minor rate_per_second` pair
type ThrottleDevice struct {
	blockIODevice
	// Rate is the IO rate limit per cgroup per device
	Rate uint64 `json:"rate"`
}

// NewThrottleDevice returns a configured ThrottleDevice pointer
func NewThrottleDevice(major, minor int64, rate uint64) *ThrottleDevice {
	td := &ThrottleDevice{}
	td.Major = major
	td.Minor = minor
	td.Rate = rate
	return td
}

// String formats the struct to be writable to the cgroup specific file
func (td *ThrottleDevice) String() string {
	return fmt.Sprintf("%d:%d %d", td.Major, td.Minor, td.Rate)
}
