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
package publisher

import "github.com/elastic/beats/libbeat/common/op"

// ClientOption allows API users to set additional options when publishing events.
type ClientOption func(option Context) Context

// Guaranteed option will retry publishing the event, until send attempt have
// been ACKed by output plugin.
func Guaranteed(o Context) Context {
	o.Guaranteed = true
	return o
}

// Sync option will block the event publisher until an event has been ACKed by
// the output plugin or failed.
func Sync(o Context) Context {
	o.Sync = true
	return o
}

func Signal(signaler op.Signaler) ClientOption {
	return func(ctx Context) Context {
		if ctx.Signal == nil {
			ctx.Signal = signaler
		} else {
			ctx.Signal = op.CombineSignalers(ctx.Signal, signaler)
		}
		return ctx
	}
}
