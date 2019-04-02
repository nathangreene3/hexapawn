package main

import "log"

// history is a set of positions that occur in a single game.
type history []*position

// position joins a board, state, and a set of available pawn options with a pawn option selection.
type position struct {
	b     board      // Board position
	s     state      // State of the game
	poSlc *pawnOpt   // Pawn option selected; set to nil before selection has been made
	po    []*pawnOpt // Available pawn options
}

// copyPosition returns a copy of a postion.
func copyPosition(psn *position) *position {
	cpy := &position{b: copyBoard(psn.b), s: psn.s, poSlc: psn.poSlc}
	n := len(psn.po)

	switch n {
	case 0:
		if psn.poSlc != nil {
			log.Fatal("copyPosition: cannot have selected a pawn option with zero pawn options available")
		}
	default:
		cpy.po = make([]*pawnOpt, 0, n)
		for i := range psn.po {
			cpy.po = append(cpy.po, psn.po[i])
		}
	}

	return cpy
}

// equalPositions returns true if each field is equal and false if otherwise.
func equalPositions(p, q *position) bool {
	switch {
	case p.s != q.s:
		return false
	case !equalPawnOpts(p.poSlc, q.poSlc):
		return false
	case len(p.po) != len(q.po):
		return false
	case !equalBoards(p.b, q.b):
		return false
	default:
		for i := range p.po {
			if !equalPawnOpts(p.po[i], q.po[i]) {
				return false
			}
		}
		return true
	}
}
