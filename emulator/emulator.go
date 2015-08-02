package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/surma-dump/nandcpu/libnandcpu"
)

func main() {
	var (
		wordsize  = flag.Uint64("wordsize", 32, "Word size of CPU")
		numCycles = flag.Uint64("cycles", 1000000, "Number of cycles to emulate")
		dumpStart = flag.Uint64("dumpStart", 0, "Start word of memory dump")
		dumpEnd   = flag.Uint64("dumpEnd", 0, "End word of memory dump")
	)
	flag.Parse()

	if flag.NArg() != 1 {
		flag.PrintDefaults()
		fmt.Println("\nExpected file to read ('-' for stdin)")
		return
	}
	if 64%*wordsize != 0 {
		flag.PrintDefaults()
		fmt.Println("\nWord size must be a multiple of 64")
		return
	}
	in := os.Stdin
	if f := flag.Arg(0); f != "-" {
		var err error
		in, err = os.Open(f)
		if err != nil {
			log.Fatalf("Error openening %s: %s", f, err)
		}
	}

	sbm, err := bitMemoryFromFile(in)
	if err != nil {
		log.Fatalf("Error reading input: %s", err)
	}
	sbm.WordSize = *wordsize
	sbm.AlignMemory()

	cpu := &libnandcpu.NandCPU{BitMemory: sbm}
	for i := uint64(0); i < *numCycles; i++ {
		cpu.Step()
	}
	// wordsize/4 = Number of hex digits per word
	pattern := fmt.Sprintf("%%%02dX\n", *wordsize/4)
	for i := *dumpStart; i < *dumpEnd; i++ {
		fmt.Printf(pattern, sbm.Word(i))
	}
}

func bitMemoryFromFile(rc io.ReadCloser) (*libnandcpu.SimpleBitMemory, error) {
	defer rc.Close()
	sbm := &libnandcpu.SimpleBitMemory{
		Buffer: make([]uint64, 0, 1024),
	}
	var err error
	for err != io.EOF {
		var chunk uint64
		err = binary.Read(rc, binary.LittleEndian, &chunk)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			continue
		}
		sbm.Buffer = append(sbm.Buffer, chunk)
	}
	return sbm, nil
}
