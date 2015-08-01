package main

import "testing"

func TestBitMemory_Bit_0xAA(t *testing.T) {
	mem := &BitMemory{
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
	mem := &BitMemory{
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
	mem := &BitMemory{
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
	mem := &BitMemory{
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
