package console

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

var once sync.Once
var viewinstans *view

type get interface {
	Get() string
}

type View interface {
	NewLine() *Line
	NewTitle(format string, a ...interface{})
	NewProgresBar() *ProgresBar
	Print(format string, a ...interface{})
	ResetView()
	ClearTerminal()
}

type view struct {
	out   io.Writer
	mutex *sync.Mutex

	node          []get
	lastLineCount int
}

func GetView() View {
	once.Do(func() {
		viewinstans = newView()
	})

	return viewinstans
}

func (v *view) ResetView() {
	v.node = make([]get, 0)
	v.lastLineCount = -1
}

func newView() *view {
	v := new(view)
	v.out = io.Writer(os.Stdout)
	v.mutex = new(sync.Mutex)
	v.node = make([]get, 0)
	v.lastLineCount = -1

	go v.loop()

	return v
}

func (v *view) NewLine() *Line {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	l := newLine(v.mutex)
	v.node = append(v.node, l)
	return l
}

func (v *view) NewTitle(format string, a ...interface{}) {
	v.mutex.Lock()

	l := newLine(v.mutex)
	v.node = append([]get{l}, v.node...)

	v.mutex.Unlock()

	l.Set(format, a...)
}

func (v *view) NewProgresBar() *ProgresBar {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	pb := newProgresBar(v.mutex)
	v.node = append(v.node, pb)

	return pb
}

func (v *view) Print(format string, a ...interface{}) {
	v.NewLine().Set(format, a...)
}

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
	v.mutex.Lock()
	defer v.mutex.Unlock()

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
