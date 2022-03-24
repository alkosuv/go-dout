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

type Console interface {
	NewLine() Line
	ClearTerminal()
}

type console struct {
	mutex *sync.Mutex

	out   io.Writer
	lines []*line
}

func NewConsole() Console {
	c := new(console)
	c.out = io.Writer(os.Stdout)
	c.lines = make([]*line, 0)
	c.mutex = new(sync.Mutex)

	go c.loop()

	return c
}

func (c *console) NewLine() Line {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	l := newLine(c.mutex)
	c.lines = append(c.lines, l)
	return l
}

func (c *console) ClearTerminal() {
	fmt.Fprint(c.out, "\033[H\033[2J")
}

func (c *console) loop() {
	ticker := time.NewTicker(refresh)

	for range ticker.C {
		c.output()
	}
}

func (c *console) output() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.lines) == 0 {
		return
	}

	c.clearLines()
	var buffer bytes.Buffer

	for _, l := range c.lines {
		buffer.WriteString(l.str + "\n")
	}
	c.out.Write(buffer.Bytes())
}

func (c *console) clearLines() {
	clear := fmt.Sprintf("\033[%dA\033[2K", len(c.lines))
	fmt.Fprint(c.out, clear)
}
