package main

import (
	"bufio"
	"fmt"
	"os"
)

func openFileScanner(filename string) (*os.File, *bufio.Scanner, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file '%s': %v", filename, err)
	}

	scanner := bufio.NewScanner(file)

	return file, scanner, nil
}

func main() {
	inputFilename := "input.txt"

	file, scanner, err := openFileScanner(inputFilename)
	if err != nil {
		fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		fmt.Println(line)
	}
}