package main

import (
	"bytes"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type GoldenCase struct {
	InSrc   string `yaml:"in_src"`
	OutCode string `yaml:"out_code"`
	OutData string `yaml:"out_data"`
	In      string `yaml:"in"`
	Stdout  string `yaml:"stdout"`
	Log     string `yaml:"log"`
	MemDump string `yaml:"mem_dump"`
}

func TestPrograms(t *testing.T) {
	testDir := filepath.Join("..", "testdata")
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("cannot read testdata dir: %v", err)
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".yml") {
			continue
		}
		program := strings.TrimSuffix(f.Name(), ".yml")
		t.Run(program, func(t *testing.T) {
			runProgramTest(t, program, filepath.Join(testDir, f.Name()))
		})
	}
}

func runProgramTest(t *testing.T, program, ymlPath string) {
	outDir := filepath.Join("..", "out", program)
	_ = os.MkdirAll(outDir, 0755)

	codePath := filepath.Join(outDir, program+"_code.bin")
	dataPath := filepath.Join(outDir, program+"_data.bin")
	//logPath := filepath.Join(outDir, program+".log")
	srcPath := filepath.Join("..", "test_programs", program+".txt")

	var golden GoldenCase
	ymlBytes, err := os.ReadFile(ymlPath)
	if err != nil {
		t.Fatalf("failed to read .yml: %v", err)
	}
	if err := yaml.Unmarshal(ymlBytes, &golden); err != nil {
		t.Fatalf("invalid yaml: %v", err)
	}

	if err := os.WriteFile(srcPath, []byte(golden.InSrc), 0644); err != nil {
		t.Fatalf("failed to write .src: %v", err)
	}

	// === Run translator ===
	translator := filepath.Join("..", "translator", "translator.exe")
	cmdTranslator := exec.Command(translator, srcPath, codePath, dataPath)
	cmdTranslator.Stdout = os.Stdout
	cmdTranslator.Stderr = os.Stderr
	if err := cmdTranslator.Run(); err != nil {
		t.Fatalf("translator failed: %v", err)
	}

	// === Check code.bin and data.bin ===
	assertBinaryEqual(t, codePath, golden.OutCode, "code.bin")
	assertBinaryEqual(t, dataPath, golden.OutData, "data.bin")

	// === Prepare input file for emulator ===
	var inputPath string
	if strings.TrimSpace(golden.In) != "" {
		tmpIn, err := os.CreateTemp("", "stdin-*.txt")
		if err != nil {
			t.Fatalf("failed to create temp input: %v", err)
		}
		defer os.Remove(tmpIn.Name())
		_, _ = tmpIn.WriteString(golden.In)
		tmpIn.Close()
		inputPath = tmpIn.Name()
	}

	// === Run machine ===
	args := []string{codePath, dataPath}
	if inputPath != "" {
		args = append(args, inputPath)
	}

	var stdout bytes.Buffer
	cmdMachine := exec.Command(filepath.Join("..", "machine", "machine.exe"), args...)
	cmdMachine.Stdout = &stdout
	cmdMachine.Stderr = os.Stderr
	if err := cmdMachine.Run(); err != nil {
		t.Fatalf("machine failed: %v", err)
	}

	// === Compare stdout ===
	gotStdout := normalize(stdout.String())
	wantStdout := normalize(golden.Stdout)
	if gotStdout != wantStdout {
		t.Errorf("stdout mismatch\nExpected:\n%s\nGot:\n%s", wantStdout, gotStdout)
	}

	// === Compare log ===
	//logBytes, err := os.ReadFile(logPath)
	//if err != nil {
	//	t.Fatalf("failed to read log: %v", err)
	//}
	//gotLog := normalize(string(logBytes))
	//wantLog := normalize(golden.Log)
	//if gotLog != wantLog {
	//	t.Errorf("log mismatch\nExpected:\n%s\nGot:\n%s", wantLog, gotLog)
	//}

	// === Compare memory dump, if provided ===
	if strings.TrimSpace(golden.MemDump) != "" {
		dumpPath := filepath.Join(outDir, program+"_MemoryDump.txt")
		dumpBytes, err := os.ReadFile(dumpPath)
		if err != nil {
			t.Fatalf("failed to read memory dump: %v", err)
		}
		gotDump := normalize(string(dumpBytes))
		wantDump := normalize(golden.MemDump)
		if gotDump != wantDump {
			t.Errorf("memory dump mismatch\nExpected:\n%s\nGot:\n%s", wantDump, gotDump)
		}
	}
}

func assertBinaryEqual(t *testing.T, path string, hexText string, label string) {
	expected, err := parseHex(hexText)
	if err != nil {
		t.Fatalf("invalid hex for %s: %v", label, err)
	}
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", label, err)
	}
	if !bytes.Equal(expected, actual) {
		t.Errorf("%s mismatch\nExpected: % X\nGot:      % X", label, expected, actual)
	}
}

func parseHex(s string) ([]byte, error) {
	var out []byte
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		b, err := hex.DecodeString(line)
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
	}
	return out, nil
}

func normalize(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\r\n", "\n")
}
