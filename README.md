# Leaky bucket rate limiter

A Go implementation of the leaky-bucket rate limit algorithm. 

## Install
`go get github.com/twiny/ratelimit`

## API

```go
Take() time.Time
Rate() int
Duration() time.Duration
String() string
```

## Example
```go
package main

import (
	"fmt"
	"time"

	"github.com/twiny/ratelimit"
)

// main
func main() {
	limiter := ratelimit.NewLimiter(10, 1*time.Second) // 10 request per second

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := limiter.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}

// output:
// 0 239ns
// 1 100ms
// 2 100ms
// 3 100ms
// 4 100ms
// 5 100ms
// 6 100ms
// 7 100ms
// 8 100ms
// 9 100ms
```