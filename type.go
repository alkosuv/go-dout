package dout

import (
	"fmt"
	"strings"
	"sync"
)

type Line struct {
	mutex *sync.Mutex
	value string
}

func newLine(mutex *sync.Mutex) *Line {
	return &Line{
		mutex: mutex,
	}
}

func (l *Line) Set(format string, a ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.value = fmt.Sprintf(format, a...)
}

func (l *Line) Get() string {
	return l.value
}

type ProgressBar struct {
	mutex *sync.Mutex
	value string
}

func newProgressBar(mutex *sync.Mutex) *ProgressBar {
	return &ProgressBar{
		mutex: mutex,
	}
}

func (p *ProgressBar) Set(description string, current, max int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.value = genProgresBar(description, current, max)
}

func (p *ProgressBar) Get() string {
	return p.value
}

func genProgresBar(description string, carrent, max int) string {
	const countLineProcess int = 25

	persent := float32(carrent) / float32(max) * 100
	count := 0
	if persent != 0 {
		count = int(persent) * countLineProcess / (100)
	}

	str1 := strings.Repeat("#", count)
	str2 := strings.Repeat("*", countLineProcess-count)

	result := fmt.Sprintf("%s [%s%s] (%d/%d)", description, str1, str2, carrent, max)

	return result
}
