// Copyright 2021 John Papandriopoulos.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package concurrent

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// RunSweepErr will use at most maxThreads (or equal to the number of CPUs on
// the system if zero), to run func f concurrently, returning the first error
// received.  If an error is reported, some func f may not be executed.
func RunSweepErr(count, maxThreads int, f func(index int) error) error {
	if count == 0 {
		return nil
	}
	if maxThreads == 0 {
		maxThreads = runtime.NumCPU()
	}
	if maxThreads > count {
		maxThreads = count
	}

	if maxThreads == 1 {
		for i := 0; i < count; i++ {
			if err := f(i); err != nil {
				return err
			}
		}
		return nil
	}

	q := make(chan struct{}, maxThreads)

	nerr := uint32(0)
	var firstErr error
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		if atomic.LoadUint32(&nerr) > 0 {
			break
		}

		// Block if we exceed maxThreads
		q <- struct{}{}
		wg.Add(1)

		// Run job
		go func(i int) {
			err := f(i)
			if err != nil {
				if atomic.AddUint32(&nerr, 1) == 1 {
					firstErr = err
				}
			}
			// We're done
			wg.Done()
			<-q
		}(i)
	}

	wg.Wait()
	close(q)

	return firstErr
}

// RunSweep is like RunSweepErr, but without errors.
func RunSweep(count, maxThreads int, f func(index int)) {
	RunSweepErr(count, maxThreads, func(i int) error {
		f(i)
		return nil
	})
}
