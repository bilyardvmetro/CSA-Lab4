package main

import (
	"errors"
	"log"
)

type IOController struct {
	inBuf  []int
	outBuf []int
}

func (controller *IOController) init(input string) IOController {
	if input == "" {
		return IOController{
			inBuf:  nil,
			outBuf: make([]int, 0),
		}
	}

	inBuffer := make([]int, len(input)+1) // init buf for pstr
	inBuffer[0] = len(input)              // set str length

	for i, char := range input {
		inBuffer[i+1] = int(char)
	}

	return IOController{
		inBuf:  inBuffer,
		outBuf: make([]int, 0),
	}
}

func (controller *IOController) write(value int) {
	controller.outBuf = append(controller.outBuf, value)
	log.Printf("В out буфер записан символ:  %d (%c)\n", value, value)
}

func (controller *IOController) read() (int, error) {
	if controller.inBuf == nil || len(controller.inBuf) == 0 {
		return 0, errors.New("попытка прочитать данные из пустого буфера")
	}
	value := controller.inBuf[0]
	controller.inBuf = controller.inBuf[1:]
	log.Printf("Из in буфера прочитан символ: %d (%c)\n", value, value)
	return value, nil
}
