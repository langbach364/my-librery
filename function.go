package main

import (
	"fmt"
	"os"
	"os/exec"
)

func write_file(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo file: %v", err)
	}
	return file, nil
}

func read_file(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi mở file: %v", err)
	}
	return file, nil
}

func run_script(scriptName string) {
	cmd := exec.Command("bash", scriptName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Lỗi khi chạy %s: %v\n", scriptName, err)
	}
}
