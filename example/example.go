package main

import (
	"time"

	"github.com/gen95mis/go-console"
)

func main() {
	v := console.GetView()

	line0 := v.NewLine()
	pb0 := v.NewProgressBar()
	pb1 := v.NewProgressBar()

	v.NewTitle("\t\t *** some title *** ")

	for i := 0; i <= 100; i++ {
		line0.Set("number of iterations: %d", i)
		pb0.Set("downloading file.txt", i, 100)
		pb1.Set("downloading count files", i*12, 100*12)
		time.Sleep(time.Millisecond * 10)
	}

	v.NewLine().Set("Completed!")
	v.Print("Exit")
	time.Sleep(time.Second * 1)
}
