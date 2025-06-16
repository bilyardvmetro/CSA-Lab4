package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

type IOController struct {
	inBuf  []int
	outBuf []int
}

func parsePString(input string) []int {
	var inBuffer []int

	// если одно целое число
	singleNumber, err := strconv.Atoi(input)
	if err == nil {
		inBuffer = make([]int, 2)
		inBuffer[0] = 1 // num is single
		inBuffer[1] = singleNumber
		return inBuffer
	}

	// если список чисел
	parts := strings.Split(input, ",")

	if len(parts) > 1 {
		inBuffer = make([]int, len(parts)+1) // init buf for pstr
		inBuffer[0] = len(parts)             // set str length

		for i, char := range parts {
			trimmedChar := strings.TrimSpace(char)
			num, err := strconv.Atoi(trimmedChar)
			if err != nil {
				panic(err)
			}
			inBuffer[i+1] = num
		}
	} else {
		inBuffer = make([]int, len(input)+1) // init buf for pstr
		inBuffer[0] = len(input)             // set str length

		for i, char := range input {
			inBuffer[i+1] = int(char)
		}
	}

	return inBuffer
}

func makeIOController(input string) IOController {
	if input == "" {
		return IOController{
			inBuf:  nil,
			outBuf: make([]int, 0),
		}
	}

	return IOController{
		inBuf:  parsePString(input),
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
