package main

type RegisterFile struct {
	registers   []int
	muxLeftOut  Signal
	muxRightOut Signal
	leftOut     int
	rightOut    int
}

func makeRegFile() RegisterFile {
	return RegisterFile{
		registers:   make([]int, 32),
		muxLeftOut:  sel_left_reg0,
		muxRightOut: sel_right_reg0,
		leftOut:     0,
		rightOut:    0,
	}
}

func (rf *RegisterFile) setLeftRegMux(selector Signal) {
	rf.muxLeftOut = selector
	index := int(selector - sel_left_reg0)
	rf.leftOut = rf.registers[index]
}

func (rf *RegisterFile) setRightRegMux(selector Signal) {
	rf.muxRightOut = selector
	index := int(selector - sel_right_reg0)
	rf.rightOut = rf.registers[index]
}

func (rf *RegisterFile) latchRegN(index int, value int) {
	if index == 0 {
		return
	}
	rf.registers[index] = value
}
