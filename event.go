package main

// event is a pawn option selected at a position.
type event struct {
	psn   *position // Position of an event
	poSlc *pawnOpt  // Pawn option selected at a position
}

// index returns the index of an event. If the event is not found, then -1 is returned.
func (hst history) index(brd board) int {
	// hst is NOT sorted, so do a linear search
	for i := range hst {
		if equalBoards(hst[i].psn.brd, brd) {
			return i
		}
	}

	return -1
}

// equalEvents compares two events by the joint comparison of the position and pawn option fields.
func equalEvents(evnt0, evnt1 *event) bool {
	return equalPositions(evnt0.psn, evnt1.psn) && equalPawnOpts(evnt0.poSlc, evnt1.poSlc)
}

// lessEvents compares two events. The position field is compared first, then the pawn options field is compared.
func lessEvents(evnt0, evnt1 *event) bool {
	if compareEvents(evnt0, evnt1) < 0 {
		return true
	}

	return false
}

// lessEqEvents compares two events. The position field is compared first, then the pawn options field is compared.
func lessEqEvents(evnt0, evnt1 *event) bool {
	if 0 < compareEvents(evnt0, evnt1) {
		return false
	}

	return true
}

// lessEvents compares two events returning -1 if evnt0 < evnt1, 0 if evnt0 = evnt1, and 1 if evnt0 > evnt1. The position field is compared first, then the pawn options field is compared.
func compareEvents(evnt0, evnt1 *event) int {
	r := comparePositions(evnt0.psn, evnt1.psn)
	if r == 0 {
		return comparePawnOpts(evnt0.poSlc, evnt1.poSlc) // Equal positions; pawn options should also be equal, but just in case...
	}

	return r
}
