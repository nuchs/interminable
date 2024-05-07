package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/nuchs/interminable"
)

type WinSize struct {
	rows, cols, x, y uint16
}

func getTerminalSize(fd uintptr) (WinSize, error) {
	ws := WinSize{}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)))

	if errno != 0 {
		return ws, errno
	}

	return ws, nil
}

func getTermios(fd uintptr) (syscall.Termios, error) {
	var t syscall.Termios
	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCGETS,
		uintptr(unsafe.Pointer(&t)))

	if err != 0 {
		return t, err
	}

	return t, nil
}

func setTermios(fd uintptr, term syscall.Termios) error {
	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCSETS,
		uintptr(unsafe.Pointer(&term)))

	if err != 0 {
		return err
	}

	return nil
}

func setRaw(termios *syscall.Termios) {
	termios.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	termios.Oflag &^= syscall.OPOST
	termios.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	termios.Cflag &^= syscall.CSIZE | syscall.PARENB
	termios.Cflag |= syscall.CS8
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0
}

func main() {
	size, err := getTerminalSize(os.Stdin.Fd())
	if err != nil {
		fmt.Println(err)
		return
	}

	old, err := getTermios(os.Stdin.Fd())
	if err != nil {
		fmt.Println(err)
		return
	}
	termios := old
	setRaw(&termios)
	defer setTermios(os.Stdin.Fd(), old)

	screen := interminable.NewScreen(int(size.cols), int(size.rows))
	x, y := ((screen.Width-1)/2)-6, (screen.Height-1)/2
	screen.SetCell(x, y, 'H')
	screen.SetCell(x+1, y, 'e')
	screen.SetCell(x+2, y, 'l')
	screen.SetCell(x+3, y, 'l')
	screen.SetCell(x+4, y, 'o')
	screen.SetCell(x+5, y, ' ')
	screen.SetCell(x+6, y, 'W')
	screen.SetCell(x+7, y, 'o')
	screen.SetCell(x+8, y, 'r')
	screen.SetCell(x+9, y, 'l')
	screen.SetCell(x+10, y, 'd')

	s := screen.Render()
	fmt.Print(s)

	time.Sleep(5 * time.Second)
	fmt.Println("Ok, I love you, byebye")
}
