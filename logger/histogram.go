package logger

import "time"

type Histo interface {
	Observe(float64)
	ObserveElapsed(time.Duration)
}

type noopHistogram struct{}
func (*noopHistogram) Observe(float64) {}
func (*noopHistogram) ObserveElapsed(time.Duration) {}
var aNoopHistogram = &noopHistogram{}

type Unit int
const (
	Seconds Unit = iota
	Milliseconds
	Nanoseconds
	Count
)

type HistogramDef struct {
	unit Unit
	tags []string
}

type Histogram struct {
	unit Unit
	ks []string
	vs []string

	// dummy impl
	values []float64
}

func (h* Histogram) Observe(v float64) {
	switch h.unit {
	case Nanoseconds, Milliseconds, Seconds:
		panic("Wrong unit!")
	}
	h.values = append(h.values, v)
}

func (h* Histogram) ObserveElapsed(v time.Duration) {
	switch h.unit {
	case Nanoseconds:
		h.values = append(h.values, float64(v))
	case Milliseconds:
		h.values = append(h.values, float64(v/time.Millisecond))
	case Seconds:
		h.values = append(h.values, float64(v/time.Second))
	default:
		panic("Wrong unit!")
	}
}
