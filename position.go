package main

import (
	"bytes"
)

// position joins a board, state, and a set of available pawn options with a pawn
// option selection.
type position struct {
	brd board    // Board position
	st  state    // State of the game
	pos pawnOpts // Available pawn options
}

func (psn *position) String() string {
	return string(psn.toBytes())
}

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
