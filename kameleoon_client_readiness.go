package kameleoon

import "sync"

type kameleoonClientReadiness struct {
	isInitializing bool
	err            error
	cond           sync.RWMutex
}

func newKameleoonClientReadiness() *kameleoonClientReadiness {
	r := new(kameleoonClientReadiness)
	r.reset()
	return r
}

func (r *kameleoonClientReadiness) reset() {
	r.err = nil
	if !r.isInitializing {
		r.isInitializing = true
		r.cond.Lock()
	}
}
func (r *kameleoonClientReadiness) set(err error) {
	r.err = err
	if r.isInitializing {
		r.cond.Unlock()
		r.isInitializing = false
	}
}

func (r *kameleoonClientReadiness) IsInitializing() bool {
	return r.isInitializing
}

func (r *kameleoonClientReadiness) Error() error {
	return r.err
}

func (r *kameleoonClientReadiness) Wait() error {
	if r.isInitializing {
		r.cond.RLock()
		r.cond.RUnlock()
	}
	return r.err
}
