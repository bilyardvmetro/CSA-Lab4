package main

import "math"

const (
	MaxNum int = 1<<31 - 1
	MinNum int = -1 << 31
)

type ALU struct {
	registerFile RegisterFile
	aluSignal    Signal
	aluResult    int
	nz           int
}

func (alu *ALU) init(rf RegisterFile) ALU {
	return ALU{registerFile: rf, aluResult: 0, nz: 0}
}
func (alu *ALU) selectOperation(signal Signal) {
	alu.aluSignal = signal
	alu.performOperation(alu.aluSignal)
}

func (alu *ALU) performOperation(signal Signal) {
	switch signal {
	case alu_add:
		alu.aluResult = alu.registerFile.leftOut + alu.registerFile.rightOut
	case alu_sub:
		alu.aluResult = alu.registerFile.leftOut - alu.registerFile.rightOut
	case alu_mul:
		alu.aluResult = alu.registerFile.leftOut * alu.registerFile.rightOut
	case alu_mulh:
		alu.aluResult = (alu.registerFile.leftOut * alu.registerFile.rightOut) >> 32
	case alu_div:
		alu.aluResult = alu.registerFile.leftOut / alu.registerFile.rightOut
	case alu_and:
		alu.aluResult = alu.registerFile.leftOut & alu.registerFile.rightOut
	case alu_or:
		alu.aluResult = alu.registerFile.leftOut | alu.registerFile.rightOut
	case alu_xor:
		alu.aluResult = alu.registerFile.leftOut ^ alu.registerFile.rightOut
	}
	alu.handleOverflow()
}

func (alu *ALU) handleOverflow() {
	if alu.aluResult > MaxNum {
		alu.aluResult %= MaxNum
	} else if alu.aluResult < MaxNum {
		alu.aluResult %= int(math.Abs(float64(MinNum)))
	}
}

func (alu *ALU) setFlags(result int) {
	if result == 0 {
		alu.nz = 0b01
	} else if result < 0 {
		alu.nz = 0b10
	} else {
		alu.nz = 0b00
	}
}
