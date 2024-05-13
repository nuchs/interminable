package interminable

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

type WinSize struct {
	Rows, Cols       uint16
	unused1, unused2 uint16 // these are just here to make the struct the right size for the syscall
}

type Terminal struct {
	Screen   Screen
	original syscall.Termios
	fd       uintptr
	closed   chan struct{}
	subs     []chan WinSize
}

func (t *Terminal) Open(fd uintptr) error {
	var err error
	t.original, err = getTermios(fd)
	if err != nil {
		return fmt.Errorf("Failed to get termios: %w", err)
	}

	termios := t.original
	setRaw(&termios)
	if err := setTermios(fd, termios); err != nil {
		return fmt.Errorf("Failed to set raw mode: %w", err)
	}

	t.fd = fd
	ws, err := getTerminalSize(fd)
	if err != nil {
		err := fmt.Errorf("Failed to get terminal size: %w", err)
		if err2 := t.Close(); err2 != nil {
			return fmt.Errorf("Multiple failures: %w, %w", err, err2)
		}
		return err
	}
	t.Screen = NewScreen(int(ws.Cols), int(ws.Rows))
	t.closed = make(chan struct{})

	go t.eventLoop()

	return nil
}

func (t *Terminal) Close() error {
	if err := setTermios(t.fd, t.original); err != nil {
		return fmt.Errorf("Failed to restore terminal: %w", err)
	}

	t.closed <- struct{}{}

	return nil
}

func (t *Terminal) Refresh() {
	frame := t.Screen.Render()
	fmt.Print(frame)
}

func (t *Terminal) SubscribeToResizes(c chan WinSize) {
	t.subs = append(t.subs, c)
}

func (t *Terminal) Fd() uintptr {
	return t.fd
}

func (t *Terminal) eventLoop() {
	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)

EventLoop:
	for {
		select {
		case <-t.closed:
			break EventLoop
		case <-winch:
			ws, err := getTerminalSize(t.fd)
			if err != nil {
				continue
			}
			t.Screen.Resize(int(ws.Cols), int(ws.Rows))
			for _, c := range t.subs {
				c <- ws
			}
			break
		}
	}

	close(t.closed)
}

func getTerminalSize(fd uintptr) (WinSize, error) {
	ws := WinSize{}

	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)))

	if err != 0 {
		return ws, err
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
