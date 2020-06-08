package bitfield

import (
	"math"
	"math/bits"
)

// BitField64 is a simple, quick, stack-based bit-field manipulator
// for 64 bits (or less) in length.
// Methods are stateless and free from side-effects.
// It was designed to be chainable. Range for position must be [0, 63]. position
// outside this range will get the modulo treatment, so 64 will point to the 0th
// element, -1 will address the last element (i.e. 63rd), -2 the one before
// (i.e. 62nd), etc.
type BitField64 uint64

// New64 returns a zeroed (all false) bit-field that can store 64 elements
func New64() BitField64 {
	return BitField64(0)
}

// the modulo treatment for positions outside of range [0,63]
func posNormalize(pos int) int {
	const n = 64
	for pos < 0 {
		pos += n
	}
	return pos % n
}

// pos => bit mask: 1 at bit position pos, 0 everywhere else
func posToBitMask(pos int) BitField64 {
	return 1 << uint64(posNormalize(pos))
}

// Set sets the bit at position pos
func (bf64 BitField64) Set(pos int) BitField64 {
	return bf64 | posToBitMask(pos)
}

// SetMul sets the bits at position pos
func (bf64 BitField64) SetMul(pos ...int) BitField64 {
	ret := BitField64(0)
	for _, p := range pos {
		ret |= posToBitMask(p)
	}
	return ret
}

// Get returns true if bit at position pos is set, false otherwise
func (bf64 BitField64) Get(pos int) bool {
	return bf64&posToBitMask(pos) > 0
}

// Clear clears the bit at position pos
func (bf64 BitField64) Clear(pos int) BitField64 {
	return bf64 & ^posToBitMask(pos)
}

// ClearMul sets the bits at position pos
func (bf64 BitField64) ClearMul(pos ...int) BitField64 {
	ret := bf64
	for _, p := range pos {
		ret &= ^posToBitMask(p)
	}
	return ret
}

// Flip inverts the bit at position pos
func (bf64 BitField64) Flip(pos int) BitField64 {
	if bf64.Get(pos) {
		return bf64.Clear(pos)
	}
	return bf64.Set(pos)
}

// SetAll returns a bitfield where all 64 bits are set
func (bf64 BitField64) SetAll() BitField64 {
	return math.MaxUint64
}

// ClearAll returns a bitfield where all 64 bits are set
func (bf64 BitField64) ClearAll() BitField64 {
	return BitField64(0)
}

// And returns the binary AND of the two bitfields
func (bf64 BitField64) And(bfo BitField64) BitField64 {
	return bf64 & bfo
}

// Or returns the binary OR of the two bitfields
func (bf64 BitField64) Or(bfo BitField64) BitField64 {
	return bf64 | bfo
}

// Not returns the bitfield with each bit inverted: 0 becomes 1, 1 becomes 0
func (bf64 BitField64) Not() BitField64 {
	return ^bf64
}

// Xor returns the binary XOR of the two bitfields
func (bf64 BitField64) Xor(bfo BitField64) BitField64 {
	return bf64 ^ bfo
}

// OnesCount returns the number of bits set
func (bf64 BitField64) OnesCount() int {
	return bits.OnesCount64(uint64(bf64))
}

// Mid returns count bits from position pos
func (bf64 BitField64) Mid(pos, count int) BitField64 {
	if count < 0 {
		panic("count is negative")
	}
	const n = 64
	count = count % n
	pos = posNormalize(pos)
	end := pos + count
	a := bf64 << (n - end)
	a = a >> (n - count)
	return a
}

// Left returns leftmost count bits: [0, count-1]
func (bf64 BitField64) Left(count int) BitField64 {
	return bf64.Mid(0, count)
}

// Right returns rightmost count bits [63-count, 63]
func (bf64 BitField64) Right(count int) BitField64 {
	const n = 64
	return bf64.Mid(n-count, count)
}

// Rotate rotates by count bits: Bits exiting at one end entering at the other end.
// If count is positive it rotates towards higher positions; If negative it rotates
// towards lower positions.
func (bf64 BitField64) Rotate(count int) BitField64 {
	if count == 0 {
		return bf64
	}
	return BitField64(bits.RotateLeft64(uint64(bf64), count))
}

// Shift shift bits by count positions. Bits exiting at one end are discarded;
// bits entering at the other end are zeroed.
// If count is positive it shifts towards higher positions; If negative it shifts
// towards lower positions.
func (bf64 BitField64) Shift(count int) BitField64 {
	if count == 0 {
		return bf64
	}
	if count < 0 {
		return bf64 >> -count
	}
	return bf64 << count
}

// Shift2 is same as Shift but it returns the discarded bits as well
func (bf64 BitField64) Shift2(count int) (ret, discarded BitField64) {
	if count == 0 {
		return bf64, New64()
	}
	const n = 64
	if count > n || count < -n {
		return BitField64(0), BitField64(0)
	}
	if count < 0 {
		ret = bf64 >> -count
		discarded = bf64 << (n + count)
	} else {
		ret = bf64 << count
		discarded = bf64 >> (n - count)
	}
	return ret, discarded
}

func (bf64 BitField64) toString(pretty bool) string {
	const n = 64
	lastOne := -1
	b := make([]byte, n)
	for i := 0; i < 64; i++ {
		if bf64&(1<<uint64(i)) > 0 {
			b[i] = '1'
			lastOne = i
		} else {
			b[i] = '0'
		}
	}
	if pretty && lastOne >= 0 {
		return string(b[0 : lastOne+1])
	}
	return string(b)
}

// String returns the bit-representation of the bitfield
// Unlike in number to binary conversion here everything is
// reversed as position grows from 0 to 63, left to right
// So e.g. New().Set(2) print
// 0010000000000000000000000000000000000000000000000000000000000000
func (bf64 BitField64) String() string {
	return bf64.toString(false)
}

// StringPretty converts the bitfield to a string, but tailing zeros
// are dropped. So e.g. New().Set(2) will print "001"
func (bf64 BitField64) StringPretty() string {
	return bf64.toString(true)
}
