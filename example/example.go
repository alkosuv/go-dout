package main

import (
	"time"

	"github.com/gen95mis/go-console"
)

func main() {
	v := console.NewView()

	line0 := v.NewLine()
	pb0 := v.NewProgresBar()
	pb1 := v.NewProgresBar()

	v.NewTitle("\t\t *** some title *** ")

	for i := 0; i <= 100; i++ {
		line0.Set("number of iterations: %d", i)
		pb0.Set("downloading file.txt", i, 100)
		pb1.Set("downloading count files", i*12, 100*12)
		time.Sleep(time.Millisecond * 50)
	}

	v.NewLine().Set("Completed!")
	time.Sleep(time.Second * 1)
}
