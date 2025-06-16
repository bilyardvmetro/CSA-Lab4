package main

import "fmt"

type Signal int

const (
	// Control Signals
	halt Signal = iota // 0

	// Latch Signals
	latch_ir  // 1 (Index Register)
	latch_rr  // 2 (Registers Register)
	latch_mpc // 3 (Microprogram Counter)
	latch_pc  // 4 (Program Counter)

	latch_regn // 5 (General latch register signal for calculation purposes)
	// Latch Register Signals (latch_reg0 to latch_reg31)
	// нужны для правильного вычисления индекса регистра при работе с register file
	latch_reg0  // 6
	latch_reg1  // 7
	latch_reg2  // 8
	latch_reg3  // 9
	latch_reg4  // 10
	latch_reg5  // 11
	latch_reg6  // 12
	latch_reg7  // 13
	latch_reg8  // 14
	latch_reg9  // 15
	latch_reg10 // 16
	latch_reg11 // 17
	latch_reg12 // 18
	latch_reg13 // 19
	latch_reg14 // 20
	latch_reg15 // 21
	latch_reg16 // 22
	latch_reg17 // 23
	latch_reg18 // 24
	latch_reg19 // 25
	latch_reg20 // 26
	latch_reg21 // 27
	latch_reg22 // 28
	latch_reg23 // 29
	latch_reg24 // 30
	latch_reg25 // 31
	latch_reg26 // 32
	latch_reg27 // 33
	latch_reg28 // 34
	latch_reg29 // 35
	latch_reg30 // 36
	latch_reg31 // 37

	// Select MPC Signals
	sel_mpc_inc_one            // 38 (Select mpc to increment by one)
	sel_mpc_inc_two_if_z       // 39
	sel_mpc_inc_two_if_greater // 40 (Select mpc to increment by two if nz=00 )
	sel_mpc_inc_two_if_lower   // 41 (Select mpc to increment by two if nz=10)
	sel_mpc_increment          // 41 (General increment for mpc)
	sel_mpc_look_up_index      // 42 (Select mpc based on look up table index)
	sel_mpc_zero               // 43 (Select mpc as zero)

	sel_left_reg // 44 (General left register select signal for calculation purposes)
	// Left Register Select Signals (sel_left_reg0 to sel_left_reg31)
	// нужны для правильного вычисления индекса регистра при работе с register file
	sel_left_reg0  // 45
	sel_left_reg1  // 46
	sel_left_reg2  // 47
	sel_left_reg3  // 48
	sel_left_reg4  // 49
	sel_left_reg5  // 50
	sel_left_reg6  // 51
	sel_left_reg7  // 52
	sel_left_reg8  // 53
	sel_left_reg9  // 54
	sel_left_reg10 // 55
	sel_left_reg11 // 56
	sel_left_reg12 // 57
	sel_left_reg13 // 58
	sel_left_reg14 // 59
	sel_left_reg15 // 60
	sel_left_reg16 // 61
	sel_left_reg17 // 62
	sel_left_reg18 // 63
	sel_left_reg19 // 64
	sel_left_reg20 // 65
	sel_left_reg21 // 66
	sel_left_reg22 // 67
	sel_left_reg23 // 68
	sel_left_reg24 // 69
	sel_left_reg25 // 70
	sel_left_reg26 // 71
	sel_left_reg27 // 72
	sel_left_reg28 // 73
	sel_left_reg29 // 74
	sel_left_reg30 // 75
	sel_left_reg31 // 76

	sel_right_reg // 77 (General right register select signal for calculation purposes)
	// Right Register Select Signals (sel_right_reg0 to sel_right_reg31)
	// нужны для правильного вычисления индекса регистра при работе с register file
	sel_right_reg0  // 78
	sel_right_reg1  // 79
	sel_right_reg2  // 80
	sel_right_reg3  // 81
	sel_right_reg4  // 82
	sel_right_reg5  // 83
	sel_right_reg6  // 84
	sel_right_reg7  // 85
	sel_right_reg8  // 86
	sel_right_reg9  // 87
	sel_right_reg10 // 88
	sel_right_reg11 // 89
	sel_right_reg12 // 90
	sel_right_reg13 // 91
	sel_right_reg14 // 92
	sel_right_reg15 // 93
	sel_right_reg16 // 94
	sel_right_reg17 // 95
	sel_right_reg18 // 96
	sel_right_reg19 // 97
	sel_right_reg20 // 98
	sel_right_reg21 // 99
	sel_right_reg22 // 100
	sel_right_reg23 // 101
	sel_right_reg24 // 102
	sel_right_reg25 // 103
	sel_right_reg26 // 104
	sel_right_reg27 // 105
	sel_right_reg28 // 106
	sel_right_reg29 // 107
	sel_right_reg30 // 108
	sel_right_reg31 // 109

	// Data Source Select Signals
	sel_data_src_alu // 110 (Select data source from alu)
	sel_data_src_mem // 111 (Select data source from Memory)
	sel_data_src_cu  // 112 (Select data source from Control Unit)

	// Program Counter Select Signals
	sel_pc_inc // 113 (Select pc to increment)
	sel_pc_alu // 114 (Select pc from alu result)

	// ALU Left/Right Operand Select Signals
	sel_alu_r_inc // 115 (Select alu right operand from incrementer)
	sel_alu_r_rf  // 116 (Select alu right operand from Register File)
	sel_alu_l_pc  // 117 (Select alu left operand from Program Counter)
	sel_alu_l_rf  // 118 (Select alu left operand from Register File)

	// ALU Operation Signals
	alu_add  // 119
	alu_sub  // 120
	alu_mul  // 121
	alu_mulh // 122
	alu_div  // 123
	alu_and  // 124
	alu_or   // 125
	alu_xor  // 126

	// Memory Control Signals
	write_data_mem // 127
	read_data_mem  // 128
)

