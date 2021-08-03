package logger

import (
	"sync"
)

type Mutex struct {
	mu *sync.Mutex

	name string
}

var mutexScope = &Scope{name: "mutex"}

func NewMutex(name string) Mutex {
	return Mutex{name: name, mu: &sync.Mutex{}}
}

func (m Mutex) Lock(l *Logger) {
	debug := l.ScopedStream(mutexScope, Debug)
	debug.Tag("mutex_name", m.name)

	start := debug.Now()
	debug.Log("Acquiring mutex")
	m.mu.Lock()

	// TODO: it's a bit weird to create a tag here since this will affect future calls. probably better to log as an event once we figure out the API
	debug.Tag("mutex_time", debug.Since(start))
	debug.Log("Acquired mutex")
}

func (m Mutex) Unlock(l *Logger) {
	m.mu.Unlock()
	debug := l.ScopedStream(mutexScope, Debug)
	debug.Tag("mutex_name", m.name)
	debug.Log("Released mutex")
}
