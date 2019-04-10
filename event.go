package main

// history is a set of positions that occur in a single game.
type history []*event

// event is a pawn option selected at a position.
type event struct {
	psn   *position // Position of an event
	poSlc *pawnOpt  // Pawn option selected at a position
}

// index returns the index of a position defined by a given board. If the board is
// not found, then len(evnt) is returned.
func (hst history) index(brd board) int {
	for i := range hst {
		if equalBoards(hst[i].psn.brd, brd) {
			return i
		}
	}

	return len(hst)
}

// equalEvents returns true if two events have equal positions and equal selected
// pawn options.
func equalEvents(evnt0, evnt1 *event) bool {
	return equalPositions(evnt0.psn, evnt1.psn) && equalPawnOpts(evnt0.poSlc, evnt1.poSlc)
}

func compareEvents(evnt0, evnt1 *event) int {
	r := comparePositions(evnt0.psn, evnt1.psn)
	if r == 0 {
		return comparePawnOpts(evnt0.poSlc, evnt1.poSlc) // Equal positions
	}

	return r
}
