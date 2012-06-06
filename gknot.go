// Author: nigelchoi@google.com (Nigel Choi)

// Representation and solver of Thinkfun's Gordion Knot.
// When assembled, orient the cube like the picture on the packaging. Imagine
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
package main

import (
	"fmt"
)

type PieceGeom [5][7]uint8
type TransformMatrix [4][4]int8
type Piece struct {
	EscColor uint8
	Geom     PieceGeom
	Transform TransformMatrix
}
type Cell [3]int8
type PieceCells []Cell

var (
	BluePiece = Piece{
		36,  // Cyan.
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
	OrangePiece = Piece{
		35,  // Magenta.
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
	PurplePiece = Piece{
		34,  // Blue/Purple in Terminal.app.
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
	GreenPiece = Piece{
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
	RedPiece = Piece{
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
	YellowPiece = Piece{
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

func (piece Piece) Print() {
	const block = '\u2588'
	fmt.Printf("%c[0;%dm", '\x1b', piece.EscColor)
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
	fmt.Printf("%c[0m", '\x1b')
}

func (piece Piece) PieceCells() PieceCells {
	numCells := 0
	for _, row := range piece.Geom {
		numCells += len(row)
	}
	pieceCells := make(PieceCells, numCells)
	return pieceCells
}

func main() {
	BluePiece.Print()
	fmt.Println()
	OrangePiece.Print()
	fmt.Println()
	PurplePiece.Print()
	fmt.Println()
	GreenPiece.Print()
	fmt.Println()
	RedPiece.Print()
	fmt.Println()
	YellowPiece.Print()

	blueCells := BluePiece.PieceCells()
	fmt.Println(blueCells)
}
