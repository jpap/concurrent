// Copyright 2019 John Papandriopoulos.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package concurrent_test

import (
	"sync"
	"testing"

	"github.com/yourbasic/bit"

	"go.jpap.org/concurrent"
)

func TestRunConcurrently(t *testing.T) {
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
		b := bit.New()
		concurrent.Run(tc.count, tc.threads, func(m, n int) {
			mux.Lock()
			for i := m; i < n; i++ {
				b.Add(i)
			}
			mux.Unlock()
		})
		if c := b.Size(); c != tc.count {
			t.Errorf("failed: got %d, expected %d (threads: %d): %v", c, tc.count, tc.threads, b)
		}
	}
}
