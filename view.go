package dout

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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
	NewProgressBar(countLineProcess int) *ProgressBar
	Print(str string)
	Printf(format string, a ...interface{})
	Println(str string)
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

	format += "\n"
	l.Set(format, a...)
}

// NewProgressBar create new ProgressBar. If the countLineProcess is 0, then the default countLineProcess is set to 25.
func (v *view) NewProgressBar(countLineProcess int) *ProgressBar {
	mutex.Lock()
	defer mutex.Unlock()

	if countLineProcess == 0 {
		countLineProcess = 25
	}

	pb := newProgressBar(countLineProcess)
	v.node = append(v.node, pb)

	return pb
}

// Print formats using the default formats for its operands and writes to standard output.
func (v *view) Print(str string) {
	v.NewLine().Set(str)
}

// Print formats according to a format specifier and writes to standard output.
func (v *view) Printf(format string, a ...interface{}) {
	v.NewLine().Set(format, a...)
}

// Println formats using the default formats for its operands and writes to standard outpu and newline is appended.
func (v *view) Println(str string) {
	v.NewLine().Set(str + "\n")
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

	v.clearLines(v.lastLineCount)
	var buffer bytes.Buffer

	count := 0
	for _, n := range v.node {
		str := n.Get()
		count += strings.Count(str, "\n")
		buffer.WriteString(n.Get())
	}
	v.out.Write(buffer.Bytes())

	v.lastLineCount = count
}

func (v *view) clearLines(count int) {
	clear := fmt.Sprintf("\033[%dA\033[2K", count)
	fmt.Fprint(v.out, clear)
}
