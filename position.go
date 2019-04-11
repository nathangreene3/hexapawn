package main

import (
	"bytes"
)

// position joins a state, board, and a set of available pawn options.
type position struct {
	st  state    // State of the game
	brd board    // Board position
	pos pawnOpts // Available pawn options
}

// String returns the formated representation of a position.
func (psn *position) String() string {
	return string(psn.toBytes())
}

// toBytes returns the formated representation of a position.
func (psn *position) toBytes() []byte {
	buf := bytes.Buffer{}

	buf.Write(psn.brd.toBytes())
	buf.Write([]byte{'\n', byte(psn.st)})

	for i := range psn.pos {
		buf.WriteString(psn.pos[i].String())
	}

	buf.WriteByte('\n')

	return buf.Bytes()
}

// copyPosition returns a copy of a postion.
func copyPosition(psn *position) *position {
	cpy := &position{brd: copyBoard(psn.brd), st: psn.st, pos: make(pawnOpts, 0, len(psn.pos))}
	for i := range psn.pos {
		cpy.pos = append(cpy.pos, copyPawnOpt(psn.pos[i]))
	}

	return cpy
}

// equalPositions returns true if each field is equal and false if otherwise.
func equalPositions(psn0, psn1 *position) bool {
	switch {
	case psn0.st != psn1.st:
		return false
	case len(psn0.pos) != len(psn1.pos):
		return false
	case !equalBoards(psn0.brd, psn1.brd):
		return false
	default:
		for i := range psn0.pos {
			if !equalPawnOpts(psn0.pos[i], psn1.pos[i]) {
				return false
			}
		}

		return true
	}
}

// lessPositions compares two positions. The state field is checked first, then boards are compared.
func lessPositions(psn0, psn1 *position) bool {
	if comparePositions(psn0, psn1) < 0 {
		return true
	}

	return false
}

// lessEqPositions compares two positions. The state field is checked first, then boards are compared.
func lessEqPositions(psn0, psn1 *position) bool {
	if 0 < comparePositions(psn0, psn1) {
		return false
	}

	return true
}

// comparePositions compares two positions returning -1 if psn0 < psn1, 0 if psn0 = psn1, and 1 if psn0 > psn1. The state property is compared first, then the board field is compared. Panics if either position is nil.
func comparePositions(psn0, psn1 *position) int {
	switch {
	case psn0 == nil, psn1 == nil:
		panic("lessPositions: cannot compare nil positions")
	case psn0.st < psn1.st:
		return -1
	case psn1.st < psn0.st:
		return 1
	default:
		return compareBoards(psn0.brd, psn1.brd) // states are equal
	}
}
