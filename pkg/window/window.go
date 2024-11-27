package window

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rickKoch/mahta/sys"
)

const (
	// ASCII space char = 32
	space = 32
	// Default column size
	defaultColumnSize = 15
)

type window struct {
	size     winsize
	elements elements
	canvas   []byte
}

func New() (*window, error) {
	size, err := getWindowSize(os.Stdout.Fd())
	if err != nil {
		return nil, err
	}

	// generate canvas from the window size
	canvas := size.canvas()

	// create the window
	w := window{
		size:   size,
		canvas: canvas,
		elements: []element{
			{
				width:  6,
				height: 8,
				value:  []byte("pero, trpe, mite, risto"),
				color:  color{code: BgGreen},
			},
			{
				width:  1,
				height: 1,
				value:  []byte("m"),
				x:      20,
				color:  color{code: BgHiMagenta},
			},
			{
				width:  6,
				height: 1,
				value:  []byte("trpe, testing"),
				x:      3,
				y:      4,
				color:  color{code: FgBlue},
			},
		},
	}

	return &w, nil
}

// Render draws and updates the window
func (w *window) Render(ctx context.Context) error {
	ws := make(chan winsize)

	if err := windowSizeChanges(ctx, ws); err != nil {
		return err
	}

	if err := w.draw(); err != nil {
		return err
	}

	for {
		select {
		case winsize := <-ws:
			fmt.Printf("WindowSizeChanges: %#v\n", winsize)
			w.size = winsize
		case <-ctx.Done():
			return nil
		}
	}
}

// draw the window on a file
func (w *window) draw() error {
	if err := w.elements.draw(w.canvas, w.size); err != nil {
		return err
	}

	if _, err := os.Stdout.Write(w.canvas); err != nil {
		return err
	}

	return nil
}

type winsize struct {
	// number of rows based on the size of the terminal
	rows int
	// number of columns based on the size of the terminal
	cols int
	// number of bytes in one row
	rowSize int
	// number of bytes in one column
	colSize int
}

func (ws winsize) canvas() []byte {
	// total number of bytes on the canvas
	total := ws.rowSize * ws.rows

	canv := make([]byte, total)
	for i := 0; i < total; i += ws.colSize {
		canv[i] = byte(space)
	}

	return canv
}

func getWindowSize(fd uintptr) (winsize, error) {
	ws := winsize{colSize: defaultColumnSize}
	var err error

	ws.rows, ws.cols, err = sys.TIOCGWINSZ(fd)

	// columns are created from multiple bytes
	ws.rowSize = ws.cols * ws.colSize

	return ws, err
}

func windowSizeChanges(ctx context.Context, ws chan winsize) error {
	ch := make(chan os.Signal, 1)

	// window resize signal
	sig := syscall.SIGWINCH
	signal.Notify(ch, sig)

	go func() {
		for {
			select {
			case <-ch:
				var err error
				s, err := getWindowSize(os.Stdout.Fd())
				if err == nil {
					ws <- s
				}
			case <-ctx.Done():
				signal.Reset(sig)
				close(ch)
				return
			}
		}
	}()

	return nil
}
