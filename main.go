package main

import (
	"errors"
	"fmt"
	"os"
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

func changeRow(buf *frameBuf, rownum uint, newrow []rune) error {
	if rownum > buf.height {
		return errors.New("row number out of bounds")
	}
	buf.buf[rownum] = newrow
	return nil
}

func clearBuf(buf *frameBuf) {
	for y := range buf.height {
		for x := range buf.width {
			buf.buf[y][x] = ' '
		}
	}
}

func setCell(buf *frameBuf, x uint, y uint, char rune) error {
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

func main() {
	fmt.Print("\x1b[?1049h") // alt screen
  fmt.Print("\x1b[?25l")   // hide cursor

  defer func() {
    fmt.Print("\x1b[?25h") // show cursor
    fmt.Print("\x1b[?1049l") // back to normal screen
  }()

	fmt.Print("\x1b[2J\x1b[H")

	fd := int(os.Stdin.Fd())
	w, h, err := term.GetSize(fd)
	if err != nil {panic(err)}
	buf := newBuf(w, h)
	for {
		fd := int(os.Stdin.Fd())
		w, h, err := term.GetSize(fd)
		if err != nil { panic(err) }
	
		fmt.Printf("\x1b[%d;%dH%s\n", h-1, w-1, "s")
		time.Sleep(33 * time.Millisecond)
	}
}
