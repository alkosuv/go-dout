package console

import (
	"fmt"
	"sync"
)

type Line interface {
	Set(format string, a ...interface{})
	Get() string
}

type line struct {
	mutex *sync.Mutex
	value string
}

func newLine(mutex *sync.Mutex) Line {
	return &line{
		mutex: mutex,
	}
}

func (l *line) Set(format string, a ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.value = fmt.Sprintf(format, a...)
}

func (l *line) Get() string {
	return l.value
}
