package main

import (
	"errors"
	"example.com/CSA-Lab4/isa"
	"log"
	"strconv"
)

const (
	typeBits      = 7
	operationBits = 3
	//opExtensionBits = 7
	registerBits = 5
	immLPartBits = 5
	//immRPartBits    = 7
	rType = "1000000"
	iType = "0100000"
	sType = "0010000"
	bType = "0001000"
	uType = "0000100"
	jType = "0000010"
)

type ControlUnit struct {
	dataPath DataPath

	ir                         int               // contains mpcMem index mapped on instruction
	rr                         [3]int            // rd, rs1, rs2
	lookUpTable                map[[3]string]int // map {type, operation, opExtension} to mpcMem index
	decoderRegistersOut        [3]string
	decoderLookUpTableIndexOut int

	mpcMux   Signal
	incValue int

	mpc   int
	ticks int
}

func makeControlUnit(dataPath DataPath) ControlUnit {
	instrEntriesMap := map[[3]string]int{
		isa.InstructionMap[isa.HALT]: 3,
		isa.InstructionMap[isa.ADD]:  4,
		isa.InstructionMap[isa.SUB]:  5,
		isa.InstructionMap[isa.MUL]:  6,
		isa.InstructionMap[isa.MULH]: 7,
		isa.InstructionMap[isa.DIV]:  8,
		isa.InstructionMap[isa.AND]:  9,
		isa.InstructionMap[isa.OR]:   10,
		isa.InstructionMap[isa.XOR]:  11,
		isa.InstructionMap[isa.LUI]:  12,
		isa.InstructionMap[isa.ADDI]: 13,
		isa.InstructionMap[isa.ORI]:  15,
		isa.InstructionMap[isa.LW]:   17,
		isa.InstructionMap[isa.SW]:   20,
		isa.InstructionMap[isa.JAL]:  23,
		isa.InstructionMap[isa.JALR]: 26,
		isa.InstructionMap[isa.BEQ]:  29,
		isa.InstructionMap[isa.BNE]:  34,
		isa.InstructionMap[isa.BLE]:  39,
		isa.InstructionMap[isa.BGT]:  44,
	}

	return ControlUnit{
		dataPath:    dataPath,
		ir:          0,
		rr:          [3]int{},
		lookUpTable: instrEntriesMap,
		incValue:    1,
		mpc:         0,
		ticks:       0,
	}
}

func (c *ControlUnit) readInstruction() string {
	return c.dataPath.instructionMem.readInstruction(c.dataPath.pc)
}

func (c *ControlUnit) decodeInstruction() {
	instruction := c.readInstruction()

	instrType := instruction[0:typeBits]
	switch instrType {
	case rType:
		c.decoderRegistersOut = [3]string{
			instruction[typeBits : typeBits+registerBits],
			instruction[typeBits+registerBits+operationBits : typeBits+registerBits+operationBits+registerBits],
			instruction[typeBits+registerBits+operationBits+registerBits : typeBits+registerBits+operationBits+2*registerBits],
		}

		operation := instruction[typeBits+registerBits : typeBits+registerBits+operationBits]
		opExtension := instruction[typeBits+registerBits+operationBits+2*registerBits : 32]
		instrEntries := [3]string{instrType, operation, opExtension}

		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	case iType:
		c.decoderRegistersOut = [3]string{
			instruction[typeBits : typeBits+registerBits],
			instruction[typeBits+registerBits+operationBits : typeBits+registerBits+operationBits+registerBits],
			"",
		}
		c.dataPath.immFromCU = instruction[typeBits+registerBits+operationBits+registerBits : 32]

		operation := instruction[typeBits+registerBits : typeBits+registerBits+operationBits]
		instrEntries := [3]string{instrType, operation, ""}

		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	case sType, bType:
		c.decoderRegistersOut = [3]string{
			"",
			instruction[typeBits+immLPartBits+operationBits : typeBits+immLPartBits+operationBits+registerBits],
			instruction[typeBits+immLPartBits+operationBits+registerBits : typeBits+registerBits+operationBits+2*registerBits],
		}
		c.dataPath.immFromCU = instruction[typeBits:typeBits+immLPartBits] + instruction[typeBits+immLPartBits+operationBits+2*registerBits:32]

		operation := instruction[typeBits+immLPartBits : typeBits+immLPartBits+operationBits]
		instrEntries := [3]string{instrType, operation, ""}

		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	case uType:
		c.decoderRegistersOut = [3]string{
			instruction[typeBits : typeBits+registerBits], "", "",
		}
		// fill lower 12-bits
		c.dataPath.immFromCU = instruction[typeBits+registerBits:32] + "000000000000"

		instrEntries := [3]string{instrType, "", ""}

		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	case jType:
		c.decoderRegistersOut = [3]string{
			instruction[typeBits : typeBits+registerBits], "", "",
		}
		c.dataPath.immFromCU = instruction[typeBits+registerBits : 32]

		instrEntries := [3]string{instrType, "", ""}

		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	}
}

