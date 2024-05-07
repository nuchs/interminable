package interminable_test

import (
	"testing"

	"github.com/nuchs/interminable"
)

func TestNewScreen(t *testing.T) {
	s := interminable.NewScreen(10, 10)
	if s.Width != 10 || s.Height != 10 {
		t.Error("wrong screen dimensions")
	}
}

func TestBadSetCell(t *testing.T) {
	err := "out of bounds"

	testCases := []struct {
		desc string
		x    int
		y    int
		err  string
	}{
		{desc: "x < 0", x: -1, y: 0, err: err},
		{desc: "y < 0", x: 0, y: -1, err: err},
		{desc: "x >= w", x: 10, y: 0, err: err},
		{desc: "y >= h", x: 0, y: 10, err: err},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := interminable.NewScreen(5, 5)
			err := s.SetCell(tc.x, tc.y, ' ')
			if err.Error() != tc.err {
				t.Errorf("%s: expected error %v, got %v", tc.desc, tc.err, err)
			}
		})
	}
}

func TestRender(t *testing.T) {
	s := interminable.NewScreen(3, 2)
	s.SetCell(0, 0, 'a')
	s.SetCell(1, 1, 'b')
	s.SetCell(2, 0, 'c')

	result := s.Render()

	expected := "\033[0;0Ha c\r\n b "

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
