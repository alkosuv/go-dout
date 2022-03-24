package console

import (
	"fmt"
	"sync"
)

type Line interface {
	Print(format string, a ...interface{})
}

type line struct {
	mutex *sync.Mutex
	str   string
}

func newLine(mutex *sync.Mutex) *line {
	return &line{
		mutex: mutex,
	}
}

func (l *line) Print(format string, a ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.str = fmt.Sprintf(format, a...)
}
