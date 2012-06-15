package gknot

import (
	"reflect"
	"testing"
)

func checkPiece(t *testing.T, piece *Piece, expected Cells) {
	for i, cell := range piece.Cells {
		if cell != expected[i] {
			t.Errorf("%v piece[%v] = %v, expected %v", piece.Definition.Name, i, cell, expected[i])
		}
	}
}

func TestBluePiece_cells(t *testing.T) {
	checkPiece(t, BluePieceDef.Piece(), Cells{
		{0,2,1},{1,2,1},{2,2,1},{3,2,1},{4,2,1},{5,2,1},{6,2,1},
		{0,2,2},                                        {6,2,2},
		{0,2,3},{1,2,3},                        {5,2,3},{6,2,3},
		{0,2,4},        {2,2,4},                        {6,2,4},
		{0,2,5},{1,2,5},{2,2,5},{3,2,5},        {5,2,5},{6,2,5}})
}

func TestOrangePiece_cells(t *testing.T) {
	checkPiece(t, OrangePieceDef.Piece(), Cells{
		{0,4,1},{1,4,1},{2,4,1},{3,4,1},{4,4,1},{5,4,1},{6,4,1},
		{0,4,2},                                        {6,4,2},
		{0,4,3},{1,4,3},                        {5,4,3},{6,4,3},
		{0,4,4},                                        {6,4,4},
		{0,4,5},{1,4,5},        {3,4,5},{4,4,5},{5,4,5},{6,4,5}})
}

func TestPurplePiece_cells(t *testing.T) {
	checkPiece(t, PurplePieceDef.Piece(), Cells{
		{2,1,0},{2,1,1},{2,1,2},{2,1,3},        {2,1,5},{2,1,6},
		{2,2,0},                                        {2,2,6},
		{2,3,0},{2,3,1},                        {2,3,5},{2,3,6},
		{2,4,0},                                        {2,4,6},
		{2,5,0},{2,5,1},{2,5,2},{2,5,3},{2,5,4},{2,5,5},{2,5,6}})
}

func TestGreenPiece_cells(t *testing.T) {
	checkPiece(t, GreenPieceDef.Piece(), Cells{
		{4,1,0},{4,1,1},{4,1,2},{4,1,3},{4,1,4},{4,1,5},{4,1,6},
		{4,2,0},                                        {4,2,6},
		{4,3,0},{4,3,1},{4,3,2},{4,3,3},{4,3,4},{4,3,5},{4,3,6},
		{4,4,0},                                        {4,4,6},
		{4,5,0},{4,5,1},{4,5,2},{4,5,3},{4,5,4},{4,5,5},{4,5,6}})
}

func TestRedPiece_cells(t *testing.T) {
	checkPiece(t, RedPieceDef.Piece(), Cells{
		{1,6,2},{1,5,2},        {1,3,2},        {1,1,2},{1,0,2},
		{2,6,2},                {2,3,2},                {2,0,2},
		{3,6,2},{3,5,2},{3,4,2},{3,3,2},{3,2,2},{3,1,2},{3,0,2},
		{4,6,2},                                        {4,0,2},
		{5,6,2},{5,5,2},{5,4,2},{5,3,2},        {5,1,2},{5,0,2}})
}

func TestYellowPiece_cells(t *testing.T) {
	checkPiece(t, YellowPieceDef.Piece(), Cells{
		{1,6,4},{1,5,4},        {1,3,4},{1,2,4},{1,1,4},{1,0,4},
		{2,6,4},                                        {2,0,4},
		{3,6,4},{3,5,4},        {3,3,4},{3,2,4},{3,1,4},{3,0,4},
		{4,6,4},                                        {4,0,4},
		{5,6,4},{5,5,4},{5,4,4},{5,3,4},{5,2,4},{5,1,4},{5,0,4}})
}

func TestOverlapError(t *testing.T) {
	badRedPiece := PieceDefinition{
		"Red",
		31,
		PieceGeom{
			{1, 1, 0, 1, 0, 1, 1},
			{1, 0, 0, 1, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1},
			{1, 1, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 0, 1, 1}},
		// Same transform as RedPiece.
		TransformMatrix{
			{0, 1, 0, 1},
			{-1, 0, 0, 6},
			{0, 0, 1, 2},
			{0, 0, 0, 1}}}
	puzzle := &Puzzle{make([]*Piece, 0, 6), make(CellMap)}

	defer func() {
		if err := recover(); err != nil {
			overlapErr, ok := err.(*OverlapError)
			if !ok {
				t.Fatalf("Expected OverlapError, actual %v.", reflect.TypeOf(overlapErr))
			}
			if piecesLen := len(overlapErr.pieces); piecesLen != 2 {
				t.Fatalf("OverlapError has incorrect number of pieces: expected 2, actual %v.", piecesLen)
			}
			expectedCell := Cell{4, 5, 2}
			if overlapCell := overlapErr.cell; *overlapCell != expectedCell {
				t.Fatalf("OverlapError has incorrect overlap cell; expected %v, actual %v.", expectedCell, overlapCell)
			}
			if errMsg := overlapErr.Error(); errMsg != "Overlapping pieces [Red Green] at [4 5 2]." {
				t.Fatalf("Incorrect error message for OverlapError; actual message '%v'.", errMsg)
			}
		}
	}()
	puzzle.add(GreenPieceDef.Piece(), badRedPiece.Piece())
	t.Fatal("Should have paniked with OverlapError since two pieces overlap.")
}
