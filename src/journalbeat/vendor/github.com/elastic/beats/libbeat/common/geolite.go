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
package common

import (
	"os"
	"path/filepath"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/nranchev/go-libGeoIP"
)

// Geoip represents a string slice of GeoIP paths
type Geoip struct {
	Paths *[]string
}

func LoadGeoIPData(config Geoip) *libgeo.GeoIP {
	geoipPaths := []string{}

	if config.Paths != nil {
		geoipPaths = *config.Paths
	}
	if len(geoipPaths) == 0 {
		// disabled
		return nil
	} else {
		logp.Warn("GeoIP lookup support is deprecated and will be removed in version 6.0.")
	}

	// look for the first existing path
	var geoipPath string
	for _, path := range geoipPaths {
		fi, err := os.Lstat(path)
		if err != nil {
			logp.Err("GeoIP path could not be loaded: %s", path)
			continue
		}

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			// follow symlink
			geoipPath, err = filepath.EvalSymlinks(path)
			if err != nil {
				logp.Warn("Could not load GeoIP data: %s", err.Error())
				return nil
			}
		} else {
			geoipPath = path
		}
		break
	}

	if len(geoipPath) == 0 {
		logp.Warn("Couldn't load GeoIP database")
		return nil
	}

	geoLite, err := libgeo.Load(geoipPath)
	if err != nil {
		logp.Warn("Could not load GeoIP data: %s", err.Error())
	}

	logp.Info("Loaded GeoIP data from: %s", geoipPath)
	return geoLite
}
