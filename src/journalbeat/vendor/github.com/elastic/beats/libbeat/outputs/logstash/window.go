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
package logstash

import (
	"math"
	"sync/atomic"
)

type window struct {
	windowSize      int32
	maxOkWindowSize int // max window size sending was successful for
	maxWindowSize   int
}

func (w *window) init(start, max int) {
	*w = window{
		windowSize:    int32(start),
		maxWindowSize: max,
	}
}

func (w *window) get() int {
	return int(atomic.LoadInt32(&w.windowSize))
}

// Increase window size by factor 1.5 until max window size
// (window size grows exponentially)
// TODO: use duration until ACK to estimate an ok max window size value
func (w *window) tryGrowWindow(batchSize int) {
	windowSize := w.get()

	if windowSize <= batchSize {
		if w.maxOkWindowSize < windowSize {
			debug("update max ok window size: %v < %v",
				w.maxOkWindowSize, w.windowSize)
			w.maxOkWindowSize = windowSize

			newWindowSize := int(math.Ceil(1.5 * float64(windowSize)))
			debug("increase window size to: %v", newWindowSize)

			if windowSize <= batchSize && batchSize < newWindowSize {
				debug("set to batchSize: %v", batchSize)
				newWindowSize = batchSize
			}
			if newWindowSize > w.maxWindowSize {
				debug("set to max window size: %v", w.maxWindowSize)
				newWindowSize = int(w.maxWindowSize)
			}

			windowSize = newWindowSize
		} else if windowSize < w.maxOkWindowSize {
			debug("update current window size: %v", w.windowSize)

			windowSize = int(math.Ceil(1.5 * float64(windowSize)))
			if windowSize > w.maxOkWindowSize {
				debug("set to max ok window size: %v", w.maxOkWindowSize)
				windowSize = w.maxOkWindowSize
			}
		}

		atomic.StoreInt32(&w.windowSize, int32(windowSize))
	}
}

func (w *window) shrinkWindow() {
	windowSize := w.get()
	orig := windowSize

	windowSize = windowSize / 2
	if windowSize < minWindowSize {
		windowSize = minWindowSize
		if windowSize == orig {
			return
		}
	}

	atomic.StoreInt32(&w.windowSize, int32(windowSize))
}
