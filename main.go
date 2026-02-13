package main

import (
	"time"

	lib "learnterm/lib"
)

func main() {
	// set up the terminal
	lib.AltMode()
	defer lib.NormalMode()

	// get screen size
	w, h, err := lib.GetSize()
	if err != nil {
		panic(err)
	}

	// create a buffer
	buf := lib.NewBuf(uint(w), uint(h))

	buf.SetCell(5, 5, 's')
	buf.SetCell(w-1, h-1, 'e')

	for {
		buf.Show()
		time.Sleep(33 * time.Millisecond)
	}
}