type MicroProgramMemoryTemplate struct {
	memory [][]Signal
}

func (receiver *MicroProgramMemoryTemplate) getMicroprogramByIndex(index int) []Signal {
	return receiver.memory[index]
}

var MicroProgramMemory = MicroProgramMemoryTemplate{
	memory: [][]Signal{
		//0 Instruction Fetch
		{sel_mpc_inc_one, sel_pc_inc, latch_pc, sel_mpc_increment, latch_mpc},
		{sel_mpc_increment, latch_mpc, latch_ir, latch_rr},
		{sel_mpc_look_up_index, latch_mpc},
		// 3 halt
		{halt},
		// 4 add
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_add,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 5 sub
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_sub,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 6 mul
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_mul,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 7 mulh
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_mulh,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 8 div
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_div,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 9 and
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_and,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 10 or
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_or,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 11 xor
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_xor,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 12 lui
		{sel_data_src_cu, latch_regn, sel_mpc_zero, latch_mpc},
		// 13 addi
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg, sel_right_reg31,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_add,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 15 ori
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg, sel_right_reg31,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_or,
			sel_data_src_alu, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 17 lw
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg, sel_right_reg31,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_add,
			sel_data_src_alu, latch_reg31,
			sel_mpc_increment, latch_mpc,
		},
		{
			sel_left_reg31,
			read_data_mem,
			sel_data_src_mem, latch_regn,
			sel_mpc_zero, latch_mpc,
		},
		// 20 sw
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg, sel_right_reg31,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_add,
			sel_data_src_alu, latch_reg31,
			sel_mpc_increment, latch_mpc,
		},
		{
			sel_left_reg31, sel_right_reg,
			write_data_mem,
			sel_mpc_zero, latch_mpc,
		},
		// 23 jal
		{
			sel_alu_l_pc, sel_alu_r_inc,
			alu_add,
			sel_data_src_alu, latch_regn,
			sel_mpc_increment, latch_mpc,
		},
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_right_reg31, sel_alu_r_rf,
			sel_alu_l_pc,
			alu_add,
			sel_pc_alu, latch_pc,
			sel_mpc_zero, latch_mpc,
		},
		// 26 jalr
		{
			sel_alu_l_pc, sel_alu_r_inc,
			alu_add,
			sel_data_src_alu, latch_regn,
			sel_mpc_increment, latch_mpc,
		},
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_right_reg31, sel_alu_r_rf,
			sel_left_reg, sel_alu_l_rf,
			alu_add,
			sel_pc_alu, latch_pc,
			sel_mpc_zero, latch_mpc,
		},
		// 29 beq
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_sub,
			sel_mpc_increment,
			latch_mpc,
		},
		{sel_mpc_inc_two_if_z, sel_mpc_increment, latch_mpc},
		// if not Z
		{sel_mpc_zero, latch_mpc},
		// if Z
		{
			sel_right_reg31,
			sel_alu_r_rf,
			sel_alu_l_pc,
			alu_add,
			sel_pc_alu,
			latch_pc,
			sel_mpc_zero,
			latch_mpc,
		},
		// 34 bne
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_sub,
			sel_mpc_increment,
			latch_mpc,
		},
		{sel_mpc_inc_two_if_z, sel_mpc_increment, latch_mpc},
		// if not Z
		{
			sel_right_reg31,
			sel_alu_r_rf,
			sel_alu_l_pc,
			alu_add,
			sel_pc_alu,
			latch_pc,
			sel_mpc_zero,
			latch_mpc,
		},
		// if Z
		{sel_mpc_zero, latch_mpc},
		// 39 blt
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_sub,
			sel_mpc_increment, latch_mpc,
		},
		{sel_mpc_inc_two_if_lower, sel_mpc_increment, latch_mpc},
		// else
		{sel_mpc_zero, latch_mpc},
		// if NZ=10
		{
			sel_right_reg31, sel_alu_r_rf,
			sel_alu_l_pc,
			alu_add,
			sel_pc_alu, latch_pc,
			sel_mpc_zero, latch_mpc,
		},
		// 44 bgt
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg, sel_right_reg,
			sel_alu_l_rf, sel_alu_r_rf,
			alu_sub,
			sel_mpc_increment, latch_mpc,
		},
		{sel_mpc_inc_two_if_greater, sel_mpc_increment, latch_mpc},
		// else
		{sel_mpc_zero, latch_mpc},
		// if NZ=00
		{
			sel_right_reg31, sel_alu_r_rf,
			sel_alu_l_pc,
			alu_add,
			sel_pc_alu, latch_pc,
			sel_mpc_zero, latch_mpc,
		},
	},
}

