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
package fileout

import (
	"fmt"

	"github.com/elastic/beats/libbeat/logp"
)

type config struct {
	Index         string `config:"index"`
	Path          string `config:"path"`
	Filename      string `config:"filename"`
	RotateEveryKb int    `config:"rotate_every_kb" validate:"min=1"`
	NumberOfFiles int    `config:"number_of_files"`
}

var (
	defaultConfig = config{
		NumberOfFiles: 7,
		RotateEveryKb: 10 * 1024,
	}
)

func (c *config) Validate() error {
	if c.Filename == "" && c.Index == "" {
		return fmt.Errorf("File logging requires filename or index being set.")
	}

	if c.NumberOfFiles < 2 || c.NumberOfFiles > logp.RotatorMaxFiles {
		return fmt.Errorf("The number_of_files to keep should be between 2 and %v",
			logp.RotatorMaxFiles)
	}

	return nil
}
