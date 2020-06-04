/*
Package bitfield is slice of bitfield64-s to make it possible to store more
than 64 bits. Most functions are chainable, positions outside the [0,len) range
will get the modulo treatment, so Get(len) will return the 0th bit, Get(-1) will
return the last bit: Get(len-1)

Most methods do not modify the underlying bitfield but create a new and return that.
You can change this behaviour by calling .Mut() method. In this case all methods
explicitely marked as 'Mutable.' will be modified in-place. This reduced allocations
(for cases where speed does matter).
*/
package bitfield

import (
	bf64 "github.com/bukshee/bitfield64"
)

// bitFieldData is a slice of BitField64-s.
type bitFieldData []bf64.BitField64

// BitField is a flexible size version of BitField64.
type BitField struct {
	data    bitFieldData
	len     int
	mutable bool
}

// New creates a new BitField of length len
func New(len int) *BitField {
	return NewBitField(len)
}

// NewBitField creates a new BitField of length len and returns it.
// Returns nil if len<0
func NewBitField(len int) *BitField {
	if len < 0 {
		panic("len cannot be negative")
	}
	return &BitField{
		data:    make(bitFieldData, 1+len/64),
		len:     len,
		mutable: false,
	}
}

// Mut sets the mutable flag. This can reduce number of copying
// if execution time is important. Methods where description contains
// 'Mutable.' will modify content in-place.
func (bf *BitField) Mut() *BitField {
	bf.mutable = true
	return bf
}

// Resize resizes the bitfield to newLen in size.
// Returns a newly allocated one, leaves the original intact.
// If newLen < Len() bits are lost at the end.
// If newLen > Len() the newly added bits will be zeroed.
func (bf *BitField) Resize(newLen int) *BitField {
	ret := NewBitField(newLen)
	if newLen == 0 {
		return ret
	}
	copy(ret.data, bf.data)
	if newLen < bf.len {
		ret.clearEnd()
	}
	return ret
}

// mClone is the mutable-clone: if Mut() set, it returns bf,
// otherwise bahaves as Clone()
func (bf *BitField) mClone() *BitField {
	if bf.mutable {
		return bf
	}
	return bf.Clone()
}

// Clone creates a copy of the bitfield and returns it
func (bf *BitField) Clone() *BitField {
	ret := New(bf.len)
	copy(ret.data, bf.data)
	return ret
}

// Copy copies the content of BitField bf to dest.
// Returns false if the two bitfields differ in size, true otherwise
func (bf *BitField) Copy(dest *BitField) bool {
	if bf.len != dest.len {
		panic("Len() of bf and dest differ")
	}
	copy(dest.data, bf.data)
	return true
}

// Len returns the number of bits the BitField holds
func (bf *BitField) Len() int {
	return bf.len
}

func (bf *BitField) posNormalize(pos int) int {
	if bf.len == 0 {
		return 0
	}
	for pos < 0 {
		pos += bf.len
	}
	pos %= bf.len
	return pos
}

func (bf *BitField) posToOffset(pos int) (index int, offset int) {
	pos = bf.posNormalize(pos)
	index = pos / 64
	offset = pos % 64
	return
}

// clearEnd zeroes the bits beyond Len() in-place
// The underlying BitField64 allocates space in 64bit increments
// and Len() might be smaller than the space allocated: it needs to be
// kept zeroed at all times to be consistent
func (bf *BitField) clearEnd() *BitField {
	const n = 64
	index, offset := bf.Len()/n, bf.Len()%n
	// offset points to after the last element
	delta := n - offset
	bf.data[index] = bf.data[index].Shift(delta).Shift(-delta)
	return bf
}

// Set sets the bit(s) at position pos. Mutable.
func (bf *BitField) Set(pos ...int) *BitField {
	ret := bf.mClone()
	for _, p := range pos {
		index, offset := bf.posToOffset(p)
		ret.data[index] = ret.data[index].Set(offset)
	}
	return ret
}

// SetAll sets all bits to 1. Mutable.
func (bf *BitField) SetAll() *BitField {
	ret := bf.mClone()
	for i := range ret.data {
		ret.data[i] = ret.data[i].SetAll()
	}
	return ret.clearEnd()
}

// Clear clears the bit(s) at position pos. Mutable.
func (bf *BitField) Clear(pos ...int) *BitField {
	ret := bf.mClone()
	for _, p := range pos {
		index, offset := bf.posToOffset(p)
		ret.data[index] = ret.data[index].Clear(offset)
	}
	return ret
}

// ClearAll sets all bits to 1. Mutable.
func (bf *BitField) ClearAll() *BitField {
	ret := bf.mClone()
	for i := range ret.data {
		ret.data[i] = ret.data[i].ClearAll()
	}
	return ret
}

// Get returns the bit (as a boolean) at position pos
func (bf *BitField) Get(pos int) bool {
	index, offset := bf.posToOffset(pos)
	return bf.data[index].Get(offset)
}

// Flip inverts the bit(s) at position pos. Mutable.
func (bf *BitField) Flip(pos ...int) *BitField {
	ret := bf.mClone()
	for _, p := range pos {
		index, offset := ret.posToOffset(p)
		ret.data[index] = ret.data[index].Flip(offset)
	}
	return ret
}

// OnesCount returns the number of bits set
func (bf *BitField) OnesCount() int {
	count := 0
	for i := range bf.data {
		count += bf.data[i].OnesCount()
	}
	return count
}

