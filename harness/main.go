package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nuchs/interminable"
)

func main() {
	run()
	time.Sleep(3 * time.Second)
	fmt.Println("--------------------------------------------------")
	fmt.Println("Ok I love you bye bye!")
}

func run() error {
	term := interminable.Terminal{}
	if err := term.Open(os.Stdin.Fd()); err != nil {
		fmt.Println(err)
		return err
	}
	defer term.Close()

	x, y := ((term.Screen.Width-1)/2)-6, 5
	term.Screen.SetRow(x, y, "Hello World")
	term.Screen.SetCol(x, y, "Hello World")

	s := term.Screen.Render()
	fmt.Print(s)
	return nil
}
