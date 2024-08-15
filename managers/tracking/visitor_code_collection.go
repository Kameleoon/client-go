package tracking

import (
	cmap "github.com/orcaman/concurrent-map/v2"
)

type VisitorCodeCollection interface {
	Range(f func(visitorCode string) bool)
}

type SliceVisitorCodeCollection struct {
	visitorCodes []string
}

func (vcc SliceVisitorCodeCollection) Range(f func(visitorCode string) bool) {
	for _, vc := range vcc.visitorCodes {
		if !f(vc) {
			break
		}
	}
}

type CMapVisitorCodeCollection struct {
	visitorCodes *cmap.ConcurrentMap[string, struct{}]
}

func (vcc CMapVisitorCodeCollection) Range(f func(visitorCode string) bool) {
	for kv := range vcc.visitorCodes.IterBuffered() {
		if !f(kv.Key) {
			break
		}
	}
}

type SingletonVisitorCodeCollection struct {
	visitorCode string
}

func (vcc SingletonVisitorCodeCollection) Range(f func(visitorCode string) bool) {
	f(vcc.visitorCode)
}
