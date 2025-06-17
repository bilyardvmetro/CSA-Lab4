package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"example.com/CSA-Lab4/isa"
)

// MacroDefinition хранит имя макроса и его тело
type MacroDefinition struct {
	Name string
	Body []string // Каждая строка тела макроса
}

// SectionType определяет тип секции (данные или код)
type SectionType int

const (
	DataSection SectionType = iota
	CodeSection
	UnknownSection
)

// SymbolTables хранит отдельные таблицы символов для данных и кода
type SymbolTables struct {
	DataSymbols map[string]int
	CodeSymbols map[string]int
}

// CodeLine представляет строку кода с ее оригинальным содержанием и назначенным адресом
type CodeLine struct {
	OriginalLine  string      // Исходная строка (после первоначальной очистки)
	Address       int         // Адрес этой строки в памяти команд (адресация по словам)
	IsInstruction bool        // Является ли эта строка инструкцией, которая будет генерировать машинный код
	IsDataDef     bool        // Является ли эта строка определением данных
	DataValue     interface{} // Хранит значение данных (int или string) для второго прохода
	DataSizeWords int         // Размер данных в словах (для определения смещений)
}

type Instruction struct {
	InstructionMemAddress int
	Line                  string
}

// DataEntry представляет собой пару адрес-данные.
type DataEntry struct {
	Address uint32
	Data    uint32
}

// Регулярные выражения для %hi и %lo
var hiRegex = regexp.MustCompile(`%hi\((\w+)\)`)
var loRegex = regexp.MustCompile(`%lo\((\w+)\)`)

// Регулярные выражения для control-flow инструкций
var jalRegex = regexp.MustCompile(`(?i)^(jal)\s+([a-zA-Z0-9_]+),\s*(\w+)$`)
var jalrRegex = regexp.MustCompile(`(?i)^(jalr)\s+([a-zA-Z0-9_]+,\s*[a-zA-Z0-9_]+,\s*)(\w+)$`)
var branchRegex = regexp.MustCompile(`(?i)^(beq|bne|bgt|ble)\s+([a-zA-Z0-9_]+,\s*[a-zA-Z0-9_]+,\s*)(\w+)$`)

var memDumpFile *os.File
var dataMemory *os.File
var instructionMemory *os.File

//var dataMemoryTxt *os.File
//var instructionMemoryTxt *os.File

