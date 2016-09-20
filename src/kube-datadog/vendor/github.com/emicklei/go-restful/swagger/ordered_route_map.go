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
package swagger

// Copyright 2015 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "github.com/emicklei/go-restful"

type orderedRouteMap struct {
	elements map[string][]restful.Route
	keys     []string
}

func newOrderedRouteMap() *orderedRouteMap {
	return &orderedRouteMap{
		elements: map[string][]restful.Route{},
		keys:     []string{},
	}
}

func (o *orderedRouteMap) Add(key string, route restful.Route) {
	routes, ok := o.elements[key]
	if ok {
		routes = append(routes, route)
		o.elements[key] = routes
		return
	}
	o.elements[key] = []restful.Route{route}
	o.keys = append(o.keys, key)
}

func (o *orderedRouteMap) Do(block func(key string, routes []restful.Route)) {
	for _, k := range o.keys {
		block(k, o.elements[k])
	}
}
