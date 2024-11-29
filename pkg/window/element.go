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
	// padding
	padding bool
}

// element draws itself
func (e element) draw(canvas []byte, ws winsize) error {
	lastXIdx := e.x + e.width  // last x index of the element
	lastYIdx := e.y + e.height // last y index of the element

	// Check if the canvas is set
	if len(canvas) == 0 || ws.rowSize == 0 {
		return errors.New("canvas not set")
	}

	// calculate the number of rows
	rowsNum := len(canvas) / ws.rowSize

	// Check if the element fits in the window
	if rowsNum < lastYIdx {
		return errors.New("element does not fit in the window")
	}

	// value data
	v := []byte(e.value)
	valueIdx := 0
	for y := e.y; y < lastYIdx; y++ {
		// handle top and bottom padding
		if (y == e.y || y == lastYIdx - 1) && e.padding  {
			var colidx int = 0
			for x := e.x; x < lastXIdx; x++ {
				colidx = y*ws.rowSize + x*ws.colSize
				// reset column
				copy(canvas[colidx:colidx+ws.colSize], make([]byte, ws.colSize))

				if x == e.x {
					color := e.color.paint()
					colorlen := len(color)

					start := colidx
					end := colidx + colorlen
					copy(canvas[start:end], color[:])

					canvas[end] = byte(space)
				} else {
					canvas[colidx] = byte(space)
				}
			}

			// reset color
			nocolor := []byte(noColor)
			start := (colidx + ws.colSize) - len(nocolor)
			end := start + len(nocolor)
			copy(canvas[start:end], nocolor[:])
			continue
		}

		// Check if the element fits in the `X` line
		if lastXIdx*ws.colSize-1 > ws.rowSize {
			return errors.New("element does not fit in the X line")
		}

		var colidx int = 0
		for x := e.x; x < lastXIdx; x++ {
			// Get the correct start of the column index.
			// One char is a group of bytes
			colidx = y*ws.rowSize + x*ws.colSize

			// handle the new line char `\n`
			if len(v) > valueIdx && v[valueIdx] == newLine {
				for i := x; i < lastXIdx; i++ {
					colidx = y*ws.rowSize + i*ws.colSize
					canvas[colidx] = byte(space)
				}
				break
			}

			// reset column
			copy(canvas[colidx:colidx+ws.colSize], make([]byte, ws.colSize))

			// Set up color if it's first char of the element
			if x == e.x {
				color := e.color.paint()
				colorlen := len(color)

				start := colidx
				end := colidx + colorlen
				copy(canvas[start:end], color[:])

				// if no value, or left padding set space
				if len(v) <= valueIdx || e.padding {
					canvas[end] = byte(space)
				} else {
					canvas[end] = v[valueIdx]
					valueIdx++
				}

				continue
			}

			// add padding `space` if the value is smaller than the width
			if len(v) <= valueIdx {
				canvas[colidx] = byte(space)

				valueIdx++
				continue
			}

			// if last index and right padding set print space
			if x == lastXIdx - 1 && e.padding {
				canvas[colidx] = byte(space)

				valueIdx++
				continue
			}

			// set the value
			canvas[colidx] = v[valueIdx]
			valueIdx++
		}

		// reset color
		nocolor := []byte(noColor)
		start := (colidx + ws.colSize) - len(nocolor)
		end := start + len(nocolor)
		copy(canvas[start:end], nocolor[:])

		// prevent wrap to the next line until new line appears
		if len(v) > valueIdx {
			w := valueIdx + 1

			for i, b := range v[valueIdx:] {
				valueIdx = w + i
				if b == newLine {
					break
				}
			}
		}
	}

	return nil
}
