package dout

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	refresh = time.Millisecond * 10
)

var (
	mutex       sync.Mutex
	once        sync.Once
	viewinstans *view
)

type get interface {
	Get() string
}

type any = interface{}

type View interface {
	NewLine() *Line
	NewTitle(format string, a ...any)
	NewProgressBar(countLineProcess int) *ProgressBar
	Print(a ...any)
	Printf(format string, a ...any)
	Println(a ...any)
	ResetView()
	ClearTerminal()
}

type view struct {
	out           io.Writer
	lines         []get
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
	v.lines = make([]get, 0)
	v.lastLineCount = -1
}

func newView() *view {
	v := new(view)
	v.out = io.Writer(os.Stdout)
	v.lines = make([]get, 0)
	v.lastLineCount = -1

	go v.loop()

	return v
}

// NewLine create new Line
func (v *view) NewLine() *Line {
	mutex.Lock()
	defer mutex.Unlock()

	l := newLine()
	v.lines = append(v.lines, l)
	return l
}

// NewTitle create new Title
func (v *view) NewTitle(format string, a ...any) {
	mutex.Lock()

	l := newLine()
	v.lines = append([]get{l}, v.lines...)

	mutex.Unlock()

	format += "\n"
	l.Set(format, a...)
}

// NewProgressBar create new ProgressBar. Negative countLineProcess values ​​will be treated as positive.
func (v *view) NewProgressBar(countLineProcess int) *ProgressBar {
	mutex.Lock()
	defer mutex.Unlock()

	countLineProcess = int(math.Abs(float64(countLineProcess)))

	pb := newProgressBar(countLineProcess)
	v.lines = append(v.lines, pb)

	return pb
}

// Print formats using the default formats for its operands and writes to standard output.
func (v *view) Print(a ...any) {
	v.NewLine().Set(fmt.Sprint(a...))
}

// Print formats according to a format specifier and writes to standard output.
func (v *view) Printf(format string, a ...any) {
	v.NewLine().Set(format, a...)
}

// Println formats using the default formats for its operands and writes to standard outpu and newline is appended.
func (v *view) Println(a ...any) {
	v.NewLine().Set(fmt.Sprintln(a...))
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
	for _, n := range v.lines {
		str := n.Get()
		count += strings.Count(str, "\n")
		buffer.WriteString(n.Get())
	}
	v.out.Write(buffer.Bytes())

	if count == 0 {
		v.lastLineCount = -1
	} else {
		v.lastLineCount = count
	}
}

func (v *view) clearLines(count int) {
	clear := fmt.Sprintf("\033[%dA\033[2K", count)
	fmt.Fprint(v.out, clear)
}
