// Copyright 2019 John Papandriopoulos.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package concurrent makes it easy to process a list of things
// concurrently using a simple closure.
package concurrent // import "go.jpap.org/concurrent"

import (
	"runtime"
	"sync"
)

// Run the given func f concurrently using up to the specified number of
// maxThreads, or equal to the number of CPUs on the system if passed zero.
// The number of jobs is given by count, and the range of jobs [m, n) are
// passed to the callback f.
func Run(count, maxThreads int, f func(m, n int)) {
	if count == 0 {
		return
	}
	if maxThreads == 0 {
		maxThreads = runtime.NumCPU()
	}
	if maxThreads > count {
		maxThreads = count
	}

	if maxThreads == 1 {
		f(0, count)
		return
	}

	countPerThread := count / maxThreads
	m := 0
	n := countPerThread

	q := make(chan struct{}, maxThreads)
	var wg sync.WaitGroup
	for i := 0; i < count; i += countPerThread {
		// Queue job, blocking if we exceed maxThreads
		q <- struct{}{}
		// Queue job
		wg.Add(1)
		go func(a, b int) {
			f(a, b)
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
	// Wait for all outstanding jobs to complete
	wg.Wait()
	close(q)
}
