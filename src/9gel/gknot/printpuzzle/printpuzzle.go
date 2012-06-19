// Author: nigelchoi@google.com (Nigel Choi)

// Prints the puzzle.

package main

import (
	"9gel/gknot"
	"fmt"
)

func main() {
	puzzle := gknot.NewPuzzle()
	puzzle.Print()

	fmt.Println("Moving the Orange piece by 1 along x axis:")
	mutation := gknot.Mutation{35, gknot.TransformMatrix{
		{1, 0, 0, 1},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1}}}
	newPuzzle := puzzle.Mutate(mutation)
	newPuzzle.Print()
	fmt.Println("Moving the Orange pice by another 1 along x axis:")
	newPuzzle.Mutate(mutation).Print()
}
