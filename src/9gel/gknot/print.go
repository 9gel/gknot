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

type Coords2D [2]int
type ProjectInfo struct {
	Depth int
	*Piece
}
type ProjectedCells map[Coords2D] ProjectInfo
type Axis int
const (
	X = Axis(0)
	Y = Axis(1)
	Z = Axis(2)
)

func ProjectPuzzle(axis1, axis2 Axis, puzzle Puzzle) ProjectedCells {
	return ProjectPieces(axis1, axis2, puzzle.BluePiece, puzzle.OrangePiece, puzzle.PurplePiece, puzzle.GreenPiece, puzzle.RedPiece, puzzle.YellowPiece)
}

func ProjectPieces(axis1, axis2 Axis, pieces ...*Piece) ProjectedCells {
	axisDepth := 3 - axis1 - axis2
	projected := make(ProjectedCells)
	for _, piece := range pieces {
		for _, cell := range piece.Cells {
			coords := Coords2D{cell[axis1], cell[axis2]}
			existing, ok := projected[coords]
			if ok && existing.Depth > cell[axisDepth] {
				continue
			}
			projected[coords] = ProjectInfo{cell[axisDepth], piece}
		}
	}
	return projected
}

// Outputs the puzzle in 3 planar projections, along the planes:
// - x-y: x to the right, y upwards
// - y-z: y upwards, z to the left
// - x-z: x to the right, z downwards
func (puzzle Puzzle) Print() {
	xyProjected := ProjectPuzzle(X, Y, puzzle)
	yzProjected := ProjectPuzzle(Y, Z, puzzle)
	xzProjected := ProjectPuzzle(X, Z, puzzle)
	fmt.Println(xyProjected)
	fmt.Println(yzProjected)
	fmt.Println(xzProjected)
}
