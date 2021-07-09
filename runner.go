// Copyright 2021 John Papandriopoulos.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// Runner runs jobs with a specified maximum limit on concurrency.
type Runner struct {
	sync.Mutex

	width     int
	nsuccess  uint32
	nfailures uint32
	ch        chan struct{}
	errors    []error
	wg        sync.WaitGroup
}

// NewRunner returns a new Runner that executes jobs concurrently with at most
// n jobs running at any one time.
func NewRunner(n int) *Runner {
	if n == 0 {
		n = runtime.NumCPU()
	}
	return &Runner{
		width: n,
		ch:    make(chan struct{}, n),
	}
}

// Errors returns all of the errors reported by jobs.
// The order is given by job completion, not submission.
func (jr *Runner) Errors() []error {
	jr.Lock()
	defer jr.Unlock()
	return append(([]error)(nil), jr.errors...) // copy
}

// Failures returns the number of errors reported by jobs so far.
// This is useful for callers to fail-fast and stop submitting
// jobs if an error has been reported.
func (jr *Runner) Failures() int {
	return int(atomic.LoadUint32(&jr.nfailures))
}

// RunErr runs a job that reports an error if it fails.
func (jr *Runner) RunErr(job func() error) {
	if jr.width == 1 {
		if err := job(); err != nil {
			jr.errors = append(jr.errors, err)
			jr.nfailures++
		} else {
			jr.nsuccess++
		}
		return
	}

	// Block if we exceed width
	jr.ch <- struct{}{}
	jr.wg.Add(1)

	// Run job
	go func() {
		if err := job(); err != nil {
			atomic.AddUint32(&jr.nfailures, 1)
			jr.Lock()
			jr.errors = append(jr.errors, err)
			jr.Unlock()
		} else {
			atomic.AddUint32(&jr.nsuccess, 1)
		}
		// Job is done
		jr.wg.Done()
		<-jr.ch
	}()
}

// Run is like RunErr, but the job func does not need to return an error.
func (jr *Runner) Run(job func()) {
	jr.RunErr(func() error { job(); return nil })
}

// Finish waits for all executing jobs to complete, and returns the number
// of successful jobs completed.
func (jr *Runner) Finish() int {
	jr.wg.Wait()
	close(jr.ch)
	return int(jr.nsuccess)
}
