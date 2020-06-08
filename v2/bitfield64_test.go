package bitfield

import (
	"fmt"
	"testing"
)

func TestPos(t *testing.T) {
	tests := []struct {
		in  int
		out uint64
	}{
		{0, 1 << 0},
		{-1, 1 << 63},
		{-64, 1 << 0},
		{3, 1 << 3},
		{65, 1 << 1},
	}
	for _, tt := range tests {
		if posToBitMask(tt.in) != BitField64(tt.out) {
			t.Errorf("got %d, wants %d", tt.in, tt.out)
		}
	}
}

func Test64_1(t *testing.T) {
	bf1 := New64().SetAll().Clear(3)
	if bf1.OnesCount() != 63 {
		t.Error("should be 63")
	}
	if bf1.Not().OnesCount() != 1 {
		t.Error("should be 1")
	}
	if bf1.Get(3) {
		t.Error("should be false")
	}
	if !bf1.Get(17) {
		t.Error("should be true")
	}
	bf2 := New64().Set(3).Set(1).Clear(5).Or(New64().Set(4))
	if bf2.StringPretty() != "01011" {
		t.Error("should be 01011")
	}

	if bf2.OnesCount() != 3 {
		t.Error("should be 3")
	}
	if bf1.Xor(bf1).OnesCount() != 0 {
		t.Error("should be zero")
	}
	if bf2.And(bf1).OnesCount() != 2 {
		t.Error("should be 2")
	}

	if New64().SetAll().ClearAll().OnesCount() != 0 {
		t.Error("should be zero")
	}
	if New64().Set(-63-64).StringPretty() != "01" {
		t.Error("should be 01")
	}

	if New64().Flip(3) != 8 {
		t.Error("shoud be 8")
	}
	if New64().Flip(4).Flip(4) != 0 {
		t.Error("should be 0")
	}
}

func Test64_2(t *testing.T) {
	a := New64().SetMul(1, 3)
	if a.Mid(0, 4).OnesCount() != 2 {
		t.Error("should be 2")
	}
	if a.Mid(1, 1).OnesCount() != 1 {
		t.Error("should be 1")
	}
	if a.Mid(3, 5).OnesCount() != 1 {
		t.Error("should be 1")
	}
	if a.Mid(4, 5).OnesCount() != 0 {
		t.Error("should be 0")
	}
	a = New64().SetMul(0, 1, 4)
	if a.Left(4).OnesCount() != 2 {
		t.Error("should be 2")
	}
	if a.Left(1).OnesCount() != 1 {
		t.Error("should be 1")
	}
	if !New64().Set(63).Right(3).Get(2) {
		t.Error("should be true")
	}

	if !a.Rotate(1).Get(5) {
		t.Error("should be true")
	}
	if !a.Rotate(-1).Get(3) {
		t.Error("should be true")
	}
	if !New64().Set(5).Rotate(0).Rotate(1).Rotate(-1).Get(5) {
		t.Error("should be true")
	}

	if !New64().Set(0).Shift(1).Get(1) {
		t.Error("should be true")
	}
	if !New64().Set(1).Shift(-1).Get(0) {
		t.Error("should be true")
	}
	if New64().SetAll().Shift(63).OnesCount() != 1 {
		t.Error("should be 1")
	}
	if New64().SetAll().Shift(-65).OnesCount() != 0 {
		t.Error("should be 0")
	}
	a = New64().Set(0).Set(5).Shift(0).Shift(-1).Shift(1)
	if !a.Get(5) {
		t.Error("should be true")
	}
	if a.Get(0) {
		t.Error("should be false")
	}

	a = New64().SetAll().ClearMul(0, -1)
	if a.Get(0) || a.Get(-1) || a.OnesCount() != 62 {
		t.Error("ClearMul fails")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	New64().Mid(0, -1)
}

func TestShift2(t *testing.T) {

	orig := New64().SetMul(0, -1)

	// case 0
	a, b := orig.Shift2(0)
	if a != orig && b != 0 {
		t.Error("Shift2(0) should leave values intact")
	}

	// case 0 < n  < 64
	a, b = orig.Shift2(1)
	if !a.Get(1) {
		t.Error("Shift2(1): 'a' should be true")
	}
	if !b.Get(0) {
		t.Error("Shift2(1): 'b' should be true")
	}

	// case -63 < n  < 0
	a, b = orig.Shift2(-1)
	if a != New64().Set(62) {
		t.Error("Shift2(1): 'a' should be true")
	}
	if b != New64().Set(63) {
		t.Error("Shift2(1): 'b' should be true")
	}

	// case n = 64
	a, b = orig.Shift2(64)
	if a.OnesCount() != 0 || b != orig {
		t.Error("Should be 0 and same")
	}
	// case n = -64
	a, b = orig.Shift2(-64)
	if a.OnesCount() != 0 || b != orig {
		t.Error("Should be 0 and same")
	}

	// case n > 64
	a, b = orig.Shift2(65)
	if a != 0 || b != 0 {
		t.Error("'a' and 'b should be 0")
	}

	// case n < -64
	a, b = orig.Shift2(-65)
	if a != 0 || b != 0 {
		t.Error("'a' and 'b should be 0")
	}

}

func ExampleBitField64_SetMul() {
	a := New64().SetMul(2, 4)
	fmt.Println(a.StringPretty())
	// Output: 00101
}

func ExampleBitField64_Xor() {
	a := New64().SetMul(0, 3)
	b := New64().Set(1)
	fmt.Println(a.Xor(b).StringPretty())
	// Output: 1101
}

func ExampleBitField64_Flip() {
	a := New64().SetAll().Flip(0)
	fmt.Println(a.String())
	// Output: 0111111111111111111111111111111111111111111111111111111111111111
}

func ExampleBitField64_Shift() {
	// SetAll(): all bits are 1
	// Shift(3): 3 zeroes enter from left
	// Left(5): takes first 3 zero bits and two 1s.
	a := New64().SetAll().Shift(3).Left(5)
	fmt.Println(a.StringPretty())
	// Output: 00011
}
