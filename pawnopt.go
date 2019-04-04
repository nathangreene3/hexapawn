package main

import "sort"

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

// availPawnOpts returns a set of pawn options available at a position (m,n).
func availPawnOpts(brd board, st state) []*pawnOpt {
	po := make([]*pawnOpt, 0, 4)
	var (
		acts    []action // Set of actions for each position (i,j)
		wght    weight   // Weight to apply to each action
		actsLen int      // Number of available actions
	)

	for i := range brd {
		for j := range brd[i] {
			acts = availActions(i, j, brd, st)
			actsLen = len(acts)
			if actsLen == 0 {
				continue
			}

			wght = 1.0 / weight(actsLen)
			for k := range acts {
				po = append(po, &pawnOpt{m: i, n: j, act: acts[k], wght: wght})
			}
		}
	}

	return po
}

// equalPawnOpts returns true if each pawn option field is equal.
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
	case po0.wght != po1.wght:
		return false
	default:
		return true
	}
}

// copyPawnOpt returns a copy of a pawn option.
func copyPawnOpt(po *pawnOpt) *pawnOpt {
	return &pawnOpt{m: po.m, n: po.n, act: po.act, wght: po.wght}
}

// less compares the weight field of two pawn options in a set.
func (pos pawnOpts) less(i, j int) bool {
	return pos[i].wght < pos[j].wght
}
