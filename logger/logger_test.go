package logger

import (
	"log"
	"io"
	"testing"
	"sync"
)

/*
goos: darwin
goarch: amd64
pkg: logger/logger
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkLogger-16          	 2313824	       507.3 ns/op
BenchmarkLogger4-16         	  608257	      1797 ns/op
BenchmarkLogger8-16         	  301158	      3880 ns/op
BenchmarkLogger16-16        	  148124	      7947 ns/op
BenchmarkLogger32-16        	   75054	     15770 ns/op
BenchmarkLoggerNoop-16      	425889259	         2.724 ns/op
BenchmarkLogger4Noop-16     	417245281	         2.860 ns/op
BenchmarkLogger8Noop-16     	400372350	         2.897 ns/op
BenchmarkLogger16Noop-16    	222931371	         5.310 ns/op
BenchmarkLogger32Noop-16    	110269083	        10.70 ns/op
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
	benchmarkLoggerN(b, 1, Debug)
}

func BenchmarkLogger4Noop(b *testing.B) {
	benchmarkLoggerN(b, 4, Debug)
}

func BenchmarkLogger8Noop(b *testing.B) {
	benchmarkLoggerN(b, 8, Debug)
}

func BenchmarkLogger16Noop(b *testing.B) {
	benchmarkLoggerN(b, 16, Debug)
}

func BenchmarkLogger32Noop(b *testing.B) {
	benchmarkLoggerN(b, 32, Debug)
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
