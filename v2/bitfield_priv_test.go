package bitfield

import (
	"testing"
)

func assert(t *testing.T, a, b interface{}) {
	if a == b {
		return
	}
	t.Helper()
	t.Errorf("%s != %s", a, b)
}

func TestPrivate1(t *testing.T) {
	tests := []struct {
		size     int // input
		pos      int //input
		expected int //output
	}{
		{67, -2, 65},
		{121, 121, 0},
		{3, -10, 2},
		{0, 2, 0},
		{0, 4, 0},
	}

	for _, tt := range tests {
		res := New(tt.size).posNormalize(tt.pos)
		assert(t, res, tt.expected)
	}
}

func TestPrivate2(t *testing.T) {
	tests := []struct {
		size        int // input
		offset      int // input
		expectedIx  int // output
		expectedPos int // output
	}{
		{65, 64, 1, 0},
		{3, -1, 0, 2},
		{65, -1, 1, 0},
	}

	for _, tt := range tests {
		ix, p := New(tt.size).posToOffset(tt.offset)
		if ix == tt.expectedIx && p == tt.expectedPos {
			continue
		}
		t.Errorf("New(%d).posToOffset(%d) should map to [%d,%d]. Got: [%d,%d]",
			tt.size, tt.offset, tt.expectedIx, tt.expectedPos, ix, p)
	}
}

