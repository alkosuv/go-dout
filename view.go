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

type get interface {
	Get() string
}

type View interface {
	NewLine() *Line
	NewTitle(format string, a ...interface{})
	NewProgresBar() *ProgresBar
	ClearTerminal()
	Complete()
}

type view struct {
	out   io.Writer
	mutex *sync.Mutex

	node          []get
	lastLineCount int

	complete chan struct{}
}

func NewView() View {
	v := new(view)
	v.out = io.Writer(os.Stdout)
	v.mutex = new(sync.Mutex)
	v.node = make([]get, 0)
	v.complete = make(chan struct{})

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

func (v *view) ClearTerminal() {
	fmt.Fprint(v.out, "\033[H\033[2J")
}

func (v *view) Complete() {
	v.complete <- struct{}{}
}

func (v *view) loop() {
	ticker := time.NewTicker(refresh)

	for {
		select {
		case <-ticker.C:
			v.output()
		case <-v.complete:
			goto brk
		}
	}

brk:
	return
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
