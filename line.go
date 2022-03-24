package console

import "fmt"

type Line interface {
	Print(format string, a ...interface{})
}

type line struct {
	str string
}

func (l *line) Print(format string, a ...interface{}) {
	l.str = fmt.Sprintf(format, a...)
}
