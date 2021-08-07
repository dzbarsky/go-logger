package logger

import (
	"io"
	"bytes"
	"sync"
	"sync/atomic"
	"time"
)

type Severity int

const (
	Debug Severity = iota
	Info
	Error
)

func (s Severity) String() string {
	switch s {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Error:
		return "ERROR"
	}
	return ""
}

var (
	severityMu sync.Mutex
	severities = map[string]Severity{}

	generation = int64(1)
)

func SetSeverity(name string, s Severity) {
	severityMu.Lock()
	defer severityMu.Unlock()

	atomic.AddInt64(&generation, 1)
	severities[name] = s
}

type Scope struct {
	name string

	cachedGeneration  int64
	cachedMinSeverity Severity
}

type Logger struct {
	out *stdoutStream
	scope  *Scope
	parent *Logger

	keys   []string
	values []string
}

func New(w io.Writer) *Logger {
	return &Logger{
		out: newStdout(w),
		scope: &Scope{},
	}
}

func (l *Logger) Tag(k string, v string) {
	l.keys = append(l.keys, k)
	l.values = append(l.values, v)
}

func (l *Logger) PushScope(name string) *Logger {
	return &Logger{
		out: l.out,
		scope:  &Scope{name: name},
		parent: l,
	}
}

func (l *Logger) PopScope() *Logger {
	return l.parent
}

func (l *Logger) Debug(msg string) {
	l.ScopedStream(l.scope, Debug).Log(msg)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.ScopedStream(l.scope, Debug).Logf(msg, args...)
}

func (l *Logger) Info(msg string) {
	l.ScopedStream(l.scope, Info).Log(msg)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.ScopedStream(l.scope, Info).Logf(msg, args...)
}

func (l *Logger) Error(msg string) {
	l.ScopedStream(l.scope, Error).Log(msg)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.ScopedStream(l.scope, Error).Logf(msg, args...)
}

func (l *Logger) valueFor(k string) string {
	for j, key := range l.keys {
		if key == k {
			return l.values[j]
		}
	}
	if l.parent != nil {
		l.parent.valueFor(k)
	}
	panic("no such key: " + k)
}

type Stream interface {
	Log(msg string)
	Logf(msg string, args ...interface{})
	LogApply(func() string)
	Tag(key string, value string)

	Now() time.Time
	Since(time.Time) time.Duration
	StringSince(time.Time) string

	Histogram(h *HistogramDef) Histo
}

func (l *Logger) ScopedStream(s *Scope, severity Severity) Stream {
	currGeneration := atomic.LoadInt64(&generation)
	if currGeneration != s.cachedGeneration {
		severityMu.Lock()
		// TODO: missing min-level should default to something reasonable?
		s.cachedMinSeverity = severities[s.name]
		severityMu.Unlock()
		s.cachedGeneration = currGeneration
	}

	if severity < s.cachedMinSeverity {
		return aNoopStream
	}

	// TODO: preallocate some keys and values in the slice and pool the streams?
	ls := &logStream{
		severity: severity,
		l: l,
		buf: bytes.NewBuffer(make([]byte, 0, 100)),
	}
	writeKV(ls.buf, "severity", severity.String())
	return ls
}

type noopStream struct{}

func (*noopStream) Log(msg string)                       {}
func (*noopStream) Logf(msg string, args ...interface{}) {}
func (*noopStream) LogApply(func() string)               {}
func (*noopStream) Tag(key string, value string)         {}
func (*noopStream) Now() time.Time                       { return time.Time{} }
func (*noopStream) Since(time.Time) time.Duration               { return 0 }
func (*noopStream) StringSince(time.Time) string               { return "" }
func (*noopStream) Histogram(*HistogramDef) Histo { return aNoopHistogram }

var aNoopStream = &noopStream{}

type logStream struct {
	severity Severity
	l        *Logger

	keys   []string
	values []string

	buf *bytes.Buffer
}

func (l *logStream) Tag(k string, v string) {
	l.keys = append(l.keys, k)
	l.values = append(l.values, v)
	writeKV(l.buf, k, v)
}

func (l *logStream) Logf(msg string, args ...interface{}) {
	l.Log(msg)
}

func (l *logStream) Log(msg string) {
	buf := bytes.NewBuffer(make([]byte, 0, 200))

	buf.WriteRune('{')
	buf.Write(l.buf.Bytes())
	writeKV(buf, "message", msg)
	buf.WriteRune('}')

	l.l.out.Write(buf.Bytes())
}

func (l *logStream) LogApply(msg func() string) {
	l.Log(msg())
}

func (l *logStream) Now() time.Time { return time.Now() }
func (l *logStream) Since(s time.Time) time.Duration {
	return time.Since(s)
}
func (l *logStream) StringSince(s time.Time) string {
	return time.Since(s).String()
}

func (l *logStream) valueFor(k string) string {
	for j, key := range l.keys {
		if key == k {
			return l.values[j]
		}
	}
	return l.l.valueFor(k)
}

func (l *logStream) Histogram(h *HistogramDef) Histo {
	ks := make([]string, len(h.tags))
	vs := make([]string, len(h.tags))
	for i, k := range h.tags {
		ks[i] = k
		vs[i] = l.valueFor(k)
	}

	return &Histogram{unit: h.unit, ks: ks, vs: vs}
}

func writeKVs(buf *bytes.Buffer, keys []string, values []string) {
	// TODO: escape strings for json?
	for i, k := range keys {
		writeKV(buf, k, values[i])
	}
}


func writeKV(buf *bytes.Buffer, k, v string) {
	buf.WriteRune('"')
	buf.WriteString(k)
	buf.WriteString(`": "`)
	buf.WriteString(v)
	buf.WriteString(`",`)
}

func (l *Logger) writeTags(buf *bytes.Buffer) {
	writeKVs(buf, l.keys, l.values)

	if l.parent != nil {
		l.parent.writeTags(buf)
	}
}

