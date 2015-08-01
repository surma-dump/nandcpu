package main

import "testing"

func TestBitMemory_Bit_0xAA(t *testing.T) {
	mem := &SimpleBitMemory{
		WordSize: 16,
		Buffer:   make([]uint64, 2),
	}
	// 1010101010101010...
	mem.Buffer[0] = uint64(0xAAAAAAAAAAAAAAAA)
	mem.Buffer[1] = uint64(0xAAAAAAAAAAAAAAAA)
	for i := uint64(0); i < 128; i++ {
		if v := mem.Bit(i); v == (i%2 == 0) {
			t.Fatalf("Bit %d is %v", i, v)
		}
	}
}

func TestBitMemory_Bit_0x55(t *testing.T) {
	mem := &SimpleBitMemory{
		WordSize: 16,
		Buffer:   make([]uint64, 2),
	}
	// 0101010101010101...
	mem.Buffer[0] = uint64(0x5555555555555555)
	mem.Buffer[1] = uint64(0x5555555555555555)
	for i := uint64(0); i < 128; i++ {
		if v := mem.Bit(i); v == (i%2 != 0) {
			t.Fatalf("Bit %d is %v", i, v)
		}
	}
}

func TestBitMemory_Bit_0xF0(t *testing.T) {
	mem := &SimpleBitMemory{
		WordSize: 16,
		Buffer:   make([]uint64, 2),
	}
	// 1111000011110000...
	mem.Buffer[0] = uint64(0xF0F0F0F0F0F0F0F0)
	mem.Buffer[1] = uint64(0xF0F0F0F0F0F0F0F0)
	for i := uint64(0); i < 128; i++ {
		if v := mem.Bit(i); v == (i%8 < 4) {
			t.Fatalf("Bit %d is %v", i, v)
		}
	}
}

func TestBitMemory_SetBit(t *testing.T) {
	mem := &SimpleBitMemory{
		WordSize: 16,
		Buffer:   make([]uint64, 2),
	}

	// 1001100110011001...
	mem.Buffer[0] = uint64(0x9999999999999999)
	mem.Buffer[1] = uint64(0x9999999999999999)
	// Unset LSB, Set LSB+1
	for i := uint64(0); i < 32; i++ {
		mem.SetBit(i*4, false)
		mem.SetBit(i*4+1, true)
	}

	// Test for 1010101010101010...
	for i := uint64(0); i < 128; i++ {
		if v := mem.Bit(i); v == (i%2 == 0) {
			t.Fatalf("Bit %d is %v", i, v)
		}
	}
}

func TestBitMemory_Word(t *testing.T) {
	mem := &SimpleBitMemory{
		WordSize: 8,
		Buffer:   make([]uint64, 2),
	}

	mem.Buffer[0] = uint64(0x0706050403020100)
	mem.Buffer[1] = uint64(0x0F0E0D0C0B0A0908)
	for i := uint64(0); i < 16; i++ {
		if v := mem.Word(i); v != i {
			t.Fatalf("Word %d is %x", i, v)
		}
	}
}

func TestNandCPU_8bit(t *testing.T) {
	cpu := &NandCPU{
		BitMemory: &SimpleBitMemory{
			WordSize: 8,
			Buffer: []uint64{
				// Word 0: Load bit 61
				// Word 1: Load bit 62
				// Word 2: Store at bit 63
				0x3F3E3D,
			},
		},
	}
	cpu.SetBit(61, true)
	cpu.SetBit(62, true)
	// Set bit 63 so we can tell it has been successfully unset after 3 clocks
	cpu.SetBit(63, true)
	cpu.Clock()
	cpu.Clock()
	cpu.Clock()
	if cpu.Bit(63) {
		t.Fatalf("Bit 63 is set")
	}
}

func TestNandCPU_4bit(t *testing.T) {
	cpu := &NandCPU{
		BitMemory: &SimpleBitMemory{
			WordSize: 4,
			Buffer: []uint64{
				// Word 0: Load bit 13
				// Word 1: Load bit 14
				// Word 2: Store at bit 15
				0xFED,
			},
		},
	}
	cpu.SetBit(13, true)
	cpu.SetBit(14, true)
	// Set bit 15 so we can tell it has been successfully unset after 3 clocks
	cpu.SetBit(15, true)
	cpu.Clock()
	cpu.Clock()
	cpu.Clock()
	if cpu.Bit(15) {
		t.Fatalf("Bit 15 is set")
	}
}

func TestNandCPU_16bit(t *testing.T) {
	cpu := &NandCPU{
		BitMemory: &SimpleBitMemory{
			WordSize: 16,
			Buffer: []uint64{
				// Word 0: Load bit 61
				// Word 1: Load bit 62
				// Word 2: Store at bit 63
				0x003F003E003D,
			},
		},
	}
	cpu.SetBit(61, true)
	cpu.SetBit(62, true)
	// Set bit 63 so we can tell it has been successfully unset after 3 clocks
	cpu.SetBit(63, true)
	cpu.Clock()
	cpu.Clock()
	cpu.Clock()
	if cpu.Bit(63) {
		t.Fatalf("Bit 63 is set")
	}
}

func TestNandCPU_wraparound(t *testing.T) {
	cpu := &NandCPU{
		BitMemory: &SimpleBitMemory{
			WordSize: 8,
			Buffer: []uint64{
				// Word 0: Load bit 29
				// Word 1: Load bit 30
				// Word 2: Store at bit 31
				// ... 0 ...
				// Word 5: Load bit 27
				// Word 6: Load bit 28
				// Word 7: Store at bit 29
				0x1B1C1D00001F1E1D,
			},
		},
		// Start at word 5
		PC: 5,
	}
	cpu.SetBit(27, true)
	cpu.SetBit(28, true)
	// Bit 29 will be set by executing word 5-7
	cpu.SetBit(30, true)
	// Set bit 31 so we can tell it has been successfully unset after 6 clocks
	cpu.SetBit(31, true)
	cpu.Clock()
	cpu.Clock()
	cpu.Clock()
	cpu.Clock()
	cpu.Clock()
	cpu.Clock()
	if cpu.Bit(31) {
		t.Fatalf("Bit 31 is set")
	}
}
