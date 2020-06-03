package bitfield

import (
	"fmt"
	"testing"
)

// test that function call do panic
func panics(f func()) bool {
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		f()
	}()
	return didPanic
}

func Test1(t *testing.T) {
	if !panics(func() { NewBitField(-3) }) {
		t.Error("should panic")
	}

	if !panics(func() { New(4).Resize(-1) }) {
		t.Error("should panic")
	}

	if NewBitField(0).Len() != 0 {
		t.Error("should be 0")
	}
	if New(65).SetAll().OnesCount() != 65 {
		t.Error("should be 65")
	}
	if New(3).Equal(New(4)) {
		t.Error("should be false")
	}
	if !New(3).Set(0, -1).Not().Equal(New(3).Set(1)) {
		t.Error("should be 1")
	}
	a := New(129).Set(0, -1).Clear(123, -3).Not().Not()
	if a.OnesCount() != 2 {
		t.Error("should be 2")
	}
	if !a.Get(0) || !a.Get(-a.Len()) || !a.Get(a.Len()) {
		t.Error("should be true")
	}

	b := New(129).Set(4, -1, 23, 11)
	if b.And(New(129).Set(0, -1)).OnesCount() != 1 {
		t.Error("should be 1")
	}

	if !panics(func() { New(5).And(New(121)) }) {
		t.Error("should panic")
	}
	if !panics(func() { New(5).Or(New(121)) }) {
		t.Error("should panic")
	}
	if !panics(func() { New(5).Xor(New(121)) }) {
		t.Error("should panic")
	}

	c := New(129).Set(73, -2).ClearAll().Set(-1)
	if !c.Equal(New(129).Set(128)) {
		t.Error("should be equal")
	}
	if c.Equal(New(129).Not()) {
		t.Error("should be not equal")
	}
	if c.Get(127) {
		t.Error("should be false")
	}
	if !c.Get(-1) {
		t.Error("should be true")
	}

	if !New(4).Flip(-1).Equal(New(4).Set(-1)) {
		t.Error("should be equal")
	}
	if !New(4).Flip(-1).Flip(-1).Equal(New(4)) {
		t.Error("should be equal")
	}
	if New(4).SetAll().Flip(0, -1).Mid(1, 2).OnesCount() != 2 {
		t.Error("should be 2")
	}

	d := New(65).Set(-1)
	e := d.Clone()
	if !d.Equal(e) {
		t.Error("should be equal")
	}
	if d.Xor(d).OnesCount() != 0 {
		t.Error("should be 0")
	}

	if !d.Set(11).Or(e).Get(11) {
		t.Error("should be true")
	}
}

