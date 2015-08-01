package main

const (
	bitAddrMask = uint64(1<<6) - 1
)

// BitMemory models a bit-addressable memory. WordSize is only
// needed if the convenience functions Word/SetWord are used.
type BitMemory struct {
	Buffer []uint64
	// WordSize must be a divisor of 64
	WordSize uint64
}

func (bm *BitMemory) cellOfBit(addr uint64) *uint64 {
	return &bm.Buffer[addr>>6]
}

// Bit returns true if the bit at addr is set
func (bm *BitMemory) Bit(addr uint64) bool {
	cell := bm.cellOfBit(addr)
	bitMask := uint64(1) << (addr & bitAddrMask)
	return *cell&bitMask != 0
}

// SetBit sets the bit at addr to v
func (bm *BitMemory) SetBit(addr uint64, v bool) {
	cell := bm.cellOfBit(addr)
	bitMask := uint64(1) << (addr & bitAddrMask)
	*cell &= ^bitMask
	if v {
		*cell |= bitMask
	}
}

// Word returns the nth word in memory
func (bm *BitMemory) Word(n uint64) uint64 {
	cell := bm.cellOfBit(n * bm.WordSize)
	mask := (uint64(1) << bm.WordSize) - 1
	wordsPerCell := 64 / bm.WordSize
	nthWordInCell := n % wordsPerCell
	firstBit := nthWordInCell * bm.WordSize
	return (*cell >> firstBit) & mask
}
