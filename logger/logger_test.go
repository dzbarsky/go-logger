package logger

import (
	"io"
	"log"
	"sync"
	"testing"
)

/*
goos: darwin
goarch: amd64
pkg: logger/logger
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkLogger-16          	 2015685	       590.8 ns/op
BenchmarkLogger4-16         	  687106	      1709 ns/op
BenchmarkLogger8-16         	  305707	      3685 ns/op
BenchmarkLogger16-16        	  140256	      7646 ns/op
BenchmarkLogger32-16        	   74095	     16142 ns/op
BenchmarkLoggerNoop-16      	265901402	         4.472 ns/op
BenchmarkLogger4Noop-16     	268158448	         4.499 ns/op
BenchmarkLogger8Noop-16     	230825785	         5.050 ns/op
BenchmarkLogger16Noop-16    	137533003	         8.709 ns/op
BenchmarkLogger32Noop-16    	57586251	        17.81 ns/op
*/

func BenchmarkLogger(b *testing.B) {
	benchmarkLoggerN(b, 1, Info)
}

func BenchmarkLogger4(b *testing.B) {
	benchmarkLoggerN(b, 4, Info)
}

func BenchmarkLogger8(b *testing.B) {
	benchmarkLoggerN(b, 8, Info)
}

func BenchmarkLogger16(b *testing.B) {
	benchmarkLoggerN(b, 16, Info)
}

func BenchmarkLogger32(b *testing.B) {
	benchmarkLoggerN(b, 32, Info)
}

func BenchmarkLoggerNoop(b *testing.B) {
	benchmarkLoggerN(b, 1, Error)
}

func BenchmarkLogger4Noop(b *testing.B) {
	benchmarkLoggerN(b, 4, Error)
}

func BenchmarkLogger8Noop(b *testing.B) {
	benchmarkLoggerN(b, 8, Error)
}

func BenchmarkLogger16Noop(b *testing.B) {
	benchmarkLoggerN(b, 16, Error)
}

func BenchmarkLogger32Noop(b *testing.B) {
	benchmarkLoggerN(b, 32, Error)
}

func benchmarkLoggerN(b *testing.B, n int, minS Severity) {
	log.SetOutput(io.Discard)

	logger := New()

	logger.Tag("version", "deadbeef")
	logger = logger.PushScope("s2")
	logger.Tag("version2", "deadbeef")
	logger = logger.PushScope("s2")
	logger.Tag("version3", "deadbeef")
	logger = logger.PushScope("s3")
	logger.Tag("version4", "deadbeef")
	logger = logger.PushScope("s5")
	logger.Tag("version5", "deadbeef")

	SetSeverity("s5", minS)

	var wg sync.WaitGroup
	for j := 0; j < n; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				logger.Info("Hello, I am alive and well")
			}
		}()
	}
	wg.Wait()
}