const errLenOther = "Len() of bf and bfOther differ"

// And does a binary AND with bfOther. Returns nil if lengths differ. Mutable.
func (bf *BitField) And(bfOther *BitField) *BitField {
	if bf.len != bfOther.len {
		panic(errLenOther)
	}
	ret := bf.mClone()
	for i := range ret.data {
		ret.data[i] = ret.data[i].And(bfOther.data[i])
	}
	return ret
}

// Or does a binary OR with bfOther. Returns nil if lengths differ. Mutable.
func (bf *BitField) Or(bfOther *BitField) *BitField {
	if bf.len != bfOther.len {
		panic(errLenOther)
	}
	ret := bf.mClone()
	for i := range ret.data {
		ret.data[i] = ret.data[i].Or(bfOther.data[i])
	}
	return ret
}

// Not does a binary NOT (inverts all bits). Mutable.
func (bf *BitField) Not() *BitField {
	ret := bf.mClone()
	for i := range bf.data {
		ret.data[i] = ret.data[i].Not()
	}
	return ret.clearEnd()
}

// Xor does a binary XOR with bfOther. Returns nil if lengths differ. Mutable.
func (bf *BitField) Xor(bfOther *BitField) *BitField {
	if bf.len != bfOther.len {
		panic(errLenOther)
	}
	ret := bf.mClone()
	for i := range bf.data {
		ret.data[i] = ret.data[i].Xor(bfOther.data[i])
	}
	return ret.clearEnd()
}

// Equal tells if two bitfields are equal or not
func (bf *BitField) Equal(bfOther *BitField) bool {
	if bf.len != bfOther.len {
		return false
	}
	for i := range bf.data {
		if bf.data[i] != bfOther.data[i] {
			return false
		}
	}
	return true
}

// Shift shifts thes bitfield by count bits and returns it.
// If count is positive it shifts towards higher bit positions;
// If negative it shifts towards lower bit positions.
// Bits exiting at one end are discarded;
// bits entering at the other end are zeroed. Mutable.
func (bf *BitField) Shift(count int) *BitField {
	ret := bf.mClone()
	if count <= -bf.Len() || count >= bf.Len() {
		return ret.ClearAll()
	}

	const n = 64
	switch {
	case count == 0:
		return ret
	case count > 0:
		ix, delta := count/n, count%n
		for i := len(ret.data) - 1; i >= 0; i-- {
			tmp := bf64.New()
			if i-ix >= 0 {
				tmp = ret.data[i-ix]
			}
			a, b := tmp.Shift2(delta)
			ret.data[i] = a
			if i+1 < len(bf.data) {
				ret.data[i+1] = ret.data[i+1].Or(b)
			}
		}
		ret.clearEnd()

	case count < 0:
		ix, delta := -count/n, -count%n
		for i := 0; i < len(ret.data); i++ {
			tmp := bf64.New()
			if i+ix < len(ret.data) {
				tmp = ret.data[i+ix]
			}
			a, b := tmp.Shift2(-delta)
			ret.data[i] = a
			if i > 0 {
				ret.data[i-1] = ret.data[i-1].Or(b)
			}
		}
	}
	return ret
}

// Mid returns counts bits from position pos as a new BitField
// Returns nil if count<0
func (bf *BitField) Mid(pos, count int) *BitField {
	switch {
	case count < 0:
		panic("count cannot be negative")

	case count == 0:
		return New(0)

	default:
		if count > bf.Len() {
			count = bf.Len()
		}
		pos = bf.posNormalize(pos)
		return bf.Shift(-pos).Resize(count)
	}
}

// Left returns count bits in the range of [0,count-1] as a new BitField
// Returns nil if count<0
func (bf *BitField) Left(count int) *BitField {
	return bf.Mid(0, count)
}

// Right returns count bits in the range of [63-count,63] as a new BitField
// Returns nil if count<0
func (bf *BitField) Right(count int) *BitField {
	return bf.Mid(bf.Len()-count, count)
}

// Append appends 'other' BitField to the end
// A newly created bitfield will be returned
func (bf *BitField) Append(other *BitField) *BitField {
	if other.Len() == 0 {
		return bf.Clone()
	}
	len := bf.Len()
	newLen := len + other.Len()
	return other.
		Resize(newLen).
		Shift(len).
		Or(bf.Resize(newLen))
}

// Rotate rotates by amount bits and returns it
// If amount>0 it rotates towards higher bit positions,
// otherwise it rotates towards lower bit positions. Mutable.
func (bf *BitField) Rotate(amount int) *BitField {
	ret := bf.mClone()
	if bf.len == 0 {
		return ret
	}
	for amount < 0 {
		amount += bf.len
	}
	amount %= bf.len

	if amount%bf.len == 0 {
		return ret
	}

	lh := bf.Left(bf.Len() - amount)
	rh := bf.Right(amount)
	rh.Append(lh).Copy(ret)
	return ret
}

func (bf *BitField) String() string {
	const n = 64
	ret := make([]byte, bf.len)
	for i := 0; i < len(bf.data); i++ {
		s := bf.data[i].String()
		for j := 0; j < len(s); j++ {
			pos := i*n + j
			if pos >= bf.len {
				break
			}
			ret[pos] = s[j]
		}
	}
	return string(ret)
}