func (c *ControlUnit) getRd() int {
	return c.rr[0]
}

func (c *ControlUnit) getRs1() int {
	return c.rr[1]
}

func (c *ControlUnit) getRs2() int {
	return c.rr[2]
}

func (c *ControlUnit) getMpcValueByInstructionEntries(instrEntries [3]string) int {
	return c.lookUpTable[instrEntries]
}

func (c *ControlUnit) latchIr() {
	c.decodeInstruction()
	c.ir = c.decoderLookUpTableIndexOut
}

func (c *ControlUnit) latchRr() error {
	for i, strRegIndex := range c.decoderRegistersOut {
		if strRegIndex == "" {
			continue
		}
		numRegIndex, err := strconv.ParseUint(strRegIndex, 2, 5)
		if err != nil {
			return err
		}
		c.rr[i] = int(numRegIndex)
	}
	return nil
}

func (c *ControlUnit) latchMpc() {
	switch c.mpcMux {
	case sel_mpc_increment:
		c.mpc += c.incValue
	case sel_mpc_look_up_index:
		c.mpc = c.ir
	case sel_mpc_zero:
		c.mpc = 0
	}
}

func (c *ControlUnit) latchReg31() {
	c.dataPath.regFile.latchRegN(int(latch_reg31-latch_reg0), c.dataPath.getRegFileInput())
}

func (c *ControlUnit) latchRegN() {
	c.dataPath.regFile.latchRegN(c.getRd(), c.dataPath.getRegFileInput())
}

// rs1 всегда лево, rs2 всегда право
func (c *ControlUnit) selLeftReg() {
	c.dataPath.regFile.setLeftRegMux(Signal(int(sel_left_reg0) + c.getRs1()))
}

func (c *ControlUnit) selRightReg() {
	c.dataPath.regFile.setRightRegMux(Signal(int(sel_right_reg0) + c.getRs2()))
}

func (c *ControlUnit) selLeftReg31() {
	c.dataPath.regFile.setLeftRegMux(sel_left_reg31)
}

func (c *ControlUnit) selRightReg31() {
	c.dataPath.regFile.setRightRegMux(sel_right_reg31)
}

func (c *ControlUnit) selDoubleIncIfN() {
	if c.dataPath.alu.nz == 0b10 {
		c.incValue = 2
	} else {
		c.incValue = 1
	}
}

func (c *ControlUnit) selDoubleIncIfZ() {
	if c.dataPath.alu.nz == 0b01 {
		c.incValue = 2
	} else {
		c.incValue = 1
	}
}

func (c *ControlUnit) dispatchSignal(signal Signal) error {
	log.Printf("Executing signal: %v", signal)

	switch signal {
	case halt:
		return errors.New("HALT")
	case latch_ir:
		c.latchIr()
	case latch_rr:
		err := c.latchRr()
		if err != nil {
			return err
		}
	case latch_mpc:
		c.latchMpc()
	case latch_pc:
		c.dataPath.latchPC()
	case latch_regn:
		c.latchRegN()
	case latch_reg31:
		c.latchReg31()
	case sel_mpc_inc_one:
		c.incValue = 1
	case sel_mpc_inc_two_if_n:
		c.selDoubleIncIfN()
	case sel_mpc_inc_two_if_z:
		c.selDoubleIncIfZ()
	case sel_mpc_increment, sel_mpc_look_up_index, sel_mpc_zero:
		c.mpcMux = signal
	case sel_left_reg:
		c.selLeftReg()
	case sel_left_reg31:
		c.selLeftReg31()
	case sel_right_reg:
		c.selRightReg()
	case sel_right_reg31:
		c.selRightReg31()
	case sel_data_src_alu, sel_data_src_mem, sel_data_src_cu:
		c.dataPath.selectDataSrcMux(signal)
	case sel_pc_inc, sel_pc_alu:
		c.dataPath.selectPcMux(signal)
	case sel_alu_r_rf, sel_alu_r_inc:
		c.dataPath.selectAluRight(signal)
	case sel_alu_l_rf, sel_alu_l_pc:
		c.dataPath.selectAluLeft(signal)
	case alu_add, alu_sub, alu_mul, alu_mulh, alu_div, alu_and, alu_or, alu_xor:
		c.dataPath.alu.selectOperation(signal)
	case read_data_mem:
		err := c.dataPath.handleReadSignal()
		if err != nil {
			return err
		}
	case write_data_mem:
		err := c.dataPath.handleWriteSignal()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ControlUnit) executeMicroProgram() (int, error) {
	mp := MicroProgramMemory.getMicroprogramByIndex(c.mpc)
	for _, signal := range mp {
		err := c.dispatchSignal(signal)
		if err != nil {
			return len(mp), err
		}
	}
	c.tick()
	return len(mp), nil
}

func (c *ControlUnit) tick() {
	c.ticks++
}
