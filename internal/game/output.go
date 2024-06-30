package game

import (
	"fmt"
	escapes "github.com/snugfox/ansi-escapes"
)

type textColor int

const (
	textColorBlack textColor = iota
	textColorRed
	textColorGreen
	textColorYellow
	textColorBlue
	textColorMagenta
	textColorCyan
	textColorWhite
)

func printColor(text string, color textColor) {
	var colorSequence string
	switch color {
	case textColorBlack:
		colorSequence = escapes.TextColorBlack
	case textColorRed:
		colorSequence = escapes.TextColorRed
	case textColorGreen:
		colorSequence = escapes.TextColorGreen
	case textColorYellow:
		colorSequence = escapes.TextColorYellow
	case textColorBlue:
		colorSequence = escapes.TextColorBlue
	case textColorMagenta:
		colorSequence = escapes.TextColorMagenta
	case textColorCyan:
		colorSequence = escapes.TextColorCyan
	case textColorWhite:
		colorSequence = escapes.TextColorWhite
	}

	fmt.Printf("%s%s%s", colorSequence, text, escapes.TextColorWhite)
}

func printColorLn(text string, color textColor) {
	printColor(text+"\n", color)
}

func ln() {
	fmt.Println()
}

func clearScreen() {
	fmt.Print(escapes.ClearScreen)
}
