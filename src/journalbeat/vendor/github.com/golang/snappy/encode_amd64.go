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
// Copyright 2016 The Snappy-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine
// +build gc
// +build !noasm

package snappy

// emitLiteral has the same semantics as in encode_other.go.
//
//go:noescape
func emitLiteral(dst, lit []byte) int

// emitCopy has the same semantics as in encode_other.go.
//
//go:noescape
func emitCopy(dst []byte, offset, length int) int

// extendMatch has the same semantics as in encode_other.go.
//
//go:noescape
func extendMatch(src []byte, i, j int) int

// encodeBlock has the same semantics as in encode_other.go.
//
//go:noescape
func encodeBlock(dst, src []byte) (d int)
