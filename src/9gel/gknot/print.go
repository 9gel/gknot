package gknot

import "fmt"

const (
	block = '\u2588'
	esc = '\x1b'
)

// Prints a piece laid flat.
func (piece PieceDefinition) Print() {
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
type ProjectedCell struct {
	Depth int
	*Piece
}
type ProjectedCells map[Coords2D] ProjectedCell
type Axis int
const (
	X = Axis(0)
	Y = Axis(1)
	Z = Axis(2)
)

// 2x3 matrix.
type Transform2D [2][3]int

func ProjectPuzzle(axis1, axis2 Axis, puzzle Puzzle) ProjectedCells {
	depthAxis := 3 - axis1 - axis2
	projected := make(ProjectedCells)
	for _, piece := range puzzle.Pieces {
		for _, cell := range piece.Cells {
			coords := Coords2D{cell[axis1], cell[axis2]}
			existing, ok := projected[coords]
			if ok && existing.Depth > cell[depthAxis] {
				continue
			}
			projected[coords] = ProjectedCell{cell[depthAxis], piece}
		}
	}
	return projected
}

func (projected ProjectedCells) axesMax() (axis1Max, axis2Max int) {
  axis1Max = -20
	axis2Max = -20
  for coords := range projected {
    if coords[0] > axis1Max {
      axis1Max = coords[0]
    }
    if coords[1] > axis2Max {
      axis2Max = coords[1]
    }
  }
	return
}

// Transforms the cells using the transformation and add to screenCells.
func (screenCells ProjectedCells) transformAndAddCells(transform Transform2D, cells ProjectedCells) {
	for cellCoords, cell := range cells {
		var screenCoords Coords2D
		for axis, row := range transform {
			for i, v := range cellCoords {
				screenCoords[axis] += row[i] * v
			}
			screenCoords[axis] += row[2]
		}
		screenCells[screenCoords] = cell
	}
}

// Outputs the puzzle in 3 projections, along the planes:
// - x-y: x to the right, y upwards
// - y-z: y upwards, z to the left
// - x-z: x to the right, z downwards
func (puzzle Puzzle) Print() {
	xyProjected := ProjectPuzzle(X, Y, puzzle)
	yzProjected := ProjectPuzzle(Y, Z, puzzle)
	xzProjected := ProjectPuzzle(X, Z, puzzle)

  // The terminal prints from top left to bottom right. i.e. the
  // coordinate system is x-y where x axis goes to the right and
  // y axis goes downwards, starting at (0,0). Translate and reflect the
  // 2D projections to this coordinate system and print.
  // Print all 3 projections side-by-side, each occupying at most 20 spaces
  // along the x axis, 60 spaces in total. They will occupy at most 20 
  // spaces down the y axis.
	screenCells := make(ProjectedCells)

	_, xyMaxY := xyProjected.axesMax()
  screenCells.transformAndAddCells(Transform2D{
    {1, 0, 0},
    {0, -1, xyMaxY}}, xyProjected)

	yzMaxY, yzMaxZ := yzProjected.axesMax()
	screenCells.transformAndAddCells(Transform2D{
		{0, -1, yzMaxZ + 20},
		{-1, 0, yzMaxY}}, yzProjected)

	screenCells.transformAndAddCells(Transform2D{
		{1, 0, 40},
		{0, 1, 0}}, xzProjected)

	screenMaxX, screenMaxY := screenCells.axesMax()
	fmt.Printf("%c[1mx-y%37cy-z%37cx-z%c[0m\n", esc, ' ', ' ', esc)
	for y := 0; y <= screenMaxY; y++ {
		spacer := ""
		for x := 0; x <= screenMaxX; x++ {
			cell, ok := screenCells[Coords2D{x, y}]
			if ok {
				fmt.Printf("%v%c[0;%dm%c%c%c[0m", spacer, esc, cell.Piece.Definition.EscColor, block, block, esc)
				spacer = ""
			} else {
				spacer += "  "
			}
		}
		fmt.Println()
	}
}
