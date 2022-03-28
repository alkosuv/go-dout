package dout

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	refresh = time.Millisecond * 10
)

var mutex sync.Mutex
var once sync.Once
var viewinstans *view

type get interface {
	Get() string
}

type View interface {
	NewLine() *Line
	NewTitle(format string, a ...interface{})
	NewProgressBar() *ProgressBar
	Print(format string, a ...interface{})
	ResetView()
	ClearTerminal()
}

type view struct {
	out           io.Writer
	node          []get
	lastLineCount int
}

// GetView returns View
func GetView() View {
	once.Do(func() {
		viewinstans = newView()
	})

	return viewinstans
}

// ResetView reset view
func (v *view) ResetView() {
	v.node = make([]get, 0)
	v.lastLineCount = -1
}

func newView() *view {
	v := new(view)
	v.out = io.Writer(os.Stdout)
	v.node = make([]get, 0)
	v.lastLineCount = -1

	go v.loop()

	return v
}

// NewLine create new Line
func (v *view) NewLine() *Line {
	mutex.Lock()
	defer mutex.Unlock()

	l := newLine()
	v.node = append(v.node, l)
	return l
}

// NewTitle create new Title
func (v *view) NewTitle(format string, a ...interface{}) {
	mutex.Lock()

	l := newLine()
	v.node = append([]get{l}, v.node...)

	mutex.Unlock()

	l.Set(format, a...)
}

// NewProgressBar create new ProgressBar
func (v *view) NewProgressBar() *ProgressBar {
	mutex.Lock()
	defer mutex.Unlock()

	pb := newProgressBar()
	v.node = append(v.node, pb)

	return pb
}

// Print formats according to a format specifier and writes to standard output.
func (v *view) Print(format string, a ...interface{}) {
	v.NewLine().Set(format, a...)
}

// ClearTerminal clears all information from the terminal
func (v *view) ClearTerminal() {
	fmt.Fprint(v.out, "\033[H\033[2J")
}

func (v *view) loop() {
	ticker := time.NewTicker(refresh)

	for range ticker.C {
		v.output()
	}
}

func (v *view) output() {
	mutex.Lock()
	defer mutex.Unlock()

	count := len(v.node)
	if count == 0 {
		return
	}

	v.clearLines(v.lastLineCount)
	var buffer bytes.Buffer

	for _, n := range v.node {
		buffer.WriteString(n.Get() + "\n")
	}
	v.out.Write(buffer.Bytes())

	v.lastLineCount = count
}

func (v *view) clearLines(count int) {
	clear := fmt.Sprintf("\033[%dA\033[2K", count)
	fmt.Fprint(v.out, clear)
}