func makeMemDumpFile(filename string) *os.File {
	dir := filepath.Dir(filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintln(file, "<memory> - <address> - <HEXCODE> - <mnemonic>/<value_dec>>")
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func writeDataToMemDump(file *os.File, address int, val int) {
	_, err := fmt.Fprintf(file, "dataMem    %d:    %X    %d\n", address, val, val)
	if err != nil {
		log.Fatal(err)
	}
}

func writeInstructionToMemDump(file *os.File, address int, binaryIns uint32, mnemonic string) {
	_, err := fmt.Fprintf(file, "progMem    %d:    %08X    %s\n", address, binaryIns, mnemonic)
	if err != nil {
		log.Fatal(err)
	}
}

func makeMemoryFile(filename string) *os.File {
	dir := filepath.Dir(filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

//func writeToMemory(file *os.File, address int, val int) {
//	fmt.Fprintf(file, "%032b:    %032b\n", address, val)
//}

// writeDataEntriesToBinaryFile записывает срез DataEntry в бинарный файл.
// Каждый адрес и данные записываются как 32-битные беззнаковые целые числа.
func writeDataEntryToBinaryFile(file *os.File, entry DataEntry) error {
	byteOrder := binary.LittleEndian

	// Записываем адрес (uint32)
	err := binary.Write(file, byteOrder, entry.Address)
	if err != nil {
		return fmt.Errorf("ошибка записи адреса %d: %v", entry.Address, err)
	}
	// Записываем данные (uint32)
	err = binary.Write(file, byteOrder, entry.Data)
	if err != nil {
		return fmt.Errorf("ошибка записи данных %d: %v", entry.Data, err)
	}

	return nil
}

func makeBinaryRTypeInstruction(tokens []string) string {
	operationEntries := isa.InstructionMap[tokens[0]]

	instructionType := operationEntries[0]
	operationCode := operationEntries[1]
	opExtension := operationEntries[2]

	rd := isa.RegisterMap[tokens[1]]
	rs1 := isa.RegisterMap[tokens[2]]
	rs2 := isa.RegisterMap[tokens[3]]

	return fmt.Sprintf("%s%s%s%s%s%s", instructionType, rd, operationCode, rs1, rs2, opExtension)
}

func makeBinaryITypeInstruction(tokens []string) string {
	rd := "11111"
	rs1 := "11111"
	var imm int64 = 4095

	operationEntries := isa.InstructionMap[tokens[0]]

	instructionType := operationEntries[0]
	operationCode := operationEntries[1]

	if tokens[0] != isa.HALT {
		rd = isa.RegisterMap[tokens[1]]
		rs1 = isa.RegisterMap[tokens[2]]
		var err error

		imm, err = strconv.ParseInt(tokens[3], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		imm = imm & 0xFFF
	}

	return fmt.Sprintf("%s%s%s%s%012b", instructionType, rd, operationCode, rs1, imm)
}

func makeBinarySBTypeInstructions(tokens []string) string {
	operationEntries := isa.InstructionMap[tokens[0]]

	instructionType := operationEntries[0]
	operationCode := operationEntries[1]

	var rs1, rs2 string

	if tokens[0] == isa.SW {
		rs2 = isa.RegisterMap[tokens[1]]
		rs1 = isa.RegisterMap[tokens[2]]
	} else {
		rs1 = isa.RegisterMap[tokens[1]]
		rs2 = isa.RegisterMap[tokens[2]]
	}

	imm, err := strconv.ParseInt(tokens[3], 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	imm = imm & 0xFFF

	immBinaryStr := fmt.Sprintf("%012b", imm)

	return fmt.Sprintf("%s%s%s%s%s%s", instructionType, immBinaryStr[:5], operationCode, rs1, rs2, immBinaryStr[5:])
}

func makeBinaryUJTypeInstruction(tokens []string) string {
	operationEntries := isa.InstructionMap[tokens[0]]

	instructionType := operationEntries[0]
	rd := isa.RegisterMap[tokens[1]]

	imm, err := strconv.ParseInt(tokens[2], 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	var immBinaryStr string
	switch tokens[0] {
	case isa.JAL:
		imm = imm & 0x000FFFFF
		immBinaryStr = fmt.Sprintf("%020b", imm)
	case isa.LUI:
		immBinaryStr = fmt.Sprintf("%020b", imm)
		if len(immBinaryStr) > 20 {
			log.Fatal("Imm is too large")
		}
	}

	return fmt.Sprintf("%s%s%s", instructionType, rd, immBinaryStr)
}

// readLines читает код из файла
func readLines(path string) []string {
	file, _ := os.Open(path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// cleanComments очищает код от комментариев
func cleanComments(lines []string) []string {
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if idx := strings.Index(line, ";"); idx != -1 {
			line = line[:idx]
			line = strings.TrimSpace(line)
		}
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

// expandMacros читает ассемблерный код, находит макросы,
// а затем заменяет их вызовы на соответствующий код.
func expandMacros(inputLines []string) ([]string, error) {
	macros := make(map[string]MacroDefinition)
	var processedLines []string
	var inMacro bool
	var currentMacroName string
	var currentMacroBody []string

	macroStartRegex := regexp.MustCompile(`%macro\s+(\w+)\(([^)]*)\)`) // %macro name(arg1, arg2, ...)
	macroEndRegex := regexp.MustCompile(`%endmacro`)                   // %endmacro
	macroCallRegex := regexp.MustCompile(`^(\w+)\(([^)]*)\)$`)         // name(arg1, arg2, ...)

	// Этап 1: Сбор определений макросов
	for _, line := range inputLines {

		if inMacro {
			if macroEndRegex.MatchString(line) {
				macros[currentMacroName] = MacroDefinition{
					Name: currentMacroName,
					Body: currentMacroBody,
				}
				inMacro = false
				currentMacroName = ""
				currentMacroBody = nil // Обнуляем для следующего макроса
			} else {
				currentMacroBody = append(currentMacroBody, line) // Сохраняем строку как есть
			}
		} else {
			if matches := macroStartRegex.FindStringSubmatch(line); len(matches) > 0 {
				inMacro = true
				currentMacroName = matches[1]
			}
		}
	}

	inMacro = false
	// Этап 2: Замена вызовов макросов
	for _, line := range inputLines {

		// Игнорируем строки, которые являются определением макроса
		if macroStartRegex.MatchString(line) {
			inMacro = true
			continue
		}
		if macroEndRegex.MatchString(line) {
			inMacro = false
			continue
		}

		if matches := macroCallRegex.FindStringSubmatch(line); len(matches) > 0 {
			macroName := matches[1]
			argsString := matches[2]
			var callArgs []string
			if argsString != "" {
				// Разделяем аргументы, учитывая возможные пробелы вокруг запятых
				for _, arg := range strings.Split(argsString, ",") {
					callArgs = append(callArgs, strings.TrimSpace(arg))
				}
			}

			if macroDef, ok := macros[macroName]; ok {
				// Найден вызов макроса, разворачиваем его
				for _, macroLine := range macroDef.Body {
					expandedLine := macroLine
					// Замена аргументов: %1, %2, ...
					for i, arg := range callArgs {
						placeholder := fmt.Sprintf("%%%d", i+1) // %1, %2, ...
						expandedLine = strings.ReplaceAll(expandedLine, placeholder, arg)
					}
					processedLines = append(processedLines, expandedLine)
				}
			} else {
				processedLines = append(processedLines, line)

			}
		} else {
			if !inMacro {
				// Это не определение макроса и не вызов макроса, добавляем строку как есть
				processedLines = append(processedLines, line)
			}
		}
	}

	return processedLines, nil
}

// ProcessAssemblyCode (первый проход): строит таблицы символов и промежуточное представление кода с адресами.
// Возвращает SymbolTables, список CodeLine для второго прохода, и ошибку.
func ProcessAssemblyCode(inputLines []string) (SymbolTables, []CodeLine, error) {
	symbolTables := SymbolTables{
		DataSymbols: make(map[string]int),
		CodeSymbols: make(map[string]int),
	}
	var codeLines []CodeLine // Промежуточное представление кода для второго прохода

	currentAddressData := 0 // Текущий адрес в памяти данных (по словам)
	currentAddressCode := 0 // Текущий адрес в памяти команд (по словам)
	currentSection := UnknownSection
	nextOrgAddress := -1 // Адрес, заданный последней .org

	orgDirectiveRegex := regexp.MustCompile(`^\.org\s+(\d+)$`)
	dataDirectiveRegex := regexp.MustCompile(`^\.data$`)
	codeDirectiveRegex := regexp.MustCompile(`^\.code$`)
	labelDefinitionRegex := regexp.MustCompile(`^(\w+):$`)
	dataValueRegex := regexp.MustCompile(`^(\w+):\s*(\d+|".*")$`) // num | "string"

	for lineNumber, line := range inputLines {
		isInstruction := false
		isDataDef := false

		// Обработка директивы .org
		if matches := orgDirectiveRegex.FindStringSubmatch(line); len(matches) > 0 {
			if currentSection == CodeSection {
				return SymbolTables{}, nil, fmt.Errorf("строка %d: директива .org не может появиться внутри уже открытой секции кода", lineNumber+1)
			}
			addrStr := matches[1]
			addr, err := strconv.Atoi(addrStr)
			if err != nil {
				return SymbolTables{}, nil, fmt.Errorf("строка %d: неверный адрес в .org: %s", lineNumber+1, addrStr)
			}
			nextOrgAddress = addr
			continue
		}

		// Обработка директивы .data
		if dataDirectiveRegex.MatchString(line) {
			if currentSection == CodeSection {
				return SymbolTables{}, nil, fmt.Errorf("строка %d: секция .не может появиться после секции .code", lineNumber+1)
			}
			if currentSection == DataSection {
				return SymbolTables{}, nil, fmt.Errorf("строка %d: секция .уже была объявлена", lineNumber+1)
			}
			currentSection = DataSection
			if nextOrgAddress != -1 {
				currentAddressData = nextOrgAddress
				nextOrgAddress = -1
			}
			continue
		}

		// Обработка директивы .code
		if codeDirectiveRegex.MatchString(line) {
			if currentSection == CodeSection {
				return SymbolTables{}, nil, fmt.Errorf("строка %d: секция .code уже была объявлена", lineNumber+1)
			}
			currentSection = CodeSection
			if nextOrgAddress != -1 {
				currentAddressCode = nextOrgAddress
				nextOrgAddress = -1
			}
			continue
		}

		// Если мы ещё не в какой-либо секции, и встретили что-то, что не является директивой
		if currentSection == UnknownSection {
			if labelDefinitionRegex.MatchString(line) || dataValueRegex.MatchString(line) {
				return SymbolTables{}, nil, fmt.Errorf("строка %d: метка или данные '%s' определены до объявления секции .или .code", lineNumber+1, line)
			}
			if strings.TrimSpace(line) == "" {
				continue
			}
			return SymbolTables{}, nil, fmt.Errorf("строка %d: нераспознанная строка или содержимое вне секции: '%s'", lineNumber+1, line)
		}

		// Обработка определений данных
		if currentSection == DataSection {
			if matches := dataValueRegex.FindStringSubmatch(line); len(matches) > 0 {
				varName := matches[1]
				valuePart := matches[2]

				if _, exists := symbolTables.DataSymbols[varName]; exists {
					return SymbolTables{}, nil, fmt.Errorf("строка %d: переменная '%s' уже определена в памяти данных", lineNumber+1, varName)
				}

				var dataSize int
				var dataVal interface{}

				trimmedPart := strings.TrimSpace(valuePart)
				if strings.HasPrefix(trimmedPart, "\"") && strings.HasSuffix(trimmedPart, "\"") {
					// Это строка
					strVal := strings.Trim(trimmedPart, `"`)
					dataVal = strVal
					dataSize = len(strVal) // Каждый символ - 1 слово
				} else {
					// Это число
					numVal, err := strconv.Atoi(trimmedPart)
					if err != nil {
						return SymbolTables{}, nil, fmt.Errorf("строка %d: неверное значение данных для '%s': '%s'", lineNumber+1, varName, trimmedPart)
					}
					dataVal = numVal
					dataSize = 1 // Одно число - 1 слово
				}

				symbolTables.DataSymbols[varName] = currentAddressData
				isDataDef = true
				codeLines = append(codeLines, CodeLine{
					OriginalLine:  line,
					Address:       currentAddressData,
					IsDataDef:     isDataDef,
					DataValue:     dataVal,
					DataSizeWords: dataSize,
				})
				currentAddressData += dataSize // Увеличиваем адрес на размер данных
				continue
			}
		}

		// Обработка меток (после обработки .org, .data, .code и проверки на UnknownSection)
		if matches := labelDefinitionRegex.FindStringSubmatch(line); len(matches) > 0 {
			labelName := matches[1]
			if currentSection == CodeSection {
				if _, exists := symbolTables.CodeSymbols[labelName]; exists {
					return SymbolTables{}, nil, fmt.Errorf("строка %d: метка '%s' уже определена в памяти команд", lineNumber+1, labelName)
				}
				symbolTables.CodeSymbols[labelName] = currentAddressCode
			}
			if currentSection == DataSection {
				if _, exists := symbolTables.DataSymbols[labelName]; exists {
					return SymbolTables{}, nil, fmt.Errorf("строка %d: метка '%s' уже определена в памяти данных", lineNumber+1, labelName)
				}
				symbolTables.DataSymbols[labelName] = currentAddressData
			}
			codeLines = append(codeLines, CodeLine{OriginalLine: line, Address: -1})
			continue
		}

		// Обработка инструкций (если в секции кода)
		if currentSection == CodeSection {
			if !orgDirectiveRegex.MatchString(line) &&
				!dataDirectiveRegex.MatchString(line) &&
				!codeDirectiveRegex.MatchString(line) &&
				!labelDefinitionRegex.MatchString(line) &&
				!dataValueRegex.MatchString(line) {
				isInstruction = true
				codeLines = append(codeLines, CodeLine{OriginalLine: line, Address: currentAddressCode, IsInstruction: isInstruction})
				currentAddressCode += 1
				continue
			}
		}
		codeLines = append(codeLines, CodeLine{OriginalLine: line, Address: -1})
	}

	return symbolTables, codeLines, nil
}

// ResolveSymbols (второй проход): заменяет ссылки на символы и вычисляет PC-относительные смещения.
func ResolveSymbols(processedCodeLines []CodeLine, symbolTables SymbolTables) ([]Instruction, error) {
	var instructions []Instruction

	for _, codeLine := range processedCodeLines {
		line := codeLine.OriginalLine

		// Если это определение данных, вывести его в специальном формате
		if codeLine.IsDataDef {
			switch v := codeLine.DataValue.(type) {
			case int:
				writeDataToMemDump(memDumpFile, codeLine.Address, v)
				//writeToMemory(dataMemoryTxt, codeLine.Address, v)
				err := writeDataEntryToBinaryFile(dataMemory, DataEntry{uint32(codeLine.Address), uint32(v)})
				if err != nil {
					return nil, err
				}
			case string:
				// Каждый символ строки
				for i, char := range v {
					writeDataToMemDump(memDumpFile, codeLine.Address+i, int(char))
					//writeToMemory(dataMemoryTxt, codeLine.Address+i, int(char))
					err := writeDataEntryToBinaryFile(dataMemory, DataEntry{uint32(codeLine.Address + i), uint32(char)})
					if err != nil {
						return nil, err
					}
				}
			}
			continue
		}

		// Пропускаем строки, которые не генерируют машинный код
		if strings.HasPrefix(line, ".org") ||
			strings.HasPrefix(line, ".data") ||
			strings.HasPrefix(line, ".code") ||
			strings.HasSuffix(line, ":") ||
			strings.TrimSpace(line) == "" {
			continue
		}

		expandedLine := line

		// Заменяем %hi(symbol)
		expandedLine = hiRegex.ReplaceAllStringFunc(expandedLine, func(match string) string {
			symbolName := hiRegex.FindStringSubmatch(match)[1]
			if addr, ok := symbolTables.DataSymbols[symbolName]; ok {
				return fmt.Sprintf("%d", (addr+0x800)>>12) // коррекция для компенсации расширения знака imm
			}
			return match
		})

		// Заменяем %lo(symbol)
		expandedLine = loRegex.ReplaceAllStringFunc(expandedLine, func(match string) string {
			symbolName := loRegex.FindStringSubmatch(match)[1]
			if addr, ok := symbolTables.DataSymbols[symbolName]; ok {
				return fmt.Sprintf("%d", addr&0xFFF)
			}
			return match
		})

		// Обработка control-flow инструкций с PC-относительными смещениями
		if codeLine.IsInstruction {
			if matches := jalRegex.FindStringSubmatch(expandedLine); len(matches) > 0 {
				instrName := matches[1]
				rd := matches[2]
				targetLabel := matches[3]

				if targetAddr, ok := symbolTables.CodeSymbols[targetLabel]; ok {
					offsetWords := targetAddr - (codeLine.Address + 1)

					instructions = append(instructions, Instruction{codeLine.Address, fmt.Sprintf("%s %s, %d", instrName, rd, offsetWords)})
				} else {
					return nil, fmt.Errorf("метка '%s' не найдена для инструкции '%s' по адресу %d", targetLabel, instrName, codeLine.Address)
				}
			} else if matches := branchRegex.FindStringSubmatch(expandedLine); len(matches) > 0 {
				instrName := matches[1]
				registers := matches[2]
				targetLabel := matches[3]

				if targetAddr, ok := symbolTables.CodeSymbols[targetLabel]; ok {
					offsetWords := targetAddr - (codeLine.Address + 1)

					instructions = append(instructions, Instruction{codeLine.Address, fmt.Sprintf("%s %s%d", instrName, registers, offsetWords)})
				} else {
					return nil, fmt.Errorf("метка '%s' не найдена для инструкции '%s' по адресу %d", targetLabel, instrName, codeLine.Address)
				}
			} else if matches := jalrRegex.FindStringSubmatch(expandedLine); len(matches) > 0 {
				instrName := matches[1]
				registers := matches[2]
				immediate := matches[3]

				if numVal, err := strconv.Atoi(immediate); err == nil {
					instructions = append(instructions, Instruction{codeLine.Address, fmt.Sprintf("%s %s%d", instrName, registers, numVal)})
				} else {
					if dataAddr, ok := symbolTables.DataSymbols[immediate]; ok {
						instructions = append(instructions, Instruction{codeLine.Address, fmt.Sprintf("%s %s%d", instrName, registers, dataAddr)})
					} else {
						return nil, fmt.Errorf("недопустимое непосредственное значение/метка '%s' для инструкции jalr по адресу %d", immediate, codeLine.Address)
					}
				}
			} else {
				instructions = append(instructions, Instruction{codeLine.Address, expandedLine})
			}
		}
	}
	return instructions, nil
}

func ConvertProgramToBinary(instructions []Instruction) {
	for _, instruction := range instructions {
		tokens := strings.Split(instruction.Line, " ")
		for i := range tokens {
			tokens[i] = strings.Trim(tokens[i], ",")
		}
		//fmt.Println(tokens)

		var binaryInstruction string

		switch tokens[0] {
		case isa.ADD, isa.SUB, isa.MUL, isa.MULH, isa.DIV, isa.AND, isa.OR, isa.XOR:
			binaryInstruction = makeBinaryRTypeInstruction(tokens)
		case isa.LW, isa.ORI, isa.ADDI, isa.JALR, isa.HALT:
			binaryInstruction = makeBinaryITypeInstruction(tokens)
		case isa.SW, isa.BEQ, isa.BNE, isa.BGT, isa.BLE:
			binaryInstruction = makeBinarySBTypeInstructions(tokens)
		case isa.LUI, isa.JAL:
			binaryInstruction = makeBinaryUJTypeInstruction(tokens)
		default:
			fmt.Printf("Неизвестная инструкция: %s\n", tokens)
		}

		binRepresent, _ := strconv.ParseUint(binaryInstruction, 2, 32)

		//writeToMemory(instructionMemoryTxt, instruction.InstructionMemAddress, int(binRepresent))
		err := writeDataEntryToBinaryFile(instructionMemory, DataEntry{
			uint32(instruction.InstructionMemAddress),
			uint32(binRepresent),
		})
		if err != nil {
			log.Fatal(err)
		}
		writeInstructionToMemDump(memDumpFile, instruction.InstructionMemAddress, uint32(binRepresent), instruction.Line)
	}
}

func writeSymTableToMemDump(tables SymbolTables) {
	_, err := fmt.Fprintln(memDumpFile, "Таблицы символов (адресация по словам):")
	if err != nil {
		log.Fatal(err)
	}

	// для гарантии порядка мапы
	dataKeys := make([]string, 0, len(tables.DataSymbols))
	for key := range tables.DataSymbols {
		dataKeys = append(dataKeys, key)
	}
	sort.Strings(dataKeys)

	codeKeys := make([]string, 0, len(tables.CodeSymbols))
	for key := range tables.CodeSymbols {
		codeKeys = append(codeKeys, key)
	}
	sort.Strings(codeKeys)

	_, err = fmt.Fprintln(memDumpFile, "  Память данных:")
	if err != nil {
		log.Fatal(err)
	}
	for _, dataKey := range dataKeys {
		_, err := fmt.Fprintf(memDumpFile, "    %s: %d\n", dataKey, tables.DataSymbols[dataKey])
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = fmt.Fprintln(memDumpFile, "  Память команд:")
	if err != nil {
		log.Fatal(err)
	}
	for _, codeKey := range codeKeys {
		_, err := fmt.Fprintf(memDumpFile, "    %s: %d\n", codeKey, tables.CodeSymbols[codeKey])
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = fmt.Fprintln(memDumpFile, "--------------------------------------------------------")
	if err != nil {
		log.Fatal(err)
	}
}

func write(lines []string, file string) {
	out, _ := os.Create(file)
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(out)
	w := bufio.NewWriter(out)
	defer func(w *bufio.Writer) {
		err := w.Flush()
		if err != nil {
			log.Fatal(err)
		}
	}(w)

	for _, line := range lines {
		_, err := fmt.Fprintf(w, "%s\n", line)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// os.Args - это слайс строк, содержащий аргументы командной строки.
	// os.Args[0] - это имя самой программы.
	// os.Args[1] - первый аргумент, os.Args[2] - второй и так далее.

	args := os.Args

	if len(args) != 4 {
		fmt.Println("Использование: translator.go <input_file> <target_code_file> <target_data_file>")
	}

	inputFile := args[1]
	targetCodeFile := args[2]
	targetDataFile := args[3]

	filename := filepath.Base(inputFile)
	filenameClean := strings.TrimSuffix(filename, filepath.Ext(inputFile))
	filenameDir := "/" + filenameClean + "/"

	memDumpFile = makeMemDumpFile("../out/" + filenameDir + filenameClean + "_MemoryDump.txt")

	dataMemory = makeMemoryFile("../out/" + filenameDir + filepath.Base(targetDataFile))
	instructionMemory = makeMemoryFile("../out/" + filenameDir + filepath.Base(targetCodeFile))

	//dataMemoryTxt = makeMemoryFile("../out/" + filenameDir + filenameClean + "_data.txt")
	//instructionMemoryTxt = makeMemoryFile("../out/" + filenameDir + filenameClean + "_code.txt")

	lines := readLines(inputFile)
	cleaned := cleanComments(lines)
	expanded, _ := expandMacros(cleaned)
	write(expanded, "../out/"+filenameDir+filenameClean+"_Preprocessed.txt")

	// Первый проход: строим таблицы символов и промежуточное представление кода с адресами
	symbolTables, codeLines, err := ProcessAssemblyCode(expanded)
	if err != nil {
		fmt.Printf("Ошибка в первом проходе: %v\n", err)
		os.Exit(1)
	}

	writeSymTableToMemDump(symbolTables)

	// Второй проход: разрешаем символы и вычисляем смещения
	instructions, err := ResolveSymbols(codeLines, symbolTables)
	if err != nil {
		fmt.Printf("Ошибка во втором проходе: %v\n", err)
		os.Exit(1)
	}

	ConvertProgramToBinary(instructions)
}
