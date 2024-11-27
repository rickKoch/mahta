package window

import (
  "fmt"
	"testing"
)

func TestDrawElement(t *testing.T) {
	tests := []struct {
		name     string
		el       element
    ws       winsize
		canvas   []byte
		expected []byte
		errMsg   string
		error    bool
	}{
		{
			name: "Error on empty canvas",
			el: element{
				color: color{code: FgBlack},
        width: 1,
        height: 1,
        value: []byte("testing"),
			},
      ws: winsize{},
			canvas:   []byte{},
			expected: []byte{},
			error:    true,
			errMsg:   "canvas not set",
		},
		{
			name: "Successfull draw element with no color",
			el: element{
				width: 1,
        height: 1,
        value: []byte("testing"),
			},
      ws: winsize{
        rows: 1,
        cols: 1,
        rowSize: defaultColumnSize,
        colSize: defaultColumnSize,
      },
			canvas: make([]byte, defaultColumnSize),
			// starts with "\033[0mt" and ends with \033[0m
			expected: []byte{
        0x1b, 0x5b, 0x30, 0x6d, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0x5b, 0x30, 0x6d,
      },
		},
		{
			name: "Successfull draw element color",
			el: element{
				color: color{code: FgBlack},
        width: 1,
        height: 1,
        value: []byte("testing"),
			},
      ws: winsize{
        rows: 1,
        cols: 1,
        rowSize: defaultColumnSize,
        colSize: defaultColumnSize,
      },
			canvas: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			// starts with "\033[0;30mt" and ends with \033[0m
			expected: []byte{
        0x1b, 0x5b, 0x30, 0x3b, 0x33, 0x30, 0x6d, 0x74, 0x0, 0x0, 0x0, 0x1b, 0x5b, 0x30, 0x6d,
      },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.el.draw(test.canvas, test.ws); err != nil {
				if !test.error {
					t.Errorf("element draw failed with error: %s", err)
				}

        fmt.Println(test.errMsg, err.Error())
				if test.errMsg != err.Error() {
					t.Errorf("element draw failed with different error: %s", err)
				}
			}

			if !testCanvasEq(test.expected, test.canvas) {
				t.Errorf(
					"canvas: %#v not equal to expected %#v", test.canvas, test.expected,
				)
			}
		})
	}
}

func testCanvasEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
    if a[i] != b[i] {
      return false
    }
	}

	return true
}
