// Copyright 2012 Google Inc. All Rights Reserved.
// Author: nigelchoi@google.com (Nigel Choi)

// Representation of Thinkfun's Gordion Knot.
// Pieces are laid flat, in the orientation as follows.
// When assembled, orient the cube like the picture on the packaging. Imagine
// an isometriclly projected cube:
//
// - Blue piece parallel to table at the bottom.
// - Orange piece (Magenta) parallel to table above Blue piece.
// - Purple piece (Cyan) left most.
// - Green piece parallel to and right of Purple piece.
// - Red piece right most.
// - Yellow piece parallel to and left of Red piece.
//
// We then use the OpenGL coordinate system. In the isometric projection, the
// origin is farthest from the view. From there, the x axis extends to the
// right side, the y axis extends upwards, and the z axis extends left.
//
// We then project the pieces along the x-y, y-z and x-z planes, in particular:
//
// - The Blue and Orange (Magenta) pieces project onto the x-z plane.
// - The Purple (Cyan) and Green pieces project onto the y-z plane.
// - The Red and Yellow pieces project onto the x-y plane.
package main

import (
	"fmt"
)

type PieceGeom [5][7]uint8
type Piece struct {
	Geom     PieceGeom
	EscColor uint8
}

var (
	BluePiece = Piece{
		// x-7, z-5
		PieceGeom{
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 0, 0, 1, 1},
			{1, 0, 1, 0, 0, 0, 1},
			{1, 1, 1, 1, 0, 1, 1}},
		34}
	OrangePiece = Piece{
		// x-7, z-5
		PieceGeom{
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 0, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 1, 1, 1, 1}},
		35}
	PurplePiece = Piece{
		// z-7, y-5
		PieceGeom{
			{1, 1, 1, 1, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 0, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1}},
		36}
	GreenPiece = Piece{
		// z-7, y-5
		PieceGeom{
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1}},
		32}
	RedPiece = Piece{
		// y-7, x-5
		PieceGeom{
			{1, 1, 0, 1, 0, 1, 1},
			{1, 0, 0, 1, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 1, 1, 1, 1}},
		31}
	YellowPiece = Piece{
		// y-7, x-5
		PieceGeom{
			{1, 1, 1, 1, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1}},
		33}
)

func (piece Piece) Print() {
	const block = '\u2588'
	fmt.Printf("%c[0;%dm", '\x1b', piece.EscColor)
	for _, row := range piece.Geom {
		for _, v := range row {
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
}
