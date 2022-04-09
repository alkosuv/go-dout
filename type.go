package dout

import (
	"fmt"
	"strings"
	"time"
)

type Line struct {
	value string
}

func newLine() *Line {
	return &Line{}
}

// Set formats according to a format specifier and writes to standard output.
func (l *Line) Set(format string, a ...any) {
	mutex.Lock()
	defer mutex.Unlock()

	l.value = fmt.Sprintf(format, a...)
}

// Get returns a Line value
func (l *Line) Get() string {
	return l.value
}

type ProgressBar struct {
	countLineProcess int
	value            string
	timeStart        *time.Time
}

func newProgressBar(countLineProcess int) *ProgressBar {
	return &ProgressBar{
		countLineProcess: countLineProcess,
	}
}

func newProgressBarWithTime(countLineProcess int) *ProgressBar {
	t := time.Now()
	return &ProgressBar{
		countLineProcess: countLineProcess,
		timeStart:        &t,
	}
}

// Set formats according to a format specifier and writes to standard output.
func (p *ProgressBar) Set(description string, current, max int) {
	mutex.Lock()
	defer mutex.Unlock()

	p.value = p.genProgresBar(description, current, max, p.countLineProcess)
}

// Get returns a ProgressBar value
func (p *ProgressBar) Get() string {
	return p.value
}

func (p *ProgressBar) genProgresBar(description string, carrent, max, countLineProcess int) string {
	persent := float32(carrent) / float32(max) * 100
	count := 0
	if persent != 0 {
		count = int(persent) * countLineProcess / (100)
	}

	str1 := strings.Repeat("#", count)
	str2 := strings.Repeat("*", countLineProcess-count)

	result := ""
	if p.timeStart == nil {
		result = fmt.Sprintf("%s [%s%s] (%d/%d)\n", description, str1, str2, carrent, max)
	} else {
		difference := time.Since(*p.timeStart).Round(time.Second)
		result = fmt.Sprintf("%s [%s%s] (%d/%d) [%s]\n", description, str1, str2, carrent, max, difference.String())
	}

	return result
}
