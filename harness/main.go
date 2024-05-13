package main

import (
	"fmt"
	"os"

	"github.com/nuchs/interminable"
)

func main() {
	run()
	fmt.Println("--------------------------------------------------")
	fmt.Println("Ok I love you bye bye!")
}

func run() {
	term := interminable.Terminal{}
	if err := term.Open(os.Stdin.Fd()); err != nil {
		fmt.Println(err)
		return
	}
	defer term.Close()

	notify := make(chan interminable.WinSize)
	term.SubscribeToResizes(notify)
	count := 0

	for {
		Draw(&term)
		<-notify
		count++

		if count == 3 {
			break
		}
	}
}

func Draw(term *interminable.Terminal) {
	x, y := ((term.Screen.Width()-1)/2)-6, 5
	msg := fmt.Sprintf("Hello World (%d, %d)", term.Screen.Width(), term.Screen.Height())
	term.Screen.SetRow(x, y, msg)
	term.Refresh()
}
