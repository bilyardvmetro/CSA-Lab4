package main

import "math"

const (
	MaxNum int = 1<<31 - 1
	MinNum int = -1 << 31
)

type ALU struct {
	//registerFile RegisterFile
	aluSignal Signal
	rightIn   int
	leftIn    int
	aluResult int
	nz        int
}

func makeALU() ALU {
	return ALU{
		aluResult: 0,
		nz:        0,
	}
}
func (alu *ALU) selectOperation(signal Signal) {
	alu.aluSignal = signal
	alu.performOperation(alu.aluSignal)
}

func (alu *ALU) performOperation(signal Signal) {
	switch signal {
	case alu_add:
		alu.aluResult = alu.leftIn + alu.rightIn
	case alu_sub:
		alu.aluResult = alu.leftIn - alu.rightIn
	case alu_mul:
		alu.aluResult = int(int32(alu.leftIn) * int32(alu.rightIn))
	case alu_mulh:
		alu.aluResult = (alu.leftIn * alu.rightIn) >> 32
	case alu_div:
		alu.aluResult = alu.leftIn / alu.rightIn
	case alu_and:
		alu.aluResult = alu.leftIn & alu.rightIn
	case alu_or:
		alu.aluResult = alu.leftIn | alu.rightIn
	case alu_xor:
		alu.aluResult = alu.leftIn ^ alu.rightIn
	}
	alu.handleOverflow()
	alu.setFlags(alu.aluResult)
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
