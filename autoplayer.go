package main

import (
	"log"
	"math/rand"
	"sort"
)

// autoPlayer is a set of positions that can be trained. The pawn option selection will be set for each field to indicate the optimal pawn option and the set of pawn options will be optimized and ranked.
type autoPlayer []*position

// train returns an auto player that is capable of playing hexapawn.
func train(m, n, numGames int, t state) autoPlayer {
	ap := make(autoPlayer, 0, 32)
	var (
		g        *game
		poSlc    *pawnOpt
		pawnOpts []*pawnOpt
		psn      *position
		choice   weight
		index    int
	)
	for ; 0 < numGames; numGames-- {
		g = newGame(m, n, cvc)

		// Alternate turns until neither side can move (that is, win, illegal, or stalemate state is reached)
		for {
			psn = &position{b: g.b, s: g.s, poSlc: nil, po: g.availPawnOpts()}
			switch g.s {
			case whiteTurn:
				switch t {
				case whiteTurn:
					ap.insert(psn)
					choice = weight(rand.Float64())
					index = ap.index(psn)
					for _, po := range ap[index].po {
						if choice <= po.p {
							g.move(po.m, po.n, po.a)
							psn.poSlc = po
							break
						}
					}
				case blackTurn:
					pawnOpts = g.availPawnOpts()
					poSlc = pawnOpts[rand.Intn(len(pawnOpts))]
					g.move(poSlc.m, poSlc.n, poSlc.a)
					psn.poSlc = poSlc
				default:
					log.Fatal("turn: invalid state entered")
				}
			case blackTurn:
				switch t {
				case whiteTurn:
					pawnOpts = g.availPawnOpts()
					poSlc = pawnOpts[rand.Intn(len(pawnOpts))]
					g.move(poSlc.m, poSlc.n, poSlc.a)
					psn.poSlc = poSlc
				case blackTurn:
					ap.insert(psn)
					choice = weight(rand.Float64())
					index = ap.index(psn)
					for _, po := range ap[index].po {
						if choice <= po.p {
							g.move(po.m, po.n, po.a)
							psn.poSlc = po
							break
						}
					}
				default:
					log.Fatal("turn: invalid state entered")
				}
			case whiteWin:
				switch t {
				case whiteTurn:
					for i := range g.h {
						index = ap.index(g.h[i])

					}
				case blackTurn:
				}
				break
			case blackWin:
				switch t {
				case whiteTurn:
				case blackTurn:
				}
				break
			case stalemate:
				switch t {
				case whiteTurn:
				case blackTurn:
				}
				break
			case illegal:
				log.Fatal("train: reached illegal state")
				break
			default:
				log.Fatal("train: reached unknown state")
				break
			}

			g.updateState()
			g.h = append(g.h, psn)
		}
	}

	return ap
}

// insert a position into an auto player if it doesn't already exist.
func (ap autoPlayer) insert(p *position) {
	switch len(ap) {
	case 0:
		ap = append(ap, copyPosition(p))
	case ap.index(p):
		ap = append(ap, copyPosition(p))
		sort.Sort(ap)
	}
}

// remove a position from an auto player's experience.
func (ap autoPlayer) remove(i int) *position {
	p := ap[i]
	ap = append(ap[:i], ap[i+1:]...)
	return p
}

// index returns the index a position is found in an auto player. If the position is not found, len(ap) is returned. Comparisions are made on the board and state only.
func (ap autoPlayer) index(p *position) int {
	return sort.Search(len(ap), func(i int) bool { return equalBoards(ap[i].b, p.b) && ap[i].s == p.s })
}

func (ap autoPlayer) Less(i, j int) bool {
	for a := range ap[i].b {
		for b := range ap[i].b[a] {
			if ap[j].b[a][b] < ap[i].b[a][b] {
				return false
			}
		}
	}

	return true
}

func (ap autoPlayer) Len() int {
	return len(ap)
}

func (ap autoPlayer) Swap(i, j int) {
	t := ap[i]
	ap[i] = ap[j]
	ap[j] = t
}
