package main

import "sort"

// TODO: redefine to something more appropriate
// weight is a probability value on the range [0,1].
type weight float64

// pawnOpts is a set of pawn options.
type pawnOpts []*pawnOpt

// pawnOpt is an available action at a position (m,n) with a probability weight of being selected.
type pawnOpt struct {
	m   int    // Row index in board
	n   int    // Column index in board
	act action // Available action
	w   weight // TODO: redefine this field to something more appropriate
}

// insert a pawn option into a set. The pawn option will be copied.
func (pos pawnOpts) insert(po *pawnOpt) {
	switch len(pos) {
	case 0:
		pos = append(pos, po)
	case pos.index(po):
		pos = append(pos, po)
		sort.Sort(pos)
	}
}

// index returns the index a pawn option is found in a set of pawn options. If the pawn option is not found, len(pos) is returned.
func (pos pawnOpts) index(po *pawnOpt) int {
	return sort.Search(len(pos), func(i int) bool { return equalPawnOpts(pos[i], po) })
}

// availPawnOpts returns a set of pawn options available at a position (m,n).
func availPawnOpts(brd board, st state) []*pawnOpt {
	po := make([]*pawnOpt, 0, 4)
	var a []action // Set of actions for each position (i,j)
	var w weight   // Weight to apply to each action
	var d weight   // Difference in each action weight
	var n int      // Number of available actions

	for i := range brd {
		for j := range brd[i] {
			a = availActions(i, j, brd, st)
			n = len(a)
			if n == 0 {
				continue
			}

			d = weight(1.0 / float64(n))
			for k := range a {
				w += d
				po = append(po, &pawnOpt{m: i, n: j, act: a[k], w: w})
			}

			w = 0
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
	case po0.w != po1.w:
		return false
	default:
		return true
	}
}

func copyPawnOpt(po *pawnOpt) *pawnOpt {
	return &pawnOpt{m: po.m, n: po.n, act: po.act, w: po.w}
}

// Len returns the length of a set of pawn options.
func (pos pawnOpts) Len() int {
	return len(pos)
}

// Less compares the weight field of two pawn options in a set.
func (pos pawnOpts) Less(i, j int) bool {
	return pos[i].w < pos[j].w
}

// Swap two pawn options in a set.
func (pos pawnOpts) Swap(i, j int) {
	temp := pos[i]
	pos[i] = pos[j]
	pos[j] = temp
}
