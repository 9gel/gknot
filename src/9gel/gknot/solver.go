// Solving logic for Gordion Knot.
package gknot

import (
	"fmt"
)

type Translation [3]int

func (t Translation) TransformMatrix() TransformMatrix {
	return TransformMatrix{
		{1, 0, 0, t[0]},
		{0, 1, 0, t[1]},
		{0, 0, 1, t[2]},
		{0, 0, 0, 1}}
}

// Solve and print each step.
func (puzzle *Puzzle) Solve() {
	visitedStates := make(map[StateID]bool)
	puzzle.nextMoves(visitedStates, "", Translation{0, 0, 0}, nil)
}

func (puzzle *Puzzle) pushedPieces(cells Cells, xlate Translation, pushedPieces map[string]*Piece) {
	// Collect all the pieces immediately affected by this piece push.
	// Form an array of cells belonging to all the pieces.
	newCells := make(Cells, len(cells))
	copy(newCells, cells)
	appended := true
	for appended {
		appended = false
		for _, cell := range newCells {
			newCell := Cell{cell[0] + xlate[0], cell[1] + xlate[1], cell[2] + xlate[2]}
			if existPiece, hasPiece := puzzle.CellMap[newCell]; hasPiece {
				if _, pushed := pushedPieces[existPiece.Definition.Name]; !pushed {
					pushedPieces[existPiece.Definition.Name] = existPiece
					newCells = append(newCells, existPiece.Cells...)
					appended = true
				}
			}
		}
	}
}

func (puzzle *Puzzle) nextMoves(visitedStates map[StateID]bool, lastStateID StateID, lastXlate Translation, lastPieces []string) (seenState bool) {
	stateID := puzzle.StateID()
	if _, ok := visitedStates[stateID]; ok {
		// Have seen this state already.
		return false
	}
	visitedStates[stateID] = true
	if lastStateID != "" {
		fmt.Println("From", lastStateID, "Mutate", lastXlate, "Pieces", lastPieces)
	}
	puzzle.Print()
	// For each piece, see if there is translation in any of the 6 directions, pushing other pieces along
	// if necessary. If all pieces are pushed, it is not a valid movement.
	// TODO: look for rotation opportunities.
	translations := []Translation{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1}}
	hasMoreMoves := false
	for _, piece := range puzzle.Pieces {
		for _, xlate := range translations {
			piecesToMutate := map[string]*Piece{piece.Definition.Name: piece}
			puzzle.pushedPieces(piece.Cells, xlate, piecesToMutate)
			if numMutations := len(piecesToMutate); numMutations < len(puzzle.Pieces) {
				mutations := make([]Mutation, 0, numMutations)
				mutatedPieceNames := make([]string, 0, len(piecesToMutate))
				for _, toMutate := range piecesToMutate {
					mutations = append(mutations, Mutation{toMutate.Definition.EscColor, xlate.TransformMatrix()})
					mutatedPieceNames = append(mutatedPieceNames, toMutate.Definition.Name)
				}
				if puzzle.Mutate(mutations...).nextMoves(visitedStates, stateID, xlate, mutatedPieceNames) {
					hasMoreMoves = true
				}
			}
		}
	}
	if !hasMoreMoves {
		fmt.Println("No more moves beyond state", stateID)
	}
	return true
}
