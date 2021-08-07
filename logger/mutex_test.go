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
BenchmarkMutex-16          	  431023	      2533 ns/op
BenchmarkMutex4-16         	  127555	      9422 ns/op
BenchmarkMutex8-16         	   63049	     19030 ns/op
BenchmarkMutex16-16        	   31110	     37865 ns/op
BenchmarkMutex32-16        	   15757	     76381 ns/op
BenchmarkMutexNoop-16      	33648066	        35.46 ns/op
BenchmarkMutex4Noop-16     	 4667289	       256.7 ns/op
BenchmarkMutex8Noop-16     	 1509859	       804.7 ns/op
BenchmarkMutex16Noop-16    	  743716	      1697 ns/op
BenchmarkMutex32Noop-16    	  352146	      3545 ns/op
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
