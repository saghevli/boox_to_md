package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type CaptureMode int

const (
	Heading    CaptureMode = 0
	Time                   = 1
	Note                   = 2
	Annotation             = 3
)

func main() {
	switch len(os.Args) {
	case 1:
		fmt.Print("Please supply an input file.")
	case 2:
	default:
	}
	lines := readLines(os.Args[1])
	notes, headings := extractNotes(lines)
	printMd(notes, headings)
}

func readLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("error opening file")
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		log.Fatal("error reading file")
	}
	return lines
}

// Given a notes file, return a 2d array consisting of heading titles and
// the notes within them.
// Annotations are emitted as notes that are not wrapped in backticks.
func extractNotes(lines []string) (notes [][]string, headings []string) {
	headingIndex := 0
	mode := Heading
	note, annotation := "", ""
	var emptySlice []string
	for i, line := range lines {
		if i < 2 {
			// the first two lines of any notes file are useless
			continue
		}
		switch mode {
		case Heading:
			// if the heading is new, add it to the headings list
			if len(headings) == 0 || headings[len(headings)-1] != line {
				headings = append(headings, line)
				notes = append(notes, emptySlice)
				headingIndex++
			}
			mode++
		case Time:
			mode++
		case Note:
			if strings.HasPrefix(line, "【Original Text】") {
				note = line[19:]
			} else if strings.HasPrefix(line, "【Annotations】") {
				notes[len(headings)-1] = append(notes[len(headings)-1], "`"+note+"`")
				note = ""
				annotation = line[17:]
				mode++
			} else {
				note = note + line
			}
		case Annotation:
			if strings.HasPrefix(line, "-------------------") {
				if len(annotation) != 0 {
					notes[len(headings)-1] = append(notes[len(headings)-1], annotation)
				}
				annotation = ""
				mode = Heading
			} else {
				annotation = annotation + line
			}
		}
	}
	return notes, headings
}

func printMd(notes [][]string, headings []string) {
	if len(notes) != len(headings) {
		log.Fatal("Notes and headings are different lengths.")
	}
	for i := 0; i < len(notes); i++ {
		fmt.Print("## " + headings[i] + "\n")
		for j := 0; j < len(notes[i]); j++ {
			fmt.Print("*  " + notes[i][j] + "\n")
		}
	}
}
