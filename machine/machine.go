package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const tickLimit = 7000

// DataEntry представляет собой пару адрес-данные.
type DataEntry struct {
	Address uint32
	Data    uint32
}

// readDataEntriesFromBinaryFile читает DataEntry из бинарного файла.
func readDataEntriesFromBinaryFile(fileName string) ([]DataEntry, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %v", fileName, err)
	}
	defer file.Close()

	// Используем тот же порядок байтов, что и при записи.
	byteOrder := binary.LittleEndian
	var entries []DataEntry

	for {
		var address uint32
		var data uint32

		// Читаем 32-битный адрес
		err := binary.Read(file, byteOrder, &address)
		if err != nil {
			if err == io.EOF {
				break // Достигнут конец файла
			}
			return nil, fmt.Errorf("ошибка чтения адреса: %v", err)
		}

		// Читаем 32-битные данные
		err = binary.Read(file, byteOrder, &data)
		if err != nil {
			if err == io.EOF {
				// Это может произойти, если файл обрезан посередине пары.
				return nil, fmt.Errorf("неожиданный конец файла после чтения адреса. Файл может быть поврежден")
			}
			return nil, fmt.Errorf("ошибка чтения данных: %v", err)
		}

		entries = append(entries, DataEntry{Address: address, Data: data})
	}

	return entries, nil
}

func readLines(path string) []string {
	file, _ := os.Open(path)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Причина: %v\n", r)
			os.Exit(1) // Выходим с кодом ошибки
		}
	}()

	args := os.Args

	if len(args) < 3 {
		panic("Использование: machine.go <compiled_code> <compiled_data> Optional(<input_file>)")
	}

	instructionsFile := args[1]
	dataFile := args[2]
	var inputStr string

	var inputFile string
	if len(args) > 3 {
		inputFile = args[3]
		input := readLines(inputFile)
		if len(input) != 0 {
			inputStr = input[0]
		}
	}

	programName := filepath.Base(instructionsFile)
	logFile, err := os.OpenFile("../out/"+programName[:len(programName)-9]+"/"+programName[:len(programName)-9]+".log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Если не удалось открыть файл, выводим ошибку в stderr
		log.Fatalf("Ошибка при открытии файла логов: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile)

	instructions := make([]DataEntry, 0)
	data := make([]DataEntry, 0)

	instructionsEntries, _ := readDataEntriesFromBinaryFile(instructionsFile)
	instructions = append(instructions, instructionsEntries...)

	dataEntries, _ := readDataEntriesFromBinaryFile(dataFile)
	data = append(data, dataEntries...)

	//for i, s := range instructions {
	//	fmt.Printf("%d %d: %032b\n", i, s.Address, s.Data)
	//}
	//fmt.Println()
	//for i, s := range data {
	//	fmt.Printf("%d %d: %032b\n", i, s.Address, s.Data)
	//}
	//fmt.Println()
	//fmt.Println()
	//
	//insMem := makeInstructionMem(instructions)
	//dataMem, err := makeDataMem(data, makeRegFile(), makeIOController(""))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println()
	//for i, ins := range insMem.cells {
	//	if ins != "" {
	//		fmt.Println(i, ins)
	//	} else {
	//		continue
	//	}
	//}
	//
	//for i, data := range dataMem.cells {
	//	if i > 20 {
	//		break
	//	}
	//	fmt.Println(i, data)
	//}
	runMmachine(instructions, data, inputStr)
}

func runMmachine(instructions []DataEntry, data []DataEntry, inputStr string) {
	dataPath := makeDataPath(instructions, data, inputStr)
	controlUnit := makeControlUnit(dataPath)

	instructionCounter := 0
	mcCounter := 0

	for controlUnit.ticks < tickLimit {
		if controlUnit.mpc == 0 {
			instructionCounter++
		}

		mcExecuted, err := controlUnit.executeMicroProgram()
		mcCounter += mcExecuted
		log.Printf(
			"Machine state: IR(%d); MPC(%d); PC(%d); NZ(%02b); Ticks(%d); Mc executed(%d)\nRegisters%v",
			controlUnit.ir, controlUnit.mpc, controlUnit.dataPath.pc, controlUnit.dataPath.alu.nz, controlUnit.ticks,
			mcCounter, controlUnit.dataPath.regFile.registers,
		)

		if err != nil {
			fmt.Printf("Stop Reason: %v\n", err)
			fmt.Printf("Instructions executed: %v\n", instructionCounter)
			fmt.Printf("Microprograms executed: %v\n", mcCounter)
			fmt.Printf("Output decimal: %v\n", controlUnit.dataPath.dataMem.ioController.outBuf)
			fmt.Printf("Output hex: ")
			for _, char := range controlUnit.dataPath.dataMem.ioController.outBuf {
				fmt.Printf("%X ", char)
			}
			os.Exit(0)
		}
	}
}