func Test2(t *testing.T) {
	if New(4).Resize(0).Len() != 0 {
		t.Error("should be 0")
	}
	if New(10).SetAll().Resize(6).OnesCount() != 6 {
		t.Error("should be 6")
	}
	if !New(27).Set(0, -1).Resize(65).Get(26) {
		t.Error("should be true")
	}
	if New(65).Set(-1).Resize(45).OnesCount() != 0 {
		t.Error("should be 0")
	}
	if New(65).SetAll().Resize(40).OnesCount() != 40 {
		t.Error("should be 40")
	}
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
		if res == tt.expected {
			continue
		}
		t.Errorf("New(%d).posNormalize(%d) should map to %d. Got: %d",
			tt.size, tt.pos, tt.expected, res)
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

func TestShift(t *testing.T) {
	tests := []struct {
		size     int
		shift    int
		expected int
	}{
		{100, 54, 46},
		{129, -3, 126},
		{6, 2, 4},
		{6, 0, 6},
		{0, 5, 0},
		{193, 192, 1},
		{6, 6, 0},
		{6, -6, 0},
		{193, -192, 1},
		{193, -1, 192},
		{191, -100, 91},
	}

	for _, tt := range tests {
		c := New(tt.size).SetAll().Shift(tt.shift).OnesCount()
		if c != tt.expected {
			t.Errorf("New(%d).SetAll().Shift(%d) had %d bits set, however %d was expected",
				tt.size, tt.shift, c, tt.expected)
		}
	}
}

func TestMid(t *testing.T) {
	if !New(5).SetAll().Mid(1, 1).Equal(New(1).Set(0)) {
		t.Error("should be equal")
	}
	if New(121).SetAll().Mid(-3, 3).OnesCount() != 3 {
		t.Error("should be 3")
	}

	if !panics(func() {
		New(5).Left(-1)
	}) {
		t.Error("should panic")
	}

	if !panics(func() {
		New(5).Right(-1)
	}) {
		t.Error("should panic")
	}

	if !panics(func() {
		New(10).Mid(3, -1)
	}) {
		t.Error("should panic")
	}

	if !New(65).Set(3).Left(3).Equal(New(3)) {
		t.Error("should be equal")
	}

	a := New(60).SetAll().Right(10)
	b := New(10).SetAll()

	if !a.Equal(b) {
		t.Error("should be equal")
	}

	a = New(10).SetAll().Right(11)
	a = New(10).SetAll().Left(11)
	if !a.Equal(b) {
		t.Error("should be equal")
	}

	if New(4).Left(0).Len() != 0 {
		t.Error("should be 0")
	}
	if New(4).SetAll().Mid(-3, 3).OnesCount() != 3 {
		t.Error("should be 3")
	}
}

func TestAppend(t *testing.T) {
	// trivial cases
	a := New(0).Append(New(0))
	if !a.Equal(New(0)) {
		t.Error("should be equal")
	}
	a = New(0).Append(New(1))
	if !a.Equal(New(1)) {
		t.Error("should be equal")
	}

	// real cases
	a = New(10).SetAll().Append(New(3))
	if a.Len() != 13 || a.OnesCount() != 10 || a.Right(3).OnesCount() != 0 {
		t.Error("Append is wrong")
	}
}

func TestRotate(t *testing.T) {

	a := New(65).Set(63).Rotate(1)
	if !a.Equal(New(65).Set(64)) {
		t.Error("should be equal")
	}

	a = New(65).Set(0).Rotate(-1)
	if !a.Equal(New(65).Set(64)) {
		t.Error("should be equal")
	}

	const len = 163
	a = New(len).Set(0)
	for i := -len * 2; i < len*2; i++ {
		r := a.Clone().Rotate(i)
		if !r.Equal(New(len).Set(i)) {
			t.Errorf("@%d rotate failed", i)
		}
	}
}

func TestMut(t *testing.T) {
	a := New(65).Mut()
	a.SetAll().Clear(0, -1).Flip(3, 4)
	if a.OnesCount() != 61 {
		t.Error("should be 61")
	}

	a = New(65).Mut().SetAll()
	a.Xor(a)
	if a.OnesCount() != 0 {
		t.Error("should be 0")
	}

	a = New(65)
	if !panics(func() {
		New(5).Copy(a)
	}) {
		t.Error("should panic")
	}

	New(65).Set(0, -1).Copy(a)
	if a.OnesCount() != 2 || !a.Get(-1) {
		t.Error("Copy fails")
	}
}

func Benchmark1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New(365).Mut().SetAll()
	}
}

func ExampleBitField_Shift_e1() {
	bf := NewBitField(3).Set(0).Shift(1)
	fmt.Printf("%t\n", bf.Get(1))
	// Output: true
}

func ExampleBitField_Shift_e2() {
	bf := NewBitField(3).Set(0).Shift(3)
	fmt.Printf("%d\n", bf.OnesCount())
	// Output: 0
}

func ExampleBitField_Clear() {
	bf := NewBitField(4).SetAll().Clear(0, -1)
	fmt.Printf("%d\n", bf.OnesCount())
	// Output: 2
}

func ExampleBitField_Mut() {
	bf := NewBitField(4).Set(0).Mut()
	bf.Set(1)
	fmt.Printf("with Mut: %d\n", bf.OnesCount())

	bf2 := NewBitField(4).Set(0)
	bf2.Set(1) // this is set then discarded!
	fmt.Printf("without Mut: %d\n", bf2.OnesCount())
	// Output: with Mut: 2
	// without Mut: 1

}
