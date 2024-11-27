package window

import (
	"errors"
)

// slice of elements
type elements []element

func (e elements) draw(canvas []byte, ws winsize) error {
	// start from 0
	for _, element := range e {
		if err := element.draw(canvas, ws); err != nil {
			return err
		}
	}

	return nil
}

// element represents window element.
type element struct {
	// color
	color color
	// width is the columns that the element will take. 0 means it'll take all
	// all the columns to the right
	width int
	// height is the rows that the element will take. 0 means it'll take all
	// all the rows to the bottom
	height int
	// element value
	value []byte
	// x point of a draw
	x int
	// y point of a draw
	y int
}

func (e element) draw(canvas []byte, ws winsize) error {
	lastXIdx := e.x + e.width  // last x index of the element
	lastYIdx := e.y + e.height // last y index of the element

	// Check if the canvas is set
	if len(canvas) == 0 || ws.rowSize == 0 {
		return errors.New("canvas not set")
	}

	rowsNum := len(canvas) / ws.rowSize

	// Check if the element fits in the window
	if rowsNum < lastYIdx {
		return errors.New("element does not fit in the window")
	}

	v := []byte(e.value)
	for y := e.y; y < lastYIdx; y++ {
		// Check if the element fits in the `X` line
		if lastXIdx*ws.colSize-1 > ws.rowSize {
			return errors.New("element does not fit in the X line")
		}

		var colidx int = 0
		for x := e.x; x < lastXIdx; x++ {
			// Get the correct start of the column index.
			// One char is a group of bytes
			colidx = y*ws.rowSize + x*ws.colSize
			valueIdx := ((y - e.y) * lastXIdx) + (x - e.x)

			// reset column
			copy(canvas[colidx:colidx+ws.colSize], make([]byte, ws.colSize))

			// Set up color if it's first char of the element
			if x == e.x {
				color := e.color.paint()
				colorlen := len(color)

				start := colidx
				end := colidx + colorlen
				copy(canvas[start:end], color[:])

				// if no value set space
				if len(v) <= valueIdx {
					canvas[end] = byte(space)
				} else {
					canvas[end] = v[valueIdx]
				}

				continue
			}

			// add padding `space` if the value is smaller than the width
			if len(v) <= valueIdx {
				canvas[colidx] = byte(space)

				continue
			}

			// set the value
			canvas[colidx] = v[valueIdx]
		}

		// reset color
		nocolor := []byte(noColor)
		start := (colidx + ws.colSize) - len(nocolor)
		end := start + len(nocolor)
		copy(canvas[start:end], nocolor[:])

	}

	return nil
}

type padding struct {
	top    int
	bottom int
	left   int
	right  int
}
