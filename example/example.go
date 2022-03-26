package main

import (
	"time"

	"github.com/gen95mis/go-console"
)

func main() {
	c := console.NewView()
	line0 := c.NewLine()
	pb0 := c.NewProgresBar()
	pb1 := c.NewProgresBar()

	for i := 0; i <= 100; i++ {
		line0.Set("number of iterations: %d", i)
		pb0.Set("downloading file.txt", i, 100)
		pb1.Set("downloading count files", i*12, 100*12)
		time.Sleep(time.Millisecond * 50)
	}

	line2 := c.NewLine()
	line2.Set("Completed!")
	time.Sleep(time.Second * 1)
}
