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

func TestSetCell(t *testing.T) {
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

func TestSetRow(t *testing.T) {
	testCases := []struct {
		desc     string
		row, col int
		value    string
		expected string
	}{
		{desc: "full top row", row: 0, col: 0, value: "aa", expected: "\033[0;0Haa\r\n  \r\n  "},
		{desc: "full middle row", row: 1, col: 0, value: "aa", expected: "\033[0;0H  \r\naa\r\n  "},
		{desc: "full bottom row", row: 2, col: 0, value: "aa", expected: "\033[0;0H  \r\n  \r\naa"},
		{desc: "clip left", row: 0, col: -1, value: "aa", expected: "\033[0;0Ha \r\n  \r\n  "},
		{desc: "clip right", row: 0, col: 1, value: "aa", expected: "\033[0;0H a\r\n  \r\n  "},
		{desc: "clip left fully", row: 0, col: -9, value: "aa", expected: "\033[0;0H  \r\n  \r\n  "},
		{desc: "clip right fully", row: 0, col: 9, value: "aa", expected: "\033[0;0H  \r\n  \r\n  "},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := interminable.NewScreen(2, 3)
			s.SetRow(tc.col, tc.row, tc.value)
			result := s.Render()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestSetCol(t *testing.T) {
	testCases := []struct {
		desc     string
		row, col int
		value    string
		expected string
	}{
		{desc: "full left col", row: 0, col: 0, value: "aa", expected: "\033[0;0Ha  \r\na  "},
		{desc: "full middle col", row: 0, col: 1, value: "aa", expected: "\033[0;0H a \r\n a "},
		{desc: "full right col", row: 0, col: 2, value: "aa", expected: "\033[0;0H  a\r\n  a"},
		{desc: "clip top", row: -1, col: 0, value: "aa", expected: "\033[0;0Ha  \r\n   "},
		{desc: "clip bottom", row: 1, col: 0, value: "aa", expected: "\033[0;0H   \r\na  "},
		{desc: "clip top fully", row: -9, col: 0, value: "aa", expected: "\033[0;0H   \r\n   "},
		{desc: "clip bottom fully", row: 9, col: 0, value: "aa", expected: "\033[0;0H   \r\n   "},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := interminable.NewScreen(3, 2)
			s.SetCol(tc.col, tc.row, tc.value)
			result := s.Render()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestResize(t *testing.T) {
	testcases := []struct {
		desc     string
		w, h     int
		expected string
	}{
		{desc: "wider", w: 3, h: 2, expected: "\033[0;0Hab \r\ncd "},
		{desc: "narrow", w: 1, h: 2, expected: "\033[0;0Ha\r\nc"},
		{desc: "taller", w: 2, h: 3, expected: "\033[0;0Hab\r\ncd\r\n  "},
		{desc: "shorter", w: 2, h: 1, expected: "\033[0;0Hab"},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			s := interminable.NewScreen(2, 2)
			s.SetRow(0, 0, "ab")
			s.SetRow(0, 1, "cd")
			s.Resize(tc.w, tc.h)
			result := s.Render()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
