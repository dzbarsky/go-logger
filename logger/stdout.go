package logger

import (
	"bytes"
	"io"
	"log"

	"code.cloudfoundry.org/go-diodes"
)

const batch = 1 << 24

type stdoutStream struct {
	w io.Writer
	stream *diodes.ManyToOne
}

func newStdout(w io.Writer) *stdoutStream {
	d := diodes.NewManyToOne(batch, diodes.AlertFunc(func(missed int) {
		log.Printf("Dropped %d messages", missed)
	}))

	s := &stdoutStream{
		w: w,
		stream: d,
	}
	go s.emit()
	return s
}

func (s *stdoutStream) Write(buf *bytes.Buffer) {
	s.stream.Set(diodes.GenericDataType(buf))
}

func (s *stdoutStream) emit() {
	poller := diodes.NewPoller(s.stream)
	for {
		buf := (*bytes.Buffer)(poller.Next())
		_, err := s.w.Write(buf.Bytes())
		if err != nil {
			panic(err)
		}
	}
}
