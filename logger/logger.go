package logger

import (
	"log"
	"sync"
	"sync/atomic"
	"strings"
)

type Severity int

const (
	Info Severity = iota
	Debug
	Error
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
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

type Logger struct {
	name string
	parent *Logger

	cachedGeneration int64
	cachedMinSeverity Severity

	keys []string
	values []interface{}
}

func New() *Logger{
	return &Logger{}
}

func (l *Logger) Tag(k string, v interface{}) {
	l.keys = append(l.keys, k)
	l.values = append(l.values, v)
}

func (l *Logger) PushScope(name string) *Logger {
	return &Logger{name: name, parent: l}
}

func (l *Logger) PopScope() *Logger {
	return l.parent
}

func (l *Logger) Info(msg string) {
	l.log(Info, msg)
}

func (l *Logger) Debug(msg string) {
	l.log(Debug, msg)
}

func (l *Logger) Error(msg string) {
	l.log(Error, msg)
}

func (l *Logger) log(severity Severity, msg string) {
	currGeneration := atomic.LoadInt64(&generation)
	if currGeneration != l.cachedGeneration {
		severityMu.Lock()
		// TODO: missing min-level should default to something reasonable?
		l.cachedMinSeverity = severities[l.name]
		severityMu.Unlock()
		l.cachedGeneration = currGeneration
	}
	
	if severity < l.cachedMinSeverity {
		return
	}

	var sb strings.Builder
	sb.Grow(200)
	sb.WriteString("{")
	l.writeTags(&sb)
	sb.WriteString(`"message": "`)
	sb.WriteString(msg)
	sb.WriteString(`, "severity": "`)
	sb.WriteString(severity.String())
	sb.WriteString(`"}`)
	log.Print(sb.String())
}

func (l *Logger) writeTags(sb *strings.Builder) {
	// TODO: escape strings for json?
	for i, k := range l.keys {
		sb.WriteString(`"`)
		sb.WriteString(k)
		sb.WriteString(`": "`)
		sb.WriteString(l.values[i].(string))
		sb.WriteString(`",`)
	}

	if l.parent != nil {
		l.parent.writeTags(sb)
	}
}
