package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func openFileScanner(filename string) (*os.File, *bufio.Scanner, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file '%s': %v", filename, err)
	}

	scanner := bufio.NewScanner(file)

	return file, scanner, nil
}

func getMulPatternRegex() (*regexp.Regexp, error) {
	mulPattern := `mul\((-?\d+),(-?\d+)\)`
	re, err := regexp.Compile(mulPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %v", err)
	}

	return re, nil
}

func getCombinedPatternRegex() (*regexp.Regexp, error) {
	combinedPattern := `do\(\)|don't\(\)|mul\((-?\d+),(-?\d+)\)`
	combinedRe, err := regexp.Compile(combinedPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile combined regex: %v", err)
	}

	return combinedRe, nil
}

type ValidMul struct {
	X int
	Y int
}

type MulMatch struct {
	ValidMul
	Start int
}

func extractMulPatterns(line string, re *regexp.Regexp) []MulMatch {
	matches := re.FindAllStringSubmatchIndex(line, -1)
	var muls []MulMatch

	for _, match := range matches {
		if len(match) < 5 {
			continue
		}
		fullStart := match[0]
		xStart := match[2]
		xEnd := match[3]
		yStart := match[4]
		yEnd := match[5]

		xStr := line[xStart:xEnd]
		yStr := line[yStart:yEnd]

		x, err1 := strconv.Atoi(xStr)
		y, err2 := strconv.Atoi(yStr)
		if err1 != nil || err2 != nil {
			log.Printf("Failed to parse integers: x='%s', y='%s'", xStr, yStr)
			continue
		}

		validMul := ValidMul{
			X: x,
			Y: y,
		}

		mulMatch := MulMatch{
			ValidMul: validMul,
			Start:    fullStart,
		}

		muls = append(muls, mulMatch)
	}

	return muls
}

