// Author: nigelchoi@google.com (Nigel Choi)

// Representation Thinkfun's Gordion Knot.
// When assembled, orient the cube like follows. Imagine
// an isometriclly projected cube on a right-handed coordinate system:
//
// - Blue (Cyan) piece parallel to table at the bottom.
// - Orange piece (Magenta) parallel to table above Blue piece.
// - Purple piece (Blue/Purple) left most.
// - Green piece parallel to and right of Purple piece.
// - Red piece right most.
// - Yellow piece parallel to and left of Red piece.
//
// The x axis extends to the right side, the y axis extends
// upwards, and the z axis extends left.
//
// In this coordinate system:
//
// - The Blue and Orange (Magenta) pieces are parallel to the x-z plane.
// - The Purple (Cyan) and Green pieces project onto the y-z plane.
// - The Red and Yellow pieces project onto the x-y plane.
//
// Positions:
//
// - The assembled puzzle rests on the x=0, y=0 and z=0 planes. Each piece is 1
//   unit high, 5 x 7 unit in dimensions.
// - The representation below represents pieces laid flat on the x-y plane,
//   where the matrix indices correspond to coordinates in x-y. i.e. Geom[y][x]
//   corresponds to what's in (x, y). If Geom[y][x] == 1 there's a solid; if
//   Geom[y][x] == 0 there's a void.
// - The transformation matrix places the piece where they should be in 3-space.
package gknot

import "fmt"

type PieceGeom [5][7]uint8
type TransformMatrix [4][4]int
type PieceDefinition struct {
	Name      string
	EscColor  uint8
	Geom      PieceGeom
	Transform TransformMatrix
}
type Cell [3]int
type PieceCells []Cell
type Piece struct {
	Definition *PieceDefinition
	Cells      PieceCells
}
type CellMap map[Cell]*Piece
type Puzzle struct {
	Pieces  []*Piece
	CellMap
}

type OverlapError struct {
	pieces []*Piece
	cell   *Cell
}

var (
	BluePieceDef = PieceDefinition{
		"Blue",
		36, // Cyan.
		PieceGeom{
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 0, 0, 1, 1},
			{1, 0, 1, 0, 0, 0, 1},
			{1, 1, 1, 1, 0, 1, 1}},
		// Rotate 90d about x axis,
		// Translate along y by 2,
		// Translate along z by 1.
		TransformMatrix{
			{1, 0, 0, 0},
			{0, 0, -1, 2},
			{0, 1, 0, 1},
			{0, 0, 0, 1}}}
	OrangePieceDef = PieceDefinition{
		"Orange",
		35, // Magenta.
		PieceGeom{
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 0, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 1, 1, 1, 1}},
		// Rotate 90d about x axis,
		// Translate along y by 4,
		// Translate along z by 1.
		TransformMatrix{
			{1, 0, 0, 0},
			{0, 0, -1, 4},
			{0, 1, 0, 1},
			{0, 0, 0, 1}}}
	PurplePieceDef = PieceDefinition{
		"Purple",
		34, // Blue/Purple in Terminal.app.
		PieceGeom{
			{1, 1, 1, 1, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 0, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1}},
		// Rotate -90d about y axis,
		// Translate along y by 1,
		// Translate along x by 2.
		TransformMatrix{
			{0, 0, -1, 2},
			{0, 1, 0, 1},
			{1, 0, 0, 0},
			{0, 0, 0, 1}}}
	GreenPieceDef = PieceDefinition{
		"Green",
		32,
		PieceGeom{
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1}},
		// Rotate -90d about y axis,
		// Translate along y by 1,
		// Translate along x by 4.
		TransformMatrix{
			{0, 0, -1, 4},
			{0, 1, 0, 1},
			{1, 0, 0, 0},
			{0, 0, 0, 1}}}
	RedPieceDef = PieceDefinition{
		"Red",
		31,
		PieceGeom{
			{1, 1, 0, 1, 0, 1, 1},
			{1, 0, 0, 1, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 0, 1, 1}},
		// Rotate -90d about z axis,
		// Translate along x by 1,
		// Translate along y by 6,
		// Translate along z by 2.
		TransformMatrix{
			{0, 1, 0, 1},
			{-1, 0, 0, 6},
			{0, 0, 1, 2},
			{0, 0, 0, 1}}}
	YellowPieceDef = PieceDefinition{
		"Yellow",
		33,
		PieceGeom{
			{1, 1, 0, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1}},
		// Rotate -90d about z axis,
		// Translate along x by 1,
		// Translate along y by 6,
		// Translate along z by 4.
		TransformMatrix{
			{0, 1, 0, 1},
			{-1, 0, 0, 6},
			{0, 0, 1, 4},
			{0, 0, 0, 1}}}
)

func (pieceDefn PieceDefinition) Piece() *Piece {
	// Build list of cells.
	numCells := 0
	for _, row := range pieceDefn.Geom {
		numCells += len(row)
	}
	pieceCells := make(PieceCells, 0, numCells)
	for y, row := range pieceDefn.Geom {
		for x, v := range row {
			if v == 1 {
				pieceCells = append(pieceCells, Cell{x, y, 0})
			}
		}
	}

	// Transform the cells.
	pieceCells.transform(&pieceDefn.Transform)

	return &Piece{&pieceDefn, pieceCells}
}

func (pieceCells PieceCells) transform(transform *TransformMatrix) {
	for i, cell := range pieceCells {
		pieceCells[i] = cell.transform(transform)
	}
}

func (cell Cell) transform(transform *TransformMatrix) Cell {
	newCell := Cell{0, 0, 0}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			newCell[i] += transform[i][j] * cell[j]
		}
		newCell[i] += transform[i][3]
	}
	return newCell
}

func (cellMap CellMap) add(pieces ...*Piece) {
	for _, piece := range pieces {
		for _, cell := range piece.Cells {
			if existPiece, ok := cellMap[cell]; ok {
				// Panic because the default puzzle should not have overlapping cells.
				panic(&OverlapError{[]*Piece{piece, existPiece}, &cell})
			}
			cellMap[cell] = piece
		}
	}
}

func NewPuzzle() *Puzzle {
	// Build the cell map
	bluePiece := BluePieceDef.Piece()
	orangePiece := OrangePieceDef.Piece()
	purplePiece := PurplePieceDef.Piece()
	greenPiece := GreenPieceDef.Piece()
	redPiece := RedPieceDef.Piece()
	yellowPiece := YellowPieceDef.Piece()
	pieces := []*Piece{bluePiece, orangePiece, purplePiece, greenPiece, redPiece, yellowPiece}
	cellMap := make(CellMap)
	cellMap.add(bluePiece, orangePiece, purplePiece, greenPiece, redPiece, yellowPiece)

	return &Puzzle{pieces, cellMap}
}

func (e *OverlapError) Error() string {
	pieceNames := make([]string, len(e.pieces))
	for i, piece := range e.pieces {
		pieceNames[i] = piece.Definition.Name
	}
	return fmt.Sprintf("Overlapping pieces %v at %v", pieceNames, *e.cell)
}
