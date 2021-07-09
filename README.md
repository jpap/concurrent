<!-- DO NOT EDIT. -->
<!-- Automatically generated with https://go.jpap.org/godoc-readme-gen -->

# Limited Width Concurrency for Go [![GoDoc](https://pkg.go.dev/badge/go.jpap.org/concurrent.svg)](https://pkg.go.dev/go.jpap.org/concurrent)

# Import

```go
import "go.jpap.org/concurrent"
```
# Overview

Package concurrent makes it easy to execute a list of jobs concurrently with
a simple closure while using a finite number of goroutines (concurrency
width).

Three broad patterns are supported and described below.  For each, the
package user can easily cap the maximum concurrency width, that is clipped
to the maximum number of CPUs on the system in all cases.

This package was created because the author kept finding the need to
implement these patterns over and over.

## Grouped Execution Pattern
The idea here is to take `n` jobs, split them into groups ("batches" or
"chunks"). Each invocation `i` of the closure is given an index range `[m_i,
n_i)` that specifies a non-overlapping group.  The union of all invocations
covers the index range `[0, n)`.

If one of the invocations returns an error, the first error received is
returned, but all invocations are executed before returning.

```go
func GroupedExample() {
  num := int[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
  var sum uint32
  concurrent.RunGrouped(len(num), 0, func(m, n int) {
    localSum := 0
    for j := m; j < n; j++ {
      localSum += num[j]
    }
    atomic.AddUint32(&sum, localSum)
  })
}
```

## Sweep Execution Pattern
The idea here is similar to the Grouped pattern, except that the closure is
invoked once per job, and if any invocation returns an error, no more
invocations are scheduled.  This allows errors to "short circuit" execution
and return earlier than Grouped equivalent.

Like the Grouped pattern, the concurrency width is limited, to reduce
goroutine overheads.

```go
func SweepExample() {
  // A trivial example, to contrast to the Grouped pattern example above.
  // You would almost surely not implement a sum in this manner. ;-)
  num := int[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
  var sum uint32
  concurrent.RunSweep(len(num), 0, func(j int) {
    atomic.AddUint32(&sum, num[j])
  })
}
```

## Runner Pattern
A `Runner` is also provided, that allows jobs to be scheduled without having
to know how many jobs are required up-front.  The implementation of the
Sweep pattern uses this functionality.




