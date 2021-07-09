// Copyright 2019 John Papandriopoulos.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concurrent_test

import (
	"sync"
	"testing"

	"go.jpap.org/concurrent"
)

func TestRunGrouped(t *testing.T) {
	tests := []struct {
		count   int
		threads int
	}{
		{6, 1},
		{7, 2},
		{3, 0},
		{7, 0},
		{20, 0},
		{7, 10},
		{11, 10},
		{3, 10},
	}

	for _, tc := range tests {
		var mux sync.Mutex
		x := make(map[int]bool)
		concurrent.RunGrouped(tc.count, tc.threads, func(m, n int) {
			mux.Lock()
			for i := m; i < n; i++ {
				x[i] = true
			}
			mux.Unlock()
		})
		if c := len(x); c != tc.count {
			t.Errorf("failed: got %d, expected %d (threads: %d): %v", c, tc.count, tc.threads, x)
		}
	}
}
