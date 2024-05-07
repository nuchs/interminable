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

	x, y := ((term.Screen.Width-1)/2)-6, (term.Screen.Height-1)/2
	term.Screen.SetCell(x, y, 'H')
	term.Screen.SetCell(x+1, y, 'e')
	term.Screen.SetCell(x+2, y, 'l')
	term.Screen.SetCell(x+3, y, 'l')
	term.Screen.SetCell(x+4, y, 'o')
	term.Screen.SetCell(x+5, y, ' ')
	term.Screen.SetCell(x+6, y, 'W')
	term.Screen.SetCell(x+7, y, 'o')
	term.Screen.SetCell(x+8, y, 'r')
	term.Screen.SetCell(x+9, y, 'l')
	term.Screen.SetCell(x+10, y, 'd')

	s := term.Screen.Render()
	fmt.Print(s)
	return nil
}
