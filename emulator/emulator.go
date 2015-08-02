package emulator

const (
	bitAddrMask = uint64(1<<6) - 1
)

// BitMemory models a bit-addressable memory with a given word-size.
type BitMemory interface {
	Bit(addr uint64) bool
	SetBit(addr uint64, v bool)
	Word(n uint64) uint64
	NumWords() uint64
}

// SimpleBitMemory implements a BitMemory on a uint64 array. SimpleBitMemory is
// not thread-safe.
type SimpleBitMemory struct {
	Buffer []uint64
	// WordSize must be a divisor of 64
	WordSize uint64
}

func (bm *SimpleBitMemory) cellOfBit(addr uint64) *uint64 {
	return &bm.Buffer[addr>>6]
}

// Bit returns true if the bit at addr is set
func (bm *SimpleBitMemory) Bit(addr uint64) bool {
	if addr >= uint64(len(bm.Buffer)<<6) {
		return false
	}
	cell := bm.cellOfBit(addr)
	bitMask := uint64(1) << (addr & bitAddrMask)
	return *cell&bitMask != 0
}

// SetBit sets the bit at addr to v.
func (bm *SimpleBitMemory) SetBit(addr uint64, v bool) {
	cell := bm.cellOfBit(addr)
	bitMask := uint64(1) << (addr & bitAddrMask)
	*cell &= ^bitMask
	if v {
		*cell |= bitMask
	}
}

// Word returns the nth word in memory.
func (bm *SimpleBitMemory) Word(n uint64) uint64 {
	cell := bm.cellOfBit(n * bm.WordSize)
	mask := (uint64(1) << bm.WordSize) - 1
	wordsPerCell := 64 / bm.WordSize
	nthWordInCell := n % wordsPerCell
	firstBit := nthWordInCell * bm.WordSize
	return (*cell >> firstBit) & mask
}

// NumWords returns the size of memory in number of words.
func (bm *SimpleBitMemory) NumWords() uint64 {
	return uint64(len(bm.Buffer)) * 64 / bm.WordSize
}

// AlignMemory resizes the memory buffer to be a multiple of 3*WordSize using
// append().
func (bm *SimpleBitMemory) AlignMemory() {
	if bm.NumWords()%3 != 0 {
		bm.Buffer = append(bm.Buffer, make([]uint64, 3-(bm.NumWords()%3))...)
	}
}

// NandCPU models a NAND minimal computer.
type NandCPU struct {
	BitMemory
	R1, R2 bool
	PC     uint64
}

// Clock gives the NandCPU a single clock flank.
func (nc *NandCPU) Clock() {
	valAtPC := nc.Word(nc.PC)
	switch nc.PC % 3 {
	case 0:
		nc.R1 = nc.Bit(valAtPC)
	case 1:
		nc.R2 = nc.Bit(valAtPC)
	case 2:
		nc.SetBit(valAtPC, !(nc.R1 && nc.R2))
	}
	nc.PC = (nc.PC + 1) % nc.NumWords()
}
