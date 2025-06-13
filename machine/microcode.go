package main

type Signal int

const (
	// Control Signals
	halt Signal = iota // 0

	// Latch Signals
	latch_ir  // 1 (Instruction Register)
	latch_ops // 2 (Operands)
	latch_mpc // 3 (Microprogram Counter)
	latch_pc  // 4 (Program Counter)

	latch_regn
	// Latch Register Signals (latch_reg0 to latch_reg31)
	// нужны для правильного вычисления индекса регистра при работе с register file
	latch_reg0  // 5
	latch_reg1  // 6
	latch_reg2  // 7
	latch_reg3  // 8
	latch_reg4  // 9
	latch_reg5  // 10
	latch_reg6  // 11
	latch_reg7  // 12
	latch_reg8  // 13
	latch_reg9  // 14
	latch_reg10 // 15
	latch_reg11 // 16
	latch_reg12 // 17
	latch_reg13 // 18
	latch_reg14 // 19
	latch_reg15 // 20
	latch_reg16 // 21
	latch_reg17 // 22
	latch_reg18 // 23
	latch_reg19 // 24
	latch_reg20 // 25
	latch_reg21 // 26
	latch_reg22 // 27
	latch_reg23 // 28
	latch_reg24 // 29
	latch_reg25 // 30
	latch_reg26 // 31
	latch_reg27 // 32
	latch_reg28 // 33
	latch_reg29 // 34
	latch_reg30 // 35
	latch_reg31 // 36

	// Select MPC Signals
	sel_mpc_inc_one      // 37 (Select mpc to increment by one)
	sel_mpc_inc_two_if_z // 38
	sel_mpc_inc_two_if_n // 39 (Select mpc to increment by two)
	sel_mpc_increment    // 40 (General increment for mpc)
	sel_mpc_operation    // 41 (Select mpc based on operation)
	sel_mpc_zero         // 42 (Select mpc as zero)

	sel_left_reg
	// Left Register Select Signals (sel_left_reg0 to sel_left_reg31)
	// нужны для правильного вычисления индекса регистра при работе с register file
	sel_left_reg0  // 43
	sel_left_reg1  // 44
	sel_left_reg2  // 45
	sel_left_reg3  // 46
	sel_left_reg4  // 47
	sel_left_reg5  // 48
	sel_left_reg6  // 49
	sel_left_reg7  // 50
	sel_left_reg8  // 51
	sel_left_reg9  // 52
	sel_left_reg10 // 53
	sel_left_reg11 // 54
	sel_left_reg12 // 55
	sel_left_reg13 // 56
	sel_left_reg14 // 57
	sel_left_reg15 // 58
	sel_left_reg16 // 59
	sel_left_reg17 // 60
	sel_left_reg18 // 61
	sel_left_reg19 // 62
	sel_left_reg20 // 63
	sel_left_reg21 // 64
	sel_left_reg22 // 65
	sel_left_reg23 // 66
	sel_left_reg24 // 67
	sel_left_reg25 // 68
	sel_left_reg26 // 69
	sel_left_reg27 // 70
	sel_left_reg28 // 71
	sel_left_reg29 // 72
	sel_left_reg30 // 73
	sel_left_reg31 // 74

	sel_right_reg
	// Right Register Select Signals (sel_right_reg0 to sel_right_reg31)
	// нужны для правильного вычисления индекса регистра при работе с register file
	sel_right_reg0  // 75
	sel_right_reg1  // 76
	sel_right_reg2  // 77
	sel_right_reg3  // 78
	sel_right_reg4  // 79
	sel_right_reg5  // 80
	sel_right_reg6  // 81
	sel_right_reg7  // 82
	sel_right_reg8  // 83
	sel_right_reg9  // 84
	sel_right_reg10 // 85
	sel_right_reg11 // 86
	sel_right_reg12 // 87
	sel_right_reg13 // 88
	sel_right_reg14 // 89
	sel_right_reg15 // 90
	sel_right_reg16 // 91
	sel_right_reg17 // 92
	sel_right_reg18 // 93
	sel_right_reg19 // 94
	sel_right_reg20 // 95
	sel_right_reg21 // 96
	sel_right_reg22 // 97
	sel_right_reg23 // 98
	sel_right_reg24 // 99
	sel_right_reg25 // 100
	sel_right_reg26 // 101
	sel_right_reg27 // 102
	sel_right_reg28 // 103
	sel_right_reg29 // 104
	sel_right_reg30 // 105
	sel_right_reg31 // 106

	// Data Source Select Signals
	sel_data_src_alu // 107 (Select data source from alu)
	sel_data_src_mem // 108 (Select data source from Memory)
	sel_data_src_cu  // 109 (Select data source from Control Unit)

	// Program Counter Select Signals
	sel_pc_inc // 110 (Select pc to increment)
	sel_pc_alu // 111 (Select pc from alu result)

	// ALU Left/Right Operand Select Signals
	sel_alu_r_inc // 112 (Select alu right operand from incrementer)
	sel_alu_r_rf  // 113 (Select alu right operand from Register File)
	sel_alu_l_pc  // 114 (Select alu left operand from Program Counter)
	sel_alu_l_rf  // 115 (Select alu left operand from Register File)

	// ALU Operation Signals
	alu_add  // 116
	alu_sub  // 117
	alu_mul  // 118
	alu_mulh // 119
	alu_div  // 120
	alu_and  // 121
	alu_or   // 122
	alu_xor  // 123

	// Memory Control Signals
	write_data_mem // 124
	read_data_mem  // 125
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
		{sel_mpc_increment, latch_mpc, latch_ir, latch_ops},
		{sel_mpc_operation, latch_mpc},
		// 3 halt
		{halt},
		// 4 add
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_add,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 5 sub
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_sub,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 6 mul
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_mul,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 7 mulh
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_mulh,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 8 div
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_div,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 9 and
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_and,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 10 or
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_or,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 11 xor
		{
			sel_left_reg,
			sel_right_reg,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_xor,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 12 lui
		{sel_data_src_cu, latch_regn, sel_mpc_increment, latch_mpc},
		// 13 addi
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg,
			sel_right_reg31,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_add,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 15 ori
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg,
			sel_right_reg31,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_or,
			sel_data_src_alu,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 17 lw
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg,
			sel_right_reg31,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_add,
			sel_data_src_alu,
			latch_reg31,
			sel_mpc_increment,
			latch_mpc,
		},
		{
			sel_left_reg31,
			read_data_mem,
			sel_data_src_mem,
			latch_regn,
			sel_mpc_zero,
			latch_mpc,
		},
		// 20 sw
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{
			sel_left_reg,
			sel_right_reg31,
			sel_alu_l_rf,
			sel_alu_r_rf,
			alu_add,
			sel_data_src_alu,
			latch_reg31,
			sel_mpc_increment,
			latch_mpc,
		},
		{
			sel_left_reg31,
			sel_right_reg,
			write_data_mem,
			sel_mpc_zero,
			latch_mpc,
		},
		// 23 jal
		{sel_alu_l_pc, sel_alu_r_inc, alu_add, sel_data_src_alu, latch_regn, sel_mpc_increment, latch_mpc},
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{sel_right_reg31, sel_alu_r_rf, sel_alu_l_pc, alu_add, sel_pc_alu, latch_pc, sel_mpc_zero, latch_mpc},
		// 26 jalr
		{sel_alu_l_pc, sel_alu_r_inc, alu_add, sel_data_src_alu, latch_regn, sel_mpc_increment, latch_mpc},
		{sel_data_src_cu, latch_reg31, sel_mpc_increment, latch_mpc},
		{sel_right_reg31, sel_alu_r_rf, sel_left_reg, sel_alu_l_rf, alu_add, sel_pc_alu, latch_pc, sel_mpc_zero, latch_mpc},
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
		// 39 ble
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
		{sel_mpc_inc_two_if_n, sel_mpc_increment, latch_mpc},
		// if not N
		{sel_mpc_zero, latch_mpc},
		// if N
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
		// 44 bgt
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
		{sel_mpc_inc_two_if_n, sel_mpc_increment, latch_mpc},
		// if not N
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
		// if N
		{sel_mpc_zero, latch_mpc},
	},
}
