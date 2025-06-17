package main

import (
	"errors"
	"strconv"

	"example.com/CSA-Lab4/isa"
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
	decoderRegistersOut        [3]int            // rd, rs1, rs2
	lookUpTable                map[[3]string]int // map {type, operation, opExtension} to mpcMem index
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
		dataPath:            dataPath,
		ir:                  0,
		decoderRegistersOut: [3]int{},
		lookUpTable:         instrEntriesMap,
		incValue:            1,
		mpc:                 0,
		ticks:               0,
	}
}

func (c *ControlUnit) readInstruction() string {
	return c.dataPath.instructionMem.readInstruction(c.dataPath.pc)
}

func (c *ControlUnit) decodeInstruction() error {
	instruction := c.readInstruction()

	instrType := instruction[0:typeBits]
	switch instrType {
	case rType:
		err := c.convertRegisterIndexesToInt([3]string{
			instruction[typeBits : typeBits+registerBits],
			instruction[typeBits+registerBits+operationBits : typeBits+registerBits+operationBits+registerBits],
			instruction[typeBits+registerBits+operationBits+registerBits : typeBits+registerBits+operationBits+2*registerBits],
		})
		if err != nil {
			return err
		}

		operation := instruction[typeBits+registerBits : typeBits+registerBits+operationBits]
		opExtension := instruction[typeBits+registerBits+operationBits+2*registerBits : 32]
		instrEntries := [3]string{instrType, operation, opExtension}

		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	case iType:
		err := c.convertRegisterIndexesToInt([3]string{
			instruction[typeBits : typeBits+registerBits],
			instruction[typeBits+registerBits+operationBits : typeBits+registerBits+operationBits+registerBits],
			"",
		})
		if err != nil {
			return err
		}
		c.dataPath.immFromCU = instruction[typeBits+registerBits+operationBits+registerBits : 32]

		operation := instruction[typeBits+registerBits : typeBits+registerBits+operationBits]
		instrEntries := [3]string{instrType, operation, ""}
		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	case sType, bType:
		err := c.convertRegisterIndexesToInt([3]string{
			"",
			instruction[typeBits+immLPartBits+operationBits : typeBits+immLPartBits+operationBits+registerBits],
			instruction[typeBits+immLPartBits+operationBits+registerBits : typeBits+registerBits+operationBits+2*registerBits],
		})
		if err != nil {
			return err
		}
		c.dataPath.immFromCU = instruction[typeBits:typeBits+immLPartBits] + instruction[typeBits+immLPartBits+operationBits+2*registerBits:32]

		operation := instruction[typeBits+immLPartBits : typeBits+immLPartBits+operationBits]
		instrEntries := [3]string{instrType, operation, ""}
		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	case uType:
		err := c.convertRegisterIndexesToInt([3]string{
			instruction[typeBits : typeBits+registerBits], "", "",
		})
		if err != nil {
			return err
		}
		// fill lower 12-bits
		c.dataPath.immFromCU = instruction[typeBits+registerBits:32] + "000000000000"

		instrEntries := [3]string{instrType, "", ""}
		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	case jType:
		err := c.convertRegisterIndexesToInt([3]string{
			instruction[typeBits : typeBits+registerBits], "", "",
		})
		if err != nil {
			return err
		}
		c.dataPath.immFromCU = instruction[typeBits+registerBits : 32]

		instrEntries := [3]string{instrType, "", ""}
		c.decoderLookUpTableIndexOut = c.getMpcValueByInstructionEntries(instrEntries)
	}
	return nil
}

func (c *ControlUnit) getRd() int {
	return c.decoderRegistersOut[0]
}

func (c *ControlUnit) getRs1() int {
	return c.decoderRegistersOut[1]
}

func (c *ControlUnit) getRs2() int {
	return c.decoderRegistersOut[2]
}

func (c *ControlUnit) getMpcValueByInstructionEntries(instrEntries [3]string) int {
	return c.lookUpTable[instrEntries]
}

func (c *ControlUnit) latchIr() error {
	err := c.decodeInstruction()
	if err != nil {
		return err
	}
	c.ir = c.decoderLookUpTableIndexOut
	return nil
}

func (c *ControlUnit) convertRegisterIndexesToInt(indexesStr [3]string) error {
	for i, strRegIndex := range indexesStr {
		if strRegIndex == "" {
			continue
		}
		numRegIndex, err := strconv.ParseUint(strRegIndex, 2, 5)
		if err != nil {
			return err
		}
		c.decoderRegistersOut[i] = int(numRegIndex)
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

func (c *ControlUnit) selRightReg0() {
	c.dataPath.regFile.setRightRegMux(sel_right_reg0)
}

func (c *ControlUnit) selDoubleIncIfGreater() {
	if c.dataPath.alu.nz == 0b00 { // bgt if a>b (only nz==00)
		c.incValue = 2
	} else {
		c.incValue = 1
	}
}
func (c *ControlUnit) selDoubleIncIfLower() {
	if c.dataPath.alu.nz == 0b10 { // blt if a<b (only nz==10)
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
	//log.Printf("Executing signal: %v", signal)

	switch signal {
	case halt:
		return errors.New("HALT")
	case latch_ir:
		err := c.latchIr()
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
	case sel_mpc_inc_two_if_greater:
		c.selDoubleIncIfGreater()
	case sel_mpc_inc_two_if_lower:
		c.selDoubleIncIfLower()
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
	case sel_right_reg0:
		c.selRightReg0()
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
