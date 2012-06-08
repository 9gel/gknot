package gknot

import "fmt"

// Prints a piece laid flat.
func (piece PieceDefinition) Print() {
	const block = '\u2588'
	const esc = '\x1b'
	fmt.Printf("%c[1;%dm%v%c[0m piece:\n", esc, piece.EscColor, piece.Name, esc)
	fmt.Printf("%c[0;%dm", esc, piece.EscColor)
	// Print higher index rows first since the coordinate has y axis going upwards.
	for i := len(piece.Geom)-1; i >= 0; i-- {
		for _, v := range piece.Geom[i] {
			if v == 1 {
				fmt.Printf("%c%c", block, block)
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
	fmt.Printf("%c[0m\n", esc)
}

// Outputs the puzzle in 3 planar projections, along the x-y, y-z and x-z planes.
func (puzzle Puzzle) Print() {
	fmt.Print(puzzle)
	fmt.Println()
}
