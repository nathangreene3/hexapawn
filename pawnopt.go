package main

import (
	"fmt"
	"sort"
)

// TODO: redefine to something more appropriate
// weight is a probability value on the range [0,1].
type weight float64

// pawnOpts is a set of pawn options.
type pawnOpts []*pawnOpt

// pawnOpt is an available action at a position (m,n) with a probability weight of
// being selected.
type pawnOpt struct {
	m    int    // Row index in board
	n    int    // Column index in board
	act  action // Available action
	wght weight // Probability of selecting action
}

func (po *pawnOpt) String() string {
	switch po.act {
	case forward:
		return fmt.Sprintf("pawnOpt: forward at (%d,%d), weight: %0.2f\n", po.m, po.n, po.wght)
	case captureLeft:
		return fmt.Sprintf("pawnOpt: capture-left at (%d,%d), weight: %0.2f\n", po.m, po.n, po.wght)
	case captureRight:
		return fmt.Sprintf("pawnOpt: capture-right at (%d,%d), weight: %0.2f\n", po.m, po.n, po.wght)
	default:
		return fmt.Sprintf("pawnOpt: unknown action at (%d,%d), weight: %0.2f\n", po.m, po.n, po.wght)
	}
}

// func (po *pawnOpt)toBytes()[]byte{
// 	switch po.act {
// 	case forward:
// 		return fmt.Sprintf("forward pawnOpt at (%d,%d): action: %d weight: %0.2f", po.m, po.n, po.act, po.wght)
// 	case captureLeft:
// 		return fmt.Sprintf("capture-left pawnOpt at (%d,%d): action: %d weight: %0.2f", po.m, po.n, po.act, po.wght)
// 	case captureRight:
// 		return fmt.Sprintf("capture-right pawnOpt at (%d,%d): action: %d weight: %0.2f", po.m, po.n, po.act, po.wght)
// 	default:
// 		return fmt.Sprintf("unknown pawnOpt at (%d,%d): action: %d weight: %0.2f", po.m, po.n, po.act, po.wght)
// 	}
// }

// insert a pawn option into a set. The pawn option will be copied.
func (pos pawnOpts) insert(po *pawnOpt) int {
	switch len(pos) {
	case 0:
		pos = append(pos, copyPawnOpt(po))
		return 0
	case pos.index(po):
		pos = append(pos, copyPawnOpt(po))
		sort.Slice(pos, pos.less)
		return pos.index(po)
	default:
		return pos.index(po)
	}
}

// index returns the index a pawn option is found in a set of pawn options. If the
// pawn option is not found, len(pos) is returned.
func (pos pawnOpts) index(po *pawnOpt) int {
	return sort.Search(len(pos), func(i int) bool { return equalPawnOpts(pos[i], po) })
}

// availPawnOpts returns a set of pawn options available given a board state.
func availPawnOpts(brd board, st state) pawnOpts {
	pos := make(pawnOpts, 0, 4)
	var (
		actsLen   int      // Number of available actions per pawn
		actsCount int      // Total number of actions available per board state
		acts      []action // Set of actions for each position (i,j)
	)

	for i := range brd {
		for j := range brd[i] {
			acts = availActions(i, j, brd, st)
			actsLen = len(acts)
			if actsLen == 0 {
				continue
			}

			actsCount += actsLen

			for k := range acts {
				pos = append(pos, &pawnOpt{m: i, n: j, act: acts[k]})
			}
		}
	}

	var wght weight
	if actsCount < 2 {
		wght = 1
	} else {
		wght = 1 / weight(actsCount)
	}

	for i := range pos {
		pos[i].wght = wght
	}

	return pos
}

// equalPawnOpts returns true if each pawn option field is equal, EXCEPT for the weight field.
func equalPawnOpts(po0, po1 *pawnOpt) bool {
	switch {
	case po0 == nil:
		return po1 == nil
	case po1 == nil:
		return po0 == nil
	case po0.m != po1.m:
		return false
	case po0.n != po1.n:
		return false
	case po0.act != po1.act:
		return false
	// Weight is not an important field for searching/comparing
	// case po0.wght != po1.wght:
	// 	return false
	default:
		return true
	}
}

// copyPawnOpt returns a copy of a pawn option.
func copyPawnOpt(po *pawnOpt) *pawnOpt {
	return &pawnOpt{m: po.m, n: po.n, act: po.act, wght: po.wght}
}

// less compares two pawn options on the position (m,n) and the action field in that order.
func (pos pawnOpts) less(i, j int) bool {
	if pos[i].m < pos[j].m {
		return true
	}

	if pos[j].m < pos[i].m {
		return false
	}

	// pos[i].m == pos[j].m
	if pos[i].n < pos[j].n {
		return true
	}

	if pos[j].n < pos[i].n {
		return false
	}

	// pos[i].n == pos[j].n
	return pos[i].act < pos[j].act
}

func lessPawnOpts(po0, po1 *pawnOpt) bool {
	switch {
	case po0 == nil, po1 == nil:
		panic("lessPawnOpts: cannot compare nil pawn options")
	case po0.m < po1.m:
		return true
	case po1.m < po0.m:
		return false
	case po0.n < po1.n:
		return true
	case po1.n < po0.n:
		return false
	default:
		return po0.act < po1.act
	}
}

func comparePawnOpts(po0, po1 *pawnOpt) int {
	switch {
	case po0 == nil, po1 == nil:
		panic("lessPawnOpts: cannot compare nil pawn options")
	case po0.m < po1.m:
		return -1
	case po1.m < po0.m:
		return 1
	case po0.n < po1.n:
		return -1
	case po1.n < po0.n:
		return 1
	case po0.act < po1.act:
		return -1
	case po1.act < po0.act:
		return 1
	default:
		return 0
	}
}
