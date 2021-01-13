// Copyright 2021 John Papandriopoulos.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package concurrent

// RunSweepErr will use at most maxThreads (or equal to the number of CPUs on
// the system if zero), to run func f concurrently, returning the first error
// received.  If an error is reported, some func f may not be executed.
func RunSweepErr(count, maxThreads int, f func(index int) error) error {
	if count == 0 {
		return nil
	}

	jr := NewRunner(maxThreads)
	for i := 0; i < count; i++ {
		if jr.Failures() > 0 {
			break
		}
		index := i
		jr.RunErr(func() error {
			return f(index)
		})
	}
	if jr.Finish() != count {
		return jr.Errors()[0]
	}
	return nil
}

// RunSweep is like RunSweepErr, but without errors.
func RunSweep(count, maxThreads int, f func(index int)) {
	RunSweepErr(count, maxThreads, func(i int) error {
		f(i)
		return nil
	})
}
