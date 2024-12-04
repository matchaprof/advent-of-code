package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type ValidMul struct {
	X int
	Y int
}

func readInputAndParse(filename string) ([]ValidMul, error) {
	mulPattern := `mul\((-?\d+),(-?\d+)\)`
	re, err := regexp.Compile(mulPattern)
    if err != nil {
        log.Fatalf("Failed to compile regex: %v", err)
    }

	file, err := os.Open(filename)
	if err != nil {
			return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var validMuls []ValidMul

	lineNumber := 0
    for scanner.Scan() {
        lineNumber++
        line := scanner.Text()

        matches := re.FindAllStringSubmatch(line, -1)

        for _, match := range matches {
            // match[0] is the full match, match[1] is x, match[2] is y
            fullMatch := match[0]
            xStr := match[1]
            yStr := match[2]

						x, err1 := strconv.Atoi(xStr)
            y, err2 := strconv.Atoi(yStr)
            if err1 != nil || err2 != nil {
                log.Printf("Failed to parse integers on line %d: x='%s', y='%s'", lineNumber, xStr, yStr)
                continue
            }

						validMul := ValidMul{
							X: x,
							Y: y,
						}
						validMuls = append(validMuls, validMul)

            fmt.Printf("Line %d: Found %s with x=%d and y=%d\n", lineNumber, fullMatch, x, y)
					}
				}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file: %v", err)
	}

	return validMuls, nil
}

func main() {
	inputFilename := "input.txt"

    validMuls, err := readInputAndParse(inputFilename)
    if err != nil {
        log.Fatalf("Error processing input: %v", err)
    }

		fmt.Printf("\nTotal 'mul(x,y)' patterns extracted: %d\n", len(validMuls))

		sum := 0
		for i, mul := range validMuls {
			result := mul.X * mul.Y
			sum += result
			fmt.Printf("MulCall %d: mul(%d, %d) = %d\n", i+1, mul.X, mul.Y, result)
		}

		fmt.Println("Total of all Valid Muls:", sum)
}