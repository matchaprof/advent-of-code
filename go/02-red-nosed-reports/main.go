package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func isSorted(list []int) bool {
	ascending := true
	descending := true

	for i := 1; i < len(list); i++ {
		if list[i] < list[i-1] {
			ascending = false
		}

		if list[i] > list[i-1] {
			descending = false
		}
	}

	return ascending || descending
}

func readAndSortInput(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
			return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var inputLists [][]int

	lineNumber := 0
	for scanner.Scan() {
			lineNumber++
			line := strings.TrimSpace(scanner.Text())

			if line == "" {
					continue
			}

			fields := strings.Fields(line)
			if len(fields) == 0 {
					log.Printf("Line %d is empty after trimming.\n", lineNumber)
					continue
			}

			// Convert each field to an integer
			numbers := make([]int, 0, len(fields))
			for _, field := range fields {
					num, err := strconv.Atoi(field)
					if err != nil {
							log.Printf("Invalid number '%s' on line %d: %v\n", field, lineNumber, err)
							continue
					}

					numbers = append(numbers, num)
			}

			 // Verify if the list is sorted (ascending or descending)
			 if !isSorted(numbers) {
				// log.Printf("Line %d is not sorted. Skipping line.", lineNumber) // Logs the skipped line
				continue
			}

			// Append the parsed slice to the master slice
			inputLists = append(inputLists, numbers)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file '%s': %v", filename, err)
	}

	return inputLists, nil
}

// returns the absolute value of an integer
func absInt(x int) int {
	if x < 0 {
			return -x
	}
	return x
}

// countSafeLists counts how many sorted lists have a maximum distance between
// any two numbers that is not greater than 3.
func countSafeLists(inputLists [][]int) int {
	count := 0

	for _, list := range inputLists {
		safe := true
		for i := 1; i < len(list); i++ {
			distance := absInt(list[i] - list[i-1])
			if distance < 1 || distance > 3 {
				safe = false
				break
			}
		}

		if safe {
			count++
		}
	}

	return count
}

func main() {
    inputFilename := "input.txt"

    inputLists, err := readAndSortInput(inputFilename)
    if err != nil {
        log.Fatalf("Error processing input: %v", err)
    }

    // Output the number of processed lines
    fmt.Printf("Processed %d lines.\n", len(inputLists))

		// Solution to Part 1
		safeListCount := countSafeLists(inputLists)
		fmt.Println("Number of Safe Reports:", safeListCount)
	
}