package window

import (
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
				color:  color{code: FgBlack},
				width:  1,
				height: 1,
				value:  []byte("testing"),
			},
			ws:       winsize{},
			canvas:   []byte{},
			expected: []byte{},
			error:    true,
			errMsg:   "canvas not set",
		},
		{
			name: "Successful draw element with no color",
			el: element{
				width:  1,
				height: 1,
				value:  []byte("testing"),
			},
			ws: winsize{
				rows:    1,
				cols:    1,
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
			name: "Successful draw element color",
			el: element{
				color:  color{code: FgBlack},
				width:  1,
				height: 1,
				value:  []byte("testing"),
			},
			ws: winsize{
				rows:    1,
				cols:    1,
				rowSize: defaultColumnSize,
				colSize: defaultColumnSize,
			},
			canvas: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			// starts with "\033[0;30mt" and ends with \033[0m
			expected: []byte{
				0x1b, 0x5b, 0x30, 0x3b, 0x33, 0x30, 0x6d, 0x74, 0x0, 0x0, 0x0, 0x1b, 0x5b, 0x30, 0x6d,
			},
		},
		{
			name: "Successful draw element with multiple rows",
			el: element{
				color:  color{code: BgBlue},
				width:  8,
				height: 2,
				value:  []byte("testing testing1"),
			},
			ws: winsize{
				rows:    2,
				cols:    8,
				rowSize: defaultColumnSize * 8,
				colSize: defaultColumnSize,
			},
			canvas:   make([]byte, 240),
			expected: []byte{0x1b, 0x5b, 0x30, 0x3b, 0x34, 0x34, 0x6d, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x73, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x69, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x67, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0x5b, 0x30, 0x6d, 0x1b, 0x5b, 0x30, 0x3b, 0x34, 0x34, 0x6d, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x73, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x69, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x67, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0x5b, 0x30, 0x6d},
		},
		{
			name: "Successful draw element with padding",
			el: element{
				color:  color{code: BgBlue},
				width:  8,
				height: 2,
				value:  []byte("testing testing1"),
			},
			ws: winsize{
				rows:    2,
				cols:    8,
				rowSize: defaultColumnSize * 8,
				colSize: defaultColumnSize,
			},
			canvas:   make([]byte, 240),
			expected: []byte{0x1b, 0x5b, 0x30, 0x3b, 0x34, 0x34, 0x6d, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x73, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x69, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x67, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0x5b, 0x30, 0x6d, 0x1b, 0x5b, 0x30, 0x3b, 0x34, 0x34, 0x6d, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x73, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x69, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x67, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0x5b, 0x30, 0x6d},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.el.draw(test.canvas, test.ws); err != nil {
				if !test.error {
					t.Errorf("element draw failed with error: %s", err)
				}

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
