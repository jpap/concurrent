// Copyright 2019 John Papandriopoulos.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package concurrent makes it easy to process a list of things
// concurrently using a simple closure.
package concurrent // import "go.jpap.org/concurrent"

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// RunGroupedErr the given func f concurrently using up to the specified number of
// maxThreads, or equal to the number of CPUs on the system if passed zero.
// The number of jobs is given by count, and the range of jobs [m, n) are
// passed to the callback f.
func RunGroupedErr(count, maxThreads int, f func(m, n int) error) error {
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
		return f(0, count)
	}

	countPerThread := count / maxThreads
	m := 0
	n := countPerThread

	q := make(chan struct{}, maxThreads)

	nerr := uint32(0)
	var firstErr error
	var wg sync.WaitGroup

	for i := 0; i < count; i += countPerThread {
		if atomic.LoadUint32(&nerr) > 0 {
			break
		}

		// Block if we exceed maxThreads
		q <- struct{}{}
		wg.Add(1)

		// Run job
		go func(a, b int) {
			err := f(a, b)
			if err != nil {
				if atomic.AddUint32(&nerr, 1) == 1 {
					firstErr = err
				}
			}
			// We're done
			wg.Done()
			<-q
		}(m, n)
		m += countPerThread
		n += countPerThread
		if n > count {
			// Truncate the last job, as required
			n = count
		}
	}

	wg.Wait()
	close(q)

	return firstErr
}

// RunGrouped is like RunGroupedErr but without errors.
func RunGrouped(count, maxThreads int, f func(m, n int)) {
	RunGroupedErr(count, maxThreads, func(m, n int) error {
		f(m, n)
		return nil
	})
}
