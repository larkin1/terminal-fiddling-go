package term

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type frameBuf struct {
	width		uint
	height	uint
	buf			[][]rune
}

func NewBuf(width uint, height uint) *frameBuf {
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

func (buf *frameBuf) ChangeRow(rownum uint, newrow []rune) error {
	if rownum >= buf.height {
		return errors.New("row number out of bounds")
	}
	buf.buf[rownum] = newrow
	return nil
}

func (buf *frameBuf) ClearBuf() {
	for y := range buf.height {
		for x := range buf.width {
			buf.buf[y][x] = ' '
		}
	}
}

func (buf *frameBuf) SetCell(x uint, y uint, char rune) error {
	if x >= buf.width {
		return errors.New("x value out of bounds")
	} else if y >= buf.height {
		return errors.New("y value out of bounds")
	}
	buf.buf[y][x] = char
	return nil
}

func ClearingResize(buf **frameBuf, width uint, height uint) {
	*buf = NewBuf(width, height)
}

func resize(buf **frameBuf, width uint, height uint) {
	newBuffer := NewBuf(width, height)

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

func (buf *frameBuf) Show() {
	var out strings.Builder

	for _, line := range buf.buf {
		out.WriteString(string(line))
		out.WriteByte('\n')
	}

	result := out.String()
	fmt.Print("\x1b[H")
	fmt.Print(result[:len(result)-1])
}

func AltMode() {
	fmt.Print("\x1b[?1049h") // alt screen
	fmt.Print("\x1b[?25l")   // hide cursor
	fmt.Print("\x1b[2J\x1b[H")
}

func NormalMode() {
	fmt.Print("\x1b[?25h") // show cursor
  fmt.Print("\x1b[?1049l") // back to normal screen
}

func GetSize() (uint, uint, error) {
	fd := int(os.Stdin.Fd())
	w, h, err := term.GetSize(fd)
	return uint(w), uint(h), err 

}