// для отображения имени сигнала в логах
func (s Signal) String() string {
	switch s {
	case halt:
		return "halt"
	case latch_ir:
		return "latch_ir"
	case latch_rr:
		return "latch_rr"
	case latch_mpc:
		return "latch_mpc"
	case latch_pc:
		return "latch_pc"
	case latch_regn:
		return "latch_regn" // Handle the specific latch_regn signal
	case sel_mpc_inc_one:
		return "sel_mpc_inc_one"
	case sel_mpc_inc_two_if_z:
		return "sel_mpc_inc_two_if_z"
	case sel_mpc_inc_two_if_greater:
		return "sel_mpc_inc_two_if_greater"
	case sel_mpc_inc_two_if_lower:
		return "sel_mpc_inc_two_if_lower"
	case sel_mpc_increment:
		return "sel_mpc_increment"
	case sel_mpc_look_up_index:
		return "sel_mpc_look_up_index"
	case sel_mpc_zero:
		return "sel_mpc_zero"
	case sel_left_reg:
		return "sel_left_reg" // Handle the specific sel_left_reg signal
	case sel_right_reg:
		return "sel_right_reg" // Handle the specific sel_right_reg signal
	case sel_data_src_alu:
		return "sel_data_src_alu"
	case sel_data_src_mem:
		return "sel_data_src_mem"
	case sel_data_src_cu:
		return "sel_data_src_cu"
	case sel_pc_inc:
		return "sel_pc_inc"
	case sel_pc_alu:
		return "sel_pc_alu"
	case sel_alu_r_inc:
		return "sel_alu_r_inc"
	case sel_alu_r_rf:
		return "sel_alu_r_rf"
	case sel_alu_l_pc:
		return "sel_alu_l_pc"
	case sel_alu_l_rf:
		return "sel_alu_l_rf"
	case alu_add:
		return "alu_add"
	case alu_sub:
		return "alu_sub"
	case alu_mul:
		return "alu_mul"
	case alu_mulh:
		return "alu_mulh"
	case alu_div:
		return "alu_div"
	case alu_and:
		return "alu_and"
	case alu_or:
		return "alu_or"
	case alu_xor:
		return "alu_xor"
	case write_data_mem:
		return "write_data_mem"
	case read_data_mem:
		return "read_data_mem"
	}

	// Handle ranged register signals dynamically
	if s >= latch_reg0 && s <= latch_reg31 {
		return fmt.Sprintf("latch_reg%d", s-latch_reg0)
	}
	if s >= sel_left_reg0 && s <= sel_left_reg31 {
		return fmt.Sprintf("sel_left_reg%d", s-sel_left_reg0)
	}
	if s >= sel_right_reg0 && s <= sel_right_reg31 {
		return fmt.Sprintf("sel_right_reg%d", s-sel_right_reg0)
	}

	return fmt.Sprintf("UNKNOWN_SIGNAL(%d)", s)
}
