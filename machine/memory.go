package main

import (
	"errors"
	"fmt"
	"log"
)

const MaxMemorySize = 65535
const mmioIn = 0
const mmioOut = 1

type DataMemory struct {
	ioController IOController
	cells        []int

	dataBus    int
	addressBus int
	memoryOut  int
}

func makeDataMem(cells []DataEntry, inputStream string) (DataMemory, error) {
	initialData := make([]int, MaxMemorySize)
	initialData[mmioIn] = 0
	initialData[mmioOut] = 0

	for _, entry := range cells {
		if entry.Address == mmioIn || entry.Address == mmioOut {
			log.Printf("Data on address %d has collision on mmio\n", entry.Address)
			return DataMemory{}, errors.New("data address has collision on mmio")
		}
		initialData[entry.Address] = int(entry.Data)
	}

	if len(initialData) > MaxMemorySize {
		initialData = initialData[:MaxMemorySize]
	}

	return DataMemory{ioController: makeIOController(inputStream), cells: initialData}, nil
}

func (dataMem *DataMemory) readCell(index int) (int, error) {
	if index == mmioIn {
		//log.Println("Reading from IN buffer")
		return dataMem.ioController.read()
	} else if index == mmioOut {
		log.Println("Cannot read from OUT buffer")
	}
	return dataMem.cells[index], nil
}

func (dataMem *DataMemory) writeCell(index int, value int) error {
	if index == mmioIn {
		//log.Println("Reading from IN buffer")
		return errors.New("cannot write to IN buffer")
	} else if index == mmioOut {
		//log.Println("Writing to OUT buffer")
		dataMem.ioController.write(value)
	} else {
		dataMem.cells[index] = value
	}
	return nil
}

func (dataMem *DataMemory) performReadSignal() error {
	var err error
	dataMem.memoryOut, err = dataMem.readCell(dataMem.addressBus)
	return err
}

func (dataMem *DataMemory) performWriteSignal() error {
	err := dataMem.writeCell(dataMem.addressBus, dataMem.dataBus)
	return err
}

type InstructionMemory struct {
	cells []string
}

func makeInstructionMem(cells []DataEntry) InstructionMemory {
	initialData := make([]string, MaxMemorySize)

	for _, entry := range cells {
		initialData[entry.Address] = fmt.Sprintf("%032b", entry.Data)
	}

	return InstructionMemory{cells: initialData}
}

func (insMem *InstructionMemory) readInstruction(index int) string {
	return insMem.cells[index-1] // читаем c pc, а он всегда на 1 больше адреса
}
