package logger

import (
	"io"
)

const batch = 100

type stdoutStream struct {
	w io.Writer
	stream chan []byte
}

func newStdout(w io.Writer) *stdoutStream {
	s := &stdoutStream{
		w: w,
		stream: make(chan []byte, 100),
	}
	go s.emit()
	return s
}

func (s *stdoutStream) Write(val []byte) {
	s.stream <- val
}

func (s *stdoutStream) emit() {
	for {
		val := <- s.stream
		_, err := s.w.Write(val)
		if err != nil {
			panic(err)
		}
	}
}
