# concurrent

```
import "go.jpap.org/concurrent"

...

  num := int[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
  var sum uint32
  concurrent.Run(len(num), 0, func(m, n int) {
    lsum := 0
    for i := m; i < n; i++ {
      lsum += num[i]
    }
    atomic.AddUint32(&sum, lsum)
  })
```

This Go package makes it easy to process a list of things concurrently, using
a finite number of goroutines, using a simple closure.

## License

MIT, see the `LICENSE.md` file.
