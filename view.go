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

var mutex = new(sync.Mutex)

type get interface {
	Get() string
}

type View interface {
	NewLine() *Line
	NewProgresBar() *ProgresBar
	ClearTerminal()
}

type view struct {
	out           io.Writer
	node          []get
	lastLineCount int
}

func NewView() View {
	v := new(view)
	v.out = io.Writer(os.Stdout)
	v.node = make([]get, 0)

	go v.loop()

	return v
}

func (v *view) NewLine() *Line {
	mutex.Lock()
	defer mutex.Unlock()

	l := newLine()
	v.node = append(v.node, l)
	return l
}

func (v *view) NewProgresBar() *ProgresBar {
	mutex.Lock()
	defer mutex.Unlock()

	pb := newProgresBar()
	v.node = append(v.node, pb)

	return pb
}

func (c *view) ClearTerminal() {
	fmt.Fprint(c.out, "\033[H\033[2J")
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
