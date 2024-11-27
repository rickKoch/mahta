package sys

import (
	"syscall"
	"unsafe"
)

// TIOCGWINSZ provides the window size.
func TIOCGWINSZ(fd uintptr) (int, int, error) {
	// window size corresponding struct
	ws := struct {
		rows uint16
		cols uint16
		_pad [4]byte
	}{}

	// window size system call
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
	if errno != 0 {
		return 0, 0, errno
	}

	return int(ws.rows), int(ws.cols), nil
}
