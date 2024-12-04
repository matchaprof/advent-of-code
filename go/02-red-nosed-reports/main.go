package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// returns the absolute value of an integer
func absInt(x int) int {
	if x < 0 {
			return -x
	}
	return x
}

// checks if list is sorted (ascending/descending) and meets distance criteria between consecutive numbers
func isSafe(list []int) bool {
	ascending := true
	descending := true

	for i := 1; i < len(list); i++ {
		if list[i] < list[i-1] {
			ascending = false
		}

		if list[i] > list[i-1] {
			descending = false
		}

		distance := absInt(list[i] - list[i-1])
    if distance < 1 || distance > 3 {
			return false
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

			// Append the parsed slice to the master slice
			inputLists = append(inputLists, numbers)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file '%s': %v", filename, err)
	}

	return inputLists, nil
}

// countSafeLists counts how many sorted lists have a maximum distance between
// any two numbers that is not greater than 3.
func countSafeListsStrict(inputLists [][]int) int {
	count := 0

	for _, list := range inputLists {
		if isSafe(list) {
			count++
		}
	}

	return count
}

func countSafeListsFlex(inputLists [][]int) int {
	count := 0

	for _, list := range inputLists {
		if isSafe(list) {
			count++
			continue
		}

		safe :=false
		for i := 0; i < len(list); i++ {
			modifiedList := append([]int{}, list[:i]...)
			modifiedList = append(modifiedList, list[i+1:]...)

			if isSafe(modifiedList) {
				safe = true
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
		strictSafeListCount := countSafeListsStrict(inputLists)
		fmt.Println("Number of Safe Reports (Strict):", strictSafeListCount)
	
		// Solution to Part 2
		flexSafeListCount := countSafeListsFlex(inputLists)
		fmt.Println("Number of Safe Reports via Problem Dampener:", flexSafeListCount)
}