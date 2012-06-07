package gknot

import "testing"

func checkPiece(t *testing.T, piece Piece, expected PieceCells) {
	for i, cell := range piece.Cells {
		if cell != expected[i] {
			t.Errorf("%v piece[%v] = %v, expected %v", piece.Definition.Name, i, cell, expected[i])
		}
	}
}

func TestBluePiece_cells(t *testing.T) {
	checkPiece(t, BluePieceDef.Piece(), PieceCells{
		{0,2,1},{1,2,1},{2,2,1},{3,2,1},{4,2,1},{5,2,1},{6,2,1},
		{0,2,2},                                        {6,2,2},
		{0,2,3},{1,2,3},                        {5,2,3},{6,2,3},
		{0,2,4},        {2,2,4},                        {6,2,4},
		{0,2,5},{1,2,5},{2,2,5},{3,2,5},        {5,2,5},{6,2,5}})
}

func TestOrangePiece_cells(t *testing.T) {
	checkPiece(t, OrangePieceDef.Piece(), PieceCells{
		{0,4,1},{1,4,1},{2,4,1},{3,4,1},{4,4,1},{5,4,1},{6,4,1},
		{0,4,2},                                        {6,4,2},
		{0,4,3},{1,4,3},                        {5,4,3},{6,4,3},
		{0,4,4},                                        {6,4,4},
		{0,4,5},{1,4,5},        {3,4,5},{4,4,5},{5,4,5},{6,4,5}})
}

func TestPurplePiece_cells(t *testing.T) {
	checkPiece(t, PurplePieceDef.Piece(), PieceCells{
		{2,1,0},{2,1,1},{2,1,2},{2,1,3},        {2,1,5},{2,1,6},
		{2,2,0},                                        {2,2,6},
		{2,3,0},{2,3,1},                        {2,3,5},{2,3,6},
		{2,4,0},                                        {2,4,6},
		{2,5,0},{2,5,1},{2,5,2},{2,5,3},{2,5,4},{2,5,5},{2,5,6}})
}

func TestGreenPiece_cells(t *testing.T) {
	checkPiece(t, GreenPieceDef.Piece(), PieceCells{
		{4,1,0},{4,1,1},{4,1,2},{4,1,3},{4,1,4},{4,1,5},{4,1,6},
		{4,2,0},                                        {4,2,6},
		{4,3,0},{4,3,1},{4,3,2},{4,3,3},{4,3,4},{4,3,5},{4,3,6},
		{4,4,0},                                        {4,4,6},
		{4,5,0},{4,5,1},{4,5,2},{4,5,3},{4,5,4},{4,5,5},{4,5,6}})
}

func TestRedPiece_cells(t *testing.T) {
	checkPiece(t, RedPieceDef.Piece(), PieceCells{
		{1,6,2},{1,5,2},        {1,3,2},        {1,1,2},{1,0,2},
		{2,6,2},                {2,3,2},                {2,0,2},
		{3,6,2},{3,5,2},{3,4,2},{3,3,2},{3,2,2},{3,1,2},{3,0,2},
		{4,6,2},                                        {4,0,2},
		{5,6,2},{5,5,2},{5,4,2},{5,3,2},        {5,1,2},{5,0,2}})
}

func TestYellowPiece_cells(t *testing.T) {
	checkPiece(t, YellowPieceDef.Piece(), PieceCells{
		{1,6,4},{1,5,4},        {1,3,4},{1,2,4},{1,1,4},{1,0,4},
		{2,6,4},                                        {2,0,4},
		{3,6,4},{3,5,4},        {3,3,4},{3,2,4},{3,1,4},{3,0,4},
		{4,6,4},                                        {4,0,4},
		{5,6,4},{5,5,4},{5,4,4},{5,3,4},{5,2,4},{5,1,4},{5,0,4}})
}
