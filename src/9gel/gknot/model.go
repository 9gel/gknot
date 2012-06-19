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

import (
	"fmt"
	"hash/fnv"
	"sort"
)

// A piece is defined by 5x7 matrix PieceGeom since each piece is 5 cells by
// 7 cells. It is done this way because the literal defining the piece will look
// just like the piece, instead of being a bunch of coordinates. The piece is
// laid flat on the x-y plane where x axis extends right
// and y axis extends upwards, starting at (0, 0) . PieceGeom[y][x] == 1 means
// there is a solid at (x, y); PieceGeom[y][x] == 0 means there's a void. Other
// values are invalid. The piece is also defined by a transformation matrix
// to transform it from the x-y plane to the starting position in the puzzle
// in 3 space.
type PieceGeom [5][7]uint8
type TransformMatrix [4][4]int
type PieceDefinition struct {
	Name      string
	// The ANSI escape color for printing in terminal. Also used as the ID of the piece.
	EscColor  uint8
	Geom      PieceGeom
	Transform TransformMatrix
}

// Piece represents a piece as a set of coordinates of its solid cells.
// Cells[0] is always the (0, 0) cell in the PieceDefinition.
type Cell [3]int
type Cells []Cell
type Piece struct {
	Definition *PieceDefinition
	Cells
}

type Pieces []*Piece

// Puzzle represents the current state of all the Pieces still entangled.
type CellMap map[Cell]*Piece
type Puzzle struct {
	// Source of truth of the pieces. All other fields in Puzzle must be consistent with this field.
	// The key is the EscColor of the piece.
	Pieces  map[uint8]*Piece
	// For looking up the Piece that a cell belongs to.
	CellMap
}

// Models a mutation to a Puzzle.
type Mutation struct {
}

// Error for when two pieces occupy the same cell.
type OverlapError struct {
	pieces Pieces
	cell   *Cell
}

func (e *OverlapError) Error() string {
	pieceNames := make([]string, len(e.pieces))
	for i, piece := range e.pieces {
		pieceNames[i] = piece.Definition.Name
	}
	return fmt.Sprintf("Overlapping pieces %v at %v.", pieceNames, *e.cell)
}

// Error for when two pieces have the same name.
type SameNameError struct {
	PieceName string
}

func (e *SameNameError) Error() string {
	return fmt.Sprintf("More than one piece with the same name %v.", e.PieceName)
}

// Error for when two pieces have the same ANSI escape color.
type SameEscColorError struct {
	EscColor uint8
}

func (e *SameEscColorError) Error() string {
	return fmt.Sprintf("More than one piece with the same ANSI Esc color %v.", e.EscColor)
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
	cells := make(Cells, 0, numCells)
	for y, row := range pieceDefn.Geom {
		for x, v := range row {
			if v == 1 {
				cells = append(cells, Cell{x, y, 0})
			}
		}
	}

	// Transform the cells.
	cells.transform(&pieceDefn.Transform)

	return &Piece{&pieceDefn, cells}
}

func (cells Cells) transform(transform *TransformMatrix) {
	for i, cell := range cells {
		cells[i] = cell.transform(transform)
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

func (puzzle *Puzzle) add(pieces ...*Piece) {
	nameSet := make(map[string]bool)
	escColorSet := make(map[uint8]bool)
	for _, piece := range pieces {
		for _, cell := range piece.Cells {
			if existPiece, ok := puzzle.CellMap[cell]; ok {
				// Panic because the default puzzle should not have overlapping cells.
				panic(&OverlapError{[]*Piece{piece, existPiece}, &cell})
			}
			if _, ok := nameSet[piece.Definition.Name]; ok {
				// Panic because the default puzzle should not have pieces with the
				// same names.
				panic(&SameNameError{piece.Definition.Name})
			}
			if _, ok := escColorSet[piece.Definition.EscColor]; ok {
				// Panic because the default puzzle should not have pieces with the
				// same esc color.
				panic(&SameEscColorError{piece.Definition.EscColor})
			}
			puzzle.CellMap[cell] = piece
		}
		puzzle.Pieces[piece.Definition.EscColor] = piece
	}
}

func NewPuzzle() *Puzzle {
	puzzle := &Puzzle{make(map[uint8]*Piece, 6), make(CellMap)}
	puzzle.add(BluePieceDef.Piece(),
		OrangePieceDef.Piece(),
		PurplePieceDef.Piece(),
		GreenPieceDef.Piece(),
		RedPieceDef.Piece(),
		YellowPieceDef.Piece())
	return puzzle
}

func (p Pieces) Len() int { return len(p) }
func (p Pieces) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type ByEscColor struct { Pieces }

func (p ByEscColor) Less(i, j int) bool { return p.Pieces[i].Definition.EscColor < p.Pieces[j].Definition.EscColor }

func (puzzle Puzzle) StateID() string {
	// Find the min x, min y and min z and move the puzzle to have min x, min y and min z == 0.
	// Calculate ID for each piece and concatenate.
	minX := 30
	minY := minX
	minZ := minY
	for cell := range puzzle.CellMap {
		if minX > cell[0] {
			minX = cell[0]
		}
		if minY > cell[1] {
			minY = cell[1]
		}
		if minZ > cell[2] {
			minZ = cell[2]
		}
	}
	pieces := make(Pieces, 0, len(puzzle.Pieces))
	for _, piece := range puzzle.Pieces {
		pieces = append(pieces, piece)
	}
	sort.Sort(ByEscColor{pieces})
	hash := fnv.New32()
	for _, piece := range pieces {
		hash.Write(piece.stateID(-minX, -minY, -minZ))
	}
	return fmt.Sprintf("%X", hash.Sum32())
}

func (piece Piece) stateID(shiftX, shiftY, shiftZ int) []byte {
	id := make([]byte, 0, 7)
	firstCell := piece.Cells[0]
	lastCell := piece.Cells[len(piece.Cells)-1]
	id = append(id, byte(piece.Definition.EscColor),
		byte(firstCell[0] + shiftX),
		byte(firstCell[1] + shiftY),
		byte(firstCell[2] + shiftZ),
		byte(lastCell[0] + shiftX),
		byte(lastCell[1] + shiftY),
		byte(lastCell[2] + shiftZ))
	return id
}
