package interminable

import "strings"

type TermError string

func (e TermError) Error() string {
	return string(e)
}

const (
	ErrOutOfBounds = TermError("out of bounds")
)

type Cell struct {
	X, Y  int
	Value rune
}

type Screen struct {
	Width  int
	Height int
	cells  [][]Cell
}

func NewScreen(w, h int) Screen {
	cells := make([][]Cell, h)
	for i := 0; i < h; i++ {
		cells[i] = make([]Cell, w)
		for j := 0; j < w; j++ {
			cells[i][j] = Cell{i, j, ' '}
		}
	}

	return Screen{
		Width:  w,
		Height: h,
		cells:  cells,
	}
}

func (s *Screen) SetCell(x, y int, value rune) error {
	if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
		return ErrOutOfBounds
	}

	s.cells[y][x] = Cell{x, y, value}

	return nil
}

func (s *Screen) Render() string {
	builder := strings.Builder{}

	builder.WriteString("\033[0;0H")

	for i := 0; i < s.Height; i++ {
		for j := 0; j < s.Width; j++ {
			builder.WriteRune(s.cells[i][j].Value)
		}

		if i < s.Height-1 {
			builder.WriteString("\r\n")
		}
	}

	return builder.String()
}