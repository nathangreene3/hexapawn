package main

import "sort"

// weight is a probability value on the range [0,1].
type weight float64

// pawnOpts is a set of pawn options.
type pawnOpts []*pawnOpt

// pawnOpt is an available action at a position (m,n) with a probability weight of being selected.
type pawnOpt struct {
	m int    // Row index
	n int    // Column index
	a action // Available action
	p weight // TODO: redefine this field to something more appropriate
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

// Len returns the length of a set of pawn options.
func (pos pawnOpts) Len() int {
	return len(pos)
}

// Less compares the weight field of two pawn options in a set.
func (pos pawnOpts) Less(i, j int) bool {
	return pos[i].p < pos[j].p
}

// Swap two pawn options in a set.
func (pos pawnOpts) Swap(i, j int) {
	t := pos[i]
	pos[i] = pos[j]
	pos[j] = t
}

// equalPawnOpts returns true if each pawn option field is equal.
func equalPawnOpts(p, q *pawnOpt) bool {
	switch {
	case p == nil:
		return q == nil
	case q == nil:
		return p == nil
	case p.m != q.m:
		return false
	case p.n != q.n:
		return false
	case p.a != q.a:
		return false
	case p.p != q.p:
		return false
	default:
		return true
	}
}
