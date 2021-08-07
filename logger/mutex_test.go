package logger

import (
	"io"
	"sync"
	"testing"
)

/*
goos: darwin
goarch: amd64
pkg: logger/logger
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkMutex-16          	  567092	      1786 ns/op
BenchmarkMutex4-16         	  141592	      8190 ns/op
BenchmarkMutex8-16         	   72516	     16779 ns/op
BenchmarkMutex16-16        	   35400	     33699 ns/op
BenchmarkMutex32-16        	   17161	     68696 ns/op
BenchmarkMutexNoop-16      	37217535	        30.94 ns/op
BenchmarkMutex4Noop-16     	 4726400	       246.3 ns/op
BenchmarkMutex8Noop-16     	 1620813	       807.3 ns/op
BenchmarkMutex16Noop-16    	  772576	      1827 ns/op
BenchmarkMutex32Noop-16    	  365168	      3630 ns/op
*/

func BenchmarkMutex(b *testing.B) {
	benchmarkMutexN(b, 1, Debug)
}

func BenchmarkMutex4(b *testing.B) {
	benchmarkMutexN(b, 4, Debug)
}

func BenchmarkMutex8(b *testing.B) {
	benchmarkMutexN(b, 8, Debug)
}

func BenchmarkMutex16(b *testing.B) {
	benchmarkMutexN(b, 16, Debug)
}

func BenchmarkMutex32(b *testing.B) {
	benchmarkMutexN(b, 32, Debug)
}

func BenchmarkMutexNoop(b *testing.B) {
	benchmarkMutexN(b, 1, Error)
}

func BenchmarkMutex4Noop(b *testing.B) {
	benchmarkMutexN(b, 4, Error)
}

func BenchmarkMutex8Noop(b *testing.B) {
	benchmarkMutexN(b, 8, Error)
}

func BenchmarkMutex16Noop(b *testing.B) {
	benchmarkMutexN(b, 16, Error)
}

func BenchmarkMutex32Noop(b *testing.B) {
	benchmarkMutexN(b, 32, Error)
}

func benchmarkMutexN(b *testing.B, n int, minS Severity) {
	logger := New(io.Discard)

	logger.Tag("version", "deadbeef")
	logger.Tag("version2", "deadbeef")
	logger.Tag("version3", "deadbeef")
	logger.Tag("version4", "deadbeef")
	logger.Tag("version5", "deadbeef")

	SetSeverity("mutex", minS)

	mu := NewMutex("important_lock")

	var wg sync.WaitGroup
	for j := 0; j < n; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				mu.Lock(logger)
				mu.Unlock(logger)
			}
		}()
	}
	wg.Wait()
}
