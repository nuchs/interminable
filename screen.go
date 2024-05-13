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
	Value rune
}

type Screen struct {
	width, height int
	capW, capH    int
	cells         [][]Cell
}

func (s *Screen) Width() int {
	return s.width
}

func (s *Screen) Height() int {
	return s.height
}

func NewScreen(w, h int) Screen {
	cells := make([][]Cell, h)
	for i := 0; i < h; i++ {
		cells[i] = make([]Cell, w)
		for j := 0; j < w; j++ {
			cells[i][j] = Cell{' '}
		}
	}

	return Screen{
		width:  w,
		height: h,
		capW:   w,
		capH:   h,
		cells:  cells,
	}
}

func (s *Screen) SetCell(x, y int, value rune) error {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return ErrOutOfBounds
	}

	s.cells[y][x] = Cell{value}

	return nil
}

func (s *Screen) SetRow(col, row int, value string) {
	x, runes := negativeClip(col, value)

	for _, r := range runes {
		err := s.SetCell(x, row, r)
		x += 1
		if err != nil {
			break
		}
	}
}

func (s *Screen) SetCol(col, row int, value string) {
	y, runes := negativeClip(row, value)

	for _, r := range runes {
		err := s.SetCell(col, y, r)
		y += 1
		if err != nil {
			break
		}
	}
}

func (s *Screen) Resize(w, h int) {
	if w > s.capW {
		for i := 0; i < s.height; i++ {
			oldCells := s.cells[i]
			s.cells[i] = make([]Cell, w)
			copy(s.cells[i], oldCells)
		}

		s.capW = w
	}

	if w > s.width {
		for i := 0; i < s.height; i++ {
			for j := s.width; j < w; j++ {
				s.cells[i][j].Value = ' '
			}
		}
	}

	s.width = w

	if h > s.capH {
		oldCells := s.cells
		s.cells = make([][]Cell, h)
		s.capH = h
		copy(s.cells, oldCells)

		for i := s.height; i < h; i++ {
			s.cells[i] = make([]Cell, s.width)
		}
	}

	if h > s.height {
		for i := s.height; i < h; i++ {
			for j := 0; j < s.width; j++ {
				s.cells[i][j].Value = ' '
			}
		}
	}

	s.height = h
}

func (s *Screen) Render() string {
	builder := strings.Builder{}
	builder.WriteString("\033[0;0H")

	for i := 0; i < s.height; i++ {
		for j := 0; j < s.width; j++ {
			builder.WriteRune(s.cells[i][j].Value)
		}

		if i < s.height-1 {
			builder.WriteString("\r\n")
		}
	}

	return builder.String()
}

func negativeClip(ix int, value string) (int, []rune) {
	runes := []rune(value)

	if ix < 0 {
		ix *= -1
		if ix > len(runes) {
			return 0, []rune{}
		}
		runes = runes[ix:]
		ix = 0
	}

	return ix, runes
}