// Parses all mul(x,y) patterns from input file (Part 1)
func parseAllValidMuls(filename string, re *regexp.Regexp) ([]MulMatch, error) {
	file, scanner, err := openFileScanner(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	var mulMatches []MulMatch

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		matches := extractMulPatterns(line, re)

		for _, match := range matches {
			mulMatches = append(mulMatches, match)
			fmt.Printf("Line %d: Found mul(%d,%d)\n", lineNumber, match.X, match.Y)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file: %v", err)
	}

	return mulMatches, nil
}

// TokenType represents the type of token identified in the line.
type TokenType int

const (
	TokenDo TokenType = iota
	TokenDont
	TokenMul
	TokenUnknown
)

// Classifies the token and returns TokenType
func getTokenType(token string) TokenType {
	switch token {
	case "do()":
		return TokenDo
	case "don't()":
		return TokenDont
	default:
		if strings.HasPrefix(token, "mul(") {
			return TokenMul
		}

		return TokenUnknown
	}
}

// Updates the state (enabled/disabled) based on TokenType
func updateState(tokenType TokenType, enabled *bool, lineNumber int) {
	switch tokenType {
	case TokenDo:
		*enabled = true
		fmt.Printf("Line %d: Encountered do(), enabling mul(x,y)\n", lineNumber)
	case TokenDont:
		*enabled = false
		fmt.Printf("Line %d: Encountered don't(), disabling mul(x,y)\n", lineNumber)
	}
}

// extractSingleMul parses the mul(x,y) token and returns a ValidMul struct.
func extractSingleMul(token string, re *regexp.Regexp) (ValidMul, error) {
	matches := re.FindStringSubmatch(token)
	if len(matches) < 3 {
		return ValidMul{}, fmt.Errorf("invalid mul format: %s", token)
	}
	x, err1 := strconv.Atoi(matches[1])
	y, err2 := strconv.Atoi(matches[2])
	if err1 != nil || err2 != nil {
		return ValidMul{}, fmt.Errorf("invalid integers in mul: %s", token)
	}
	return ValidMul{X: x, Y: y}, nil
}

// Processes a single token and updates the state or collects mul(x,y) patterns.
func processToken(token string, mulRe *regexp.Regexp, enabled *bool, lineNumber int, conditionalMuls *[]MulMatch, start int) {
	tokenType := getTokenType(token)

	switch tokenType {
	case TokenDo, TokenDont:
		updateState(tokenType, enabled, lineNumber)
	case TokenMul:
		mul, err := extractSingleMul(token, mulRe)
		if err != nil {
			log.Printf("Line %d: %v", lineNumber, err)
			return
		}

		mulMatch := MulMatch{
			ValidMul: mul,
			Start:    start,
		}

		if *enabled {
			*conditionalMuls = append(*conditionalMuls, mulMatch)
			fmt.Printf("Line %d: Found mul(%d,%d) while enabled\n", lineNumber, mul.X, mul.Y)
		} else {
			fmt.Printf("Line %d: Found mul(%d,%d) but currently disabled\n", lineNumber, mul.X, mul.Y)
		}
	default:
		// Ignore TokenUnknown type
	}
}

// Processes all combined matches within a line.
func processCombinedMatches(line string, combinedRe *regexp.Regexp, mulRe *regexp.Regexp, lineNumber int, enabled *bool, conditionalMuls *[]MulMatch) {
	combinedMatches := combinedRe.FindAllStringSubmatchIndex(line, -1)

	for _, match := range combinedMatches {
		if len(match) < 2 {
			continue
		}
		start := match[0]
		end := match[1]
		token := line[start:end]

		processToken(token, mulRe, enabled, lineNumber, conditionalMuls, start)
	}
}

// Extract mul(x,y) patterns based on Enabled state (part 2)
func parseMulsWithEnabledState(filename string, mulRe *regexp.Regexp, combinedRe *regexp.Regexp) ([]MulMatch, error) {
	file, scanner, err := openFileScanner(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	var conditionalMuls []MulMatch
	enabled := true

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		processCombinedMatches(
			line,
			combinedRe,
			mulRe,
			lineNumber,
			&enabled,
			&conditionalMuls,
		)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file: %v", err)
	}

	return conditionalMuls, nil
}

//	Calculates the sum of all x * y products from the mulMatches input.
// Prints each multiplication and result
func sumUpMuls(mulMatches []MulMatch) int {
	sum := 0
	for i, mul := range mulMatches {
		result := mul.X * mul.Y
		sum += result
		fmt.Printf("MulCall %d: mul(%d, %d) = %d\n", i+1, mul.X, mul.Y, result)
	}

	return sum
}

func main() {
	inputFilename := "input.txt"

	mulPatternPart1, err := getMulPatternRegex()
	if err != nil {
		log.Fatalf("Error compiling Mul Pattern regex: %v", err)
	}

	mulPatternPart2, err := getCombinedPatternRegex()
	if err != nil {
		log.Fatalf("Error compiling Enabled Mul Pattern regex: %v", err)
	}

	// Part 1 Logic
	validMulsPart1, err := parseAllValidMuls(inputFilename, mulPatternPart1)
	if err != nil {
		log.Fatalf("Error processing input for Part 1: %v", err)
	}

	part1Sum := sumUpMuls(validMulsPart1)

	// Part 2 Logic
	validMulsPart2, err := parseMulsWithEnabledState(inputFilename, mulPatternPart1, mulPatternPart2)
	if err != nil {
		log.Fatalf("Error processing input for Part 2: %v", err)
	}

	part2Sum := sumUpMuls(validMulsPart2)

	// Solution to Part 1
	fmt.Printf("\nTotal 'mul(x,y)' patterns extracted for Part 1: %d\n", len(validMulsPart1))
	fmt.Println("Total of all Valid Muls:", part1Sum)

	// Solution to Part 2
	fmt.Printf("\nTotal 'mul(x,y)' patterns extracted for Part 2: %d\n", len(validMulsPart2))
	fmt.Println("Total of all Valid Muls in Enabled do() State:", part2Sum)
}