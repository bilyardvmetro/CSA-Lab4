package main

import "strconv"

type DataPath struct {
	dataMem        DataMemory
	instructionMem InstructionMemory
	regFile        RegisterFile
	alu            ALU

	aluRightSelector Signal
	aluLeftSelector  Signal

	pc            int
	pcMuxSelector Signal

	dataSrcMuxSelector Signal

	immFromCU string
}

func makeDataPath(instructions []DataEntry, data []DataEntry, inputStream string) DataPath {
	rf := makeRegFile()
	alu := makeALU()
	insMem := makeInstructionMem(instructions)
	dataMem, err := makeDataMem(data, inputStream)
	if err != nil {
		panic(err)
	}

	return DataPath{
		dataMem:            dataMem,
		instructionMem:     insMem,
		regFile:            rf,
		alu:                alu,
		aluRightSelector:   sel_alu_r_rf,
		aluLeftSelector:    sel_alu_l_rf,
		pc:                 int(instructions[0].Address),
		pcMuxSelector:      sel_pc_inc,
		dataSrcMuxSelector: sel_data_src_alu,
		immFromCU:          "",
	}
}

func extendImmSign(imm string) int {
	bits := len(imm)
	if bits == 32 {
		imm, _ := strconv.Atoi(imm)
		return imm
	}
	signBitMask := int32(1 << (bits - 1))
	rawNumber, _ := strconv.ParseUint(imm, 2, bits)
	number := int32(rawNumber)

	var result int32
	if (number & signBitMask) != 0 {
		result = number | (^((int32(1) << bits) - 1))
	} else {
		result = number
	}
	return int(result)
}

func (d *DataPath) selectPcMux(selector Signal) {
	d.pcMuxSelector = selector
}

func (d *DataPath) selectDataSrcMux(selector Signal) {
	d.dataSrcMuxSelector = selector
}

func (d *DataPath) latchPC() {
	switch d.pcMuxSelector {
	case sel_pc_inc:
		d.pc++
	case sel_pc_alu:
		d.pc = d.alu.aluResult
	}
}

func (d *DataPath) getRegFileInput() int {
	switch d.dataSrcMuxSelector {
	case sel_data_src_alu:
		return d.alu.aluResult
	case sel_data_src_mem:
		return d.dataMem.memoryOut
	case sel_data_src_cu:
		return extendImmSign(d.immFromCU)
	}
	return 0
}

func (d *DataPath) selectAluRight(selector Signal) {
	d.aluRightSelector = selector
	switch d.aluRightSelector {
	case sel_alu_r_rf:
		d.alu.rightIn = d.regFile.rightOut
	case sel_alu_r_inc:
		d.alu.rightIn = 1
	}
}

func (d *DataPath) selectAluLeft(selector Signal) {
	d.aluLeftSelector = selector
	switch d.aluLeftSelector {
	case sel_alu_l_rf:
		d.alu.leftIn = d.regFile.leftOut
	case sel_alu_l_pc:
		d.alu.leftIn = d.pc
	}
}

func (d *DataPath) handleReadSignal() error {
	d.dataMem.addressBus = d.regFile.leftOut
	return d.dataMem.performReadSignal()
}

func (d *DataPath) handleWriteSignal() error {
	d.dataMem.addressBus = d.regFile.leftOut
	d.dataMem.dataBus = d.regFile.rightOut
	return d.dataMem.performWriteSignal()
}
