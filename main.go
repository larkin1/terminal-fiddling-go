package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

type frameBuf struct {
	width		uint
	height	uint
	buf			[][]rune
}

func newBuf(width uint, height uint) *frameBuf {
	arr := make([][]rune, height)
	for y := range arr {
		row := make([]rune, width)
		for x := range row {
			row[x] = ' '
		}
		arr[y] = row
	}
	p := frameBuf{width, height, arr}
	return &p
}

func (buf *frameBuf) changeRow(rownum uint, newrow []rune) error {
	if rownum > buf.height {
		return errors.New("row number out of bounds")
	}
	buf.buf[rownum] = newrow
	return nil
}

func (buf *frameBuf) clearBuf() {
	for y := range buf.height {
		for x := range buf.width {
			buf.buf[y][x] = ' '
		}
	}
}

func (buf *frameBuf) setCell(x uint, y uint, char rune) error {
	if x > buf.width {
		return errors.New("x value out of bounds")
	} else if y > buf.height {
		return errors.New("y value out of bounds")
	}
	buf.buf[y][x] = char
	return nil
}

func clearingResize(buf **frameBuf, width uint, height uint) {
	*buf = newBuf(width, height)
}

func resize(buf **frameBuf, width uint, height uint) {
	newBuffer := newBuf(width, height)

	old := *buf
	
	maxY := min(old.width, newBuffer.width)
	maxX := min(old.height, newBuffer.height)
	for y := range maxY {
		for x := range maxX {
			newBuffer.buf[y][x] = old.buf[y][x]
		}
	}
	*buf = newBuffer
}

func (buf *frameBuf) showBuf() {
	var out strings.Builder

	for _, line := range buf.buf {
		out.WriteString(string(line))
		out.WriteByte('\n')
	}

	fmt.Print("\x1b[2J\x1b[H")
	fmt.Println(out.String())
}

func altMode() {
	fmt.Print("\x1b[?1049h") // alt screen
  fmt.Print("\x1b[?25l")   // hide cursor
	fmt.Print("\x1b[2J\x1b[H")
}

func normalMode() {
  fmt.Print("\x1b[?25h") // show cursor
  fmt.Print("\x1b[?1049l") // back to normal screen
}

func getSize() (uint, uint, error) {
	fd := int(os.Stdin.Fd())
	w, h, err := term.GetSize(fd)
	return uint(w), uint(h), err 

}

func main() {
	// set up the terminal
	altMode()
	defer normalMode()

	// get screen size
	w, h, err := getSize()
	if err != nil {
		panic(err)
	}

	// create a buffer
	buf := newBuf(uint(w), uint(h))

	buf.setCell(5, 5, 's')

	for {
		fmt.Printf("\x1b[%d;%dH%s\n", h-1, w-1, "s")
		time.Sleep(33 * time.Millisecond)
	}
}
