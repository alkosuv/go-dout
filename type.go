package console

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

type ProgresBar struct {
	mutex *sync.Mutex
	value string
}

func newProgresBar(mutex *sync.Mutex) *ProgresBar {
	return &ProgresBar{
		mutex: mutex,
	}
}

func (p *ProgresBar) Set(desctioption string, current, max int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.value = genProgresBar(desctioption, current, max)
}

func (p *ProgresBar) Get() string {
	return p.value
}

func genProgresBar(desctioption string, carrent, max int) string {
	const countLineProces int = 25

	persent := float32(carrent) / float32(max) * 100
	count := 0
	if persent != 0 {
		count = int(persent) * countLineProces / (100)
	}

	str1 := strings.Repeat("#", count)
	str2 := strings.Repeat("*", countLineProces-count)

	result := fmt.Sprintf("%s [%s%s] (%d/%d)", desctioption, str1, str2, carrent, max)

	return result
}
