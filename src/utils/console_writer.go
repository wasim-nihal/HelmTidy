package utils

import (
	"fmt"
	"strings"
)

func PrintTable(data [][]string) {
	// Calculate column widths
	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, cell := range row {
			width := len(cell)
			if width > colWidths[i] {
				colWidths[i] = width
			}
		}
	}

	// Print top border
	printBorder(colWidths)

	// Print table content
	for _, row := range data {
		for i, cell := range row {
			fmt.Printf("| %-*s ", colWidths[i], cell)
		}
		fmt.Println("|")
		printBorder(colWidths)
	}
}

func printBorder(colWidths []int) {
	for _, width := range colWidths {
		fmt.Print("+" + strings.Repeat("-", width+2))
	}
	fmt.Println("+")
}
