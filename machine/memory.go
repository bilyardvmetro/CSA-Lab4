package main

import (
	"errors"
	"log"
)

const MaxMemorySize = 65535

type DataMemory struct {
	ioController IOController
	cells        []int
	regFile      RegisterFile
	memoryOut    int
}

func (dataMem *DataMemory) init(cells map[int]int, rf RegisterFile, controller IOController) DataMemory {
	initialData := make([]int, MaxMemorySize)
	initialData[0] = 0
	initialData[1] = 0

	for addr, data := range cells {
		if addr > 1 { // without .org code start at address 2
			initialData[addr+2] = data
		} else { // with .org code start at address >= 2
			initialData[addr] = data
		}
	}

	if len(initialData) > MaxMemorySize {
		initialData = initialData[:MaxMemorySize]
	}

	return DataMemory{ioController: controller, cells: initialData, regFile: rf}
}

func (dataMem *DataMemory) readCell(index int) (int, error) {
	if index == 0 {
		log.Println("Reading from IN buffer")
		return dataMem.ioController.read()
	}
	if index == 1 {
		log.Println("Cannot read from OUT buffer")
	}
	return dataMem.cells[index], nil
}

func (dataMem *DataMemory) writeCell(index int, value int) error {
	if index == 0 {
		log.Println("Reading from IN buffer")
		return errors.New("cannot write to IN buffer")
	}
	if index == 1 {
		log.Println("Writing to OUT buffer")
		dataMem.ioController.write(value)
	}
	dataMem.cells[index] = value
	return nil
}

func (dataMem *DataMemory) performReadSignal() error {
	address := dataMem.regFile.leftOut
	var err error
	dataMem.memoryOut, err = dataMem.readCell(address)
	return err
}

func (dataMem *DataMemory) performWriteSignal() error {
	address := dataMem.regFile.leftOut
	value := dataMem.regFile.rightOut
	err := dataMem.writeCell(address, value)
	return err
}

type InstructionMemory struct {
	cells []string
}

func (insMem *InstructionMemory) init(cells map[int]string) InstructionMemory {
	initialData := make([]string, MaxMemorySize)

	for address, instruction := range cells {
		initialData[address] = instruction
	}

	return InstructionMemory{cells: initialData}
}

func (insMem *InstructionMemory) readInstruction(index int) string {
	return insMem.cells[index-1] // читаем c pc, а он всегда на 1 больше адреса
}
