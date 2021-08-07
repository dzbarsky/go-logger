package logger

import (
	"sync"
)

type Mutex struct {
	mu *sync.Mutex

	name string
}

const (
	nameTag = "mutex_name"
)

var (
	mutexScope = &Scope{name: "mutex"}
	mutexHisto = &HistogramDef{
		unit: Nanoseconds,
		// Add version tag so we demo some tags from logger
		tags: []string{"version",  nameTag},
	}
)

func NewMutex(name string) Mutex {
	return Mutex{name: name, mu: &sync.Mutex{}}
}

func (m Mutex) Lock(l *Logger) {
	debug := l.ScopedStream(mutexScope, Debug)
	debug.Tag(nameTag, m.name)

	start := debug.Now()
	debug.Log("Acquiring mutex")
	m.mu.Lock()

	// TODO: it's a bit weird to create a tag here since this will affect future calls. probably better to log as an event once we figure out the API
	debug.Tag("mutex_time", debug.StringSince(start))
	debug.Log("Acquired mutex")
	debug.Histogram(mutexHisto).ObserveElapsed(debug.Since(start))
}

func (m Mutex) Unlock(l *Logger) {
	m.mu.Unlock()
	debug := l.ScopedStream(mutexScope, Debug)
	debug.Tag(nameTag, m.name)
	debug.Log("Released mutex")
}
