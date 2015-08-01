package main

const (
	BitAddrMask = uint64(1<<6) - 1
	// CellMask = ^uint64(0) & ^BitMask
)

type BitMemory struct {
	WordSize byte
	Buffer   []uint64
}

func (bm *BitMemory) cell(addr uint64) *uint64 {
	return &bm.Buffer[addr>>6]
}

func (bm *BitMemory) Bit(addr uint64) bool {
	cell := bm.cell(addr)
	bitMask := uint64(1) << uint64(addr&BitAddrMask)
	return *cell&bitMask != 0
}

func (bm *BitMemory) SetBit(addr uint64, v bool) {
	cell := bm.cell(addr)
	bitMask := uint64(1) << uint64(addr&BitAddrMask)
	*cell &= ^bitMask
	if v {
		*cell |= bitMask
	}
}
