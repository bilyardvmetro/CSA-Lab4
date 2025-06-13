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

	// Latch Register Signals (latch_regn, latch_reg0 to latch_reg31)
	latch_regn  // 5 (General latch register signal, if intended)
	latch_reg31 // 6 (latch reg31)

	// Select MPC Signals
	sel_mpc_inc_one      // 7 (Select mpc to increment by one)
	sel_mpc_inc_two_if_z // 8 (Select mpc to increment by two if Z)
	sel_mpc_inc_two_if_n // 9 (Select mpc to increment by two if N)
	sel_mpc_increment    // 10 (General increment for mpc)
	sel_mpc_operation    // 11 (Select mpc based on operation)
	sel_mpc_zero         // 12 (Select mpc as zero)

	// Register Select Signals
	sel_left_reg    // 13 (Select left reg file register)
	sel_left_reg31  // 14 (Select reg31 as left reg file register)
	sel_right_reg   // 15 (Select right reg file register)
	sel_right_reg31 // 16 (Select reg31 as right reg file register)

	// Data Source Select Signals
	sel_data_src_alu // 17 (Select data source from alu)
	sel_data_src_mem // 18 (Select data source from Memory)
	sel_data_src_cu  // 19 (Select data source from Control Unit)

	// Program Counter Select Signals
	sel_pc_inc // 20 (Select pc to increment)
	sel_pc_alu // 21 (Select pc from alu result)

	// ALU Left/Right Operand Select Signals
	sel_alu_r_inc // 22 (Select alu right operand from incrementer)
	sel_alu_r_rf  // 23 (Select alu right operand from Register File)
	sel_alu_l_pc  // 24 (Select alu left operand from Program Counter)
	sel_alu_l_rf  // 25 (Select alu left operand from Register File)

	// ALU Operation Signals
	alu_add  // 26
	alu_sub  // 27
	alu_mul  // 28
	alu_mulh // 29
	alu_div  // 30
	alu_and  // 31
	alu_or   // 32
	alu_xor  // 33

	// Memory Control Signals
	write_data_mem // 34
	read_data_mem  // 35
)

type MicroProgramMemoryTemplate struct {
	memory [][]Signal
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
