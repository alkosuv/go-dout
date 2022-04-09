package main

import (
	"time"

	"github.com/gen95mis/go-dout"
)

func main() {
	v := dout.GetView()

	line0 := v.NewLine()
	pb0 := v.NewProgressBar(10)
	pb1 := v.NewProgressBar(50)
	pb2 := v.NewProgressBarWithTime(50)

	v.NewTitle("\t\t *** some title *** ")

	for i := 0; i <= 100; i++ {
		line0.Set("number of iterations: %d\n", i)
		pb0.Set("downloading file.txt", i, 100)
		pb1.Set("downloading count files", i*12, 100*12)
		pb2.Set("downloading file.csv with time", i, 100)
		time.Sleep(time.Millisecond * 10)
	}

	v.Print("Completed! ")
	v.Println("Exit!")
	time.Sleep(time.Second * 1)
}
