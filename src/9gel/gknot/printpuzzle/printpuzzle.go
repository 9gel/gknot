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

	fmt.Println("Moving the all but Orange pice by -1 along x axis:")
	transform := gknot.TransformMatrix{
		{1, 0, 0, -1},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1}}
	mutations := []gknot.Mutation{{31, transform}, {32, transform}, {33, transform}, {34, transform}, {36, transform}}
	newPuzzle.Mutate(mutations...).Print()
}
