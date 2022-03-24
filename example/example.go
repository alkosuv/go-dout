package main

import (
	"time"

	"github.com/gen95mis/go-console"
)

func main() {
	c := console.NewConsole()
	line0 := c.NewLine()
	line1 := c.NewLine()

	for i := 0; i <= 100; i++ {
		line0.Set("downloading File 0 - %d %%", i)
		line1.Set("downloading File 1 - %d %%", i)
		time.Sleep(time.Millisecond * 100)
	}

	line2 := c.NewLine()
	line2.Set("Completed!")
	time.Sleep(time.Second * 1)
}
