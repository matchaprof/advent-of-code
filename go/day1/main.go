package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// readInput reads the input file and returns two slices containing the first and second numbers respectively.
func readInput(scanner *bufio.Scanner) ([]int, []int, error) {
    leftList := make([]int, 0, 1000)  // Pre-allocated capacity
    rightList := make([]int, 0, 1000)
    lineNumber := 0

    for scanner.Scan() {
        lineNumber++
        line := strings.TrimSpace(scanner.Text())

        if line == "" {
            continue
        }

        fields := strings.Fields(line)

        if len(fields) != 2 {
            return nil, nil, fmt.Errorf("line %d: expected 2 numbers, got %d fields", lineNumber, len(fields))
        }

        num1, err := strconv.Atoi(fields[0])
        if err != nil {
            return nil, nil, fmt.Errorf("line %d: invalid first number '%s': %v", lineNumber, fields[0], err)
        }

        num2, err := strconv.Atoi(fields[1])
        if err != nil {
            return nil, nil, fmt.Errorf("line %d: invalid second number '%s': %v", lineNumber, fields[1], err)
        }

        leftList = append(leftList, num1)
        rightList = append(rightList, num2)
    }

    if err := scanner.Err(); err != nil {
        return nil, nil, fmt.Errorf("error reading input file: %w", err)
    }

    return leftList, rightList, nil
}

// sortLists sorts both slices in ascending order.
func sortLists(leftList, rightList []int) {
    sort.Ints(leftList)
    sort.Ints(rightList)
}

// printSortedLists prints the sorted lists.
func printSortedLists(leftList, rightList []int) {
    fmt.Printf("First List Sorted Ascending (%d elements):\n%v\n\n", len(leftList), leftList)
    fmt.Printf("Second List Sorted Ascending (%d elements):\n%v\n", len(rightList), rightList)
}

func absInt(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

// Compute new list calculating the element-wise difference between leftList and rightList, aligned by their indices
func calculateDifferenceList(leftList, rightList []int) ([]int, error) {
    if len(leftList) != len(rightList) {
        return nil, fmt.Errorf("slices have different lengths: leftList has %d elements, rightList has %d elements", len(leftList), len(rightList))
    }
    
    differenceList := make([]int, len(leftList))

    for i, left := range leftList {
        differenceList[i] = absInt(rightList[i] - left)
    }

    return differenceList, nil
}

func addAll(differenceList []int) int {
    result := 0
    for _, value := range differenceList {
        result += value
    }

    return result
}

func main() {
    // Open the input file
    file, err := os.Open("input.txt")
    if err != nil {
        log.Fatalf("Failed to open input file: %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Read and parse the input
    leftList, rightList, err := readInput(scanner)
    if err != nil {
        log.Fatalf("Error processing input: %v", err)
    }

    // Verify the number of lines
    totalLines := len(leftList)
    if totalLines != 1000 {
        log.Printf("Warning: Expected 1000 lines, but got %d lines", totalLines)
    }

    // Sort the lists
    sortLists(leftList, rightList)

    // Print the sorted lists
    // printSortedLists(leftList, rightList)

    // Calculate and print the new differenceList
    differenceList, err := calculateDifferenceList(leftList, rightList)
    if err != nil {
        log.Fatalf("Error calculating differences: %v", err)
    }

    // fmt.Println("Difference List:", differenceList)

    total := addAll(differenceList)
    fmt.Println("Total Sum of Differences:", total)
}
