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
package humanize

import "strconv"

func stripTrailingZeros(s string) string {
	offset := len(s) - 1
	for offset > 0 {
		if s[offset] == '.' {
			offset--
			break
		}
		if s[offset] != '0' {
			break
		}
		offset--
	}
	return s[:offset+1]
}

// Ftoa converts a float to a string with no trailing zeros.
func Ftoa(num float64) string {
	return stripTrailingZeros(strconv.FormatFloat(num, 'f', 6, 64))
}
