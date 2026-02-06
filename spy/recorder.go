package spy

import (
	"path"
	"runtime"
	"strings"
	"sync"
)

type Recorder struct {
	mu      sync.Mutex
	History [][]any
}

func (r *Recorder) Record(args ...any) {
	if r == nil {
		r = &Recorder{}
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	event := make([]any, 0, 1+len(args))

	funcName := "unknown"
	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = cleanName(fn.Name())
		}
	}

	event = append(event, funcName)
	event = append(event, args...)

	if r.History == nil {
		r.History = make([][]any, 0, 10)
	}
	r.History = append(r.History, event)
}

func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.History = nil
}

func cleanName(name string) string {
	base := path.Base(name)
	if dot := strings.LastIndex(base, "."); dot >= 0 {
		return base[dot+1:]
	}
	return base
}
