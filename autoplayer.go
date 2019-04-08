package main

import (
	"log"
	"math/rand"
	"sort"
)

// Action weights are effectively the space between decision boundaries. To select
// an action, choose a random number and determine which boundary it lands in.
// Weights must sum to 1.0.
//
// <-|-|--|----|->
//   0 w0 w1   w2=1

// autoPlayer is a set of positions that can be trained. The pawn option selection
// will be set for each field to indicate the optimal pawn option and the set of
// pawn options will be optimized and ranked.
type autoPlayer struct {
	sd   side        // White or black side
	psns []*position // Set of positions experienced
}

// type (ap *autoPlayer)String()string{
// 	bldr:=strings.Builder{}
// 	for i:=range ap.psns{
// 		bldr.WriteString(ap.psns.String())
// 	}

// 	return bldr.String()
// }

// newAutoPlayer returns an autoPlayer associated with a side.
func newAutoPlayer(sd side) *autoPlayer {
	if sd != whiteSide && sd != blackSide {
		panic("newAutoPlayer: invalid side")
	}

	return &autoPlayer{sd: sd, psns: make([]*position, 0, 32)}
}

// train returns an auto player that is capable of playing hexapawn.
func (ap *autoPlayer) train(m, n, numGames int) {
	reward := weight(0.01)
	var (
		index      int
		apPosLen   int
		punishment weight
		gm         *game
		psn        *position
		gameOver   bool
	)

	for k := 0; k < numGames; k++ {
		gm = newGame(m, n, cvc)
		gameOver = false

		// Alternate turns until neither side can move (that is, win, illegal, or stalemate state is reached)
		for !gameOver {
			psn = &position{
				brd: gm.brd,
				st:  gm.st,
				pos: availPawnOpts(gm.brd, gm.st),
			}

			switch gm.st {
			case whiteTurn:
				switch ap.sd {
				case whiteSide:
					gm.move(ap.move(psn))
				case blackSide:
					gm.move(&event{psn: copyPosition(psn), poSlc: randPawnOpt(psn)})
				}
			case blackTurn:
				switch ap.sd {
				case whiteSide:
					gm.move(&event{psn: copyPosition(psn), poSlc: randPawnOpt(psn)})
				case blackSide:
					gm.move(ap.move(psn))
				}
			case whiteWin:
				switch ap.sd {
				case whiteSide:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						if index == len(ap.psns) {
							continue
						}

						apPosLen = len(ap.psns[index].pos)
						if apPosLen < 2 {
							continue // Either zero or one pawn option to select; nothing to train on
						}

						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap.psns[index].pos {
							if appo.act == hstpsn.poSlc.act {
								ap.psns[index].pos[i].wght += reward
								continue
							}

							ap.psns[index].pos[i].wght -= punishment
						}
					}
				case blackSide:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						if index == len(ap.psns) {
							continue
						}

						apPosLen = len(ap.psns[index].pos)
						if apPosLen < 2 {
							continue // Either zero or one pawn option to select; nothing to train on
						}

						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap.psns[index].pos {
							if appo.act == hstpsn.poSlc.act {
								ap.psns[index].pos[i].wght -= reward
								continue
							}

							ap.psns[index].pos[i].wght += punishment
						}
					}
				}

				gameOver = true
			case blackWin:
				switch ap.sd {
				case whiteSide:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						if index == len(ap.psns) {
							continue
						}

						apPosLen = len(ap.psns[index].pos)
						if apPosLen < 2 {
							continue // Either zero or one pawn option to select; nothing to train on
						}

						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap.psns[index].pos {
							if appo.act == hstpsn.poSlc.act {
								ap.psns[index].pos[i].wght -= reward
								continue
							}

							ap.psns[index].pos[i].wght += punishment
						}
					}
				case blackSide:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						if index == len(ap.psns) {
							continue
						}

						apPosLen = len(ap.psns[index].pos)
						if apPosLen < 2 {
							continue // Either zero or one pawn option to select; nothing to train on
						}

						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap.psns[index].pos {
							if appo.act == hstpsn.poSlc.act {
								ap.psns[index].pos[i].wght += reward
								continue
							}

							ap.psns[index].pos[i].wght -= punishment
						}
					}
				}

				gameOver = true
			case stalemate:
				for _, hstpsn := range gm.hst {
					index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
					if index == len(ap.psns) {
						continue
					}

					apPosLen = len(ap.psns[index].pos)
					if apPosLen < 2 {
						continue // Either zero or one pawn option to select; nothing to train on
					}

					punishment = reward / weight(apPosLen-1)
					for i, appo := range ap.psns[index].pos {
						if appo.act == hstpsn.poSlc.act {
							ap.psns[index].pos[i].wght -= reward
							continue
						}

						ap.psns[index].pos[i].wght += punishment
					}
				}

				gameOver = true
			case illegal:
				log.Fatal("train: reached illegal state")
			default:
				log.Fatal("train: reached unknown state")
			}
		}
	}
}

// randPawnOpt returns a random pawn option (nil if none available).
func randPawnOpt(psn *position) *pawnOpt {
	n := len(psn.pos)
	if n == 0 {
		return nil
	}

	return psn.pos[rand.Intn(n)]
}

// move returns an event representing an action taken on a given position. An event
// with no pawn option selected is returned if a position has no available pawn
// options.
func (ap *autoPlayer) move(psn *position) *event {
	index := ap.index(psn.brd, psn.st)
	if index == len(ap.psns) {
		index = ap.insert(psn.brd, psn.st)
	}

	choice := weight(rand.Float64())
	var sum weight
	for _, po := range ap.psns[index].pos {
		sum += po.wght
		if choice <= sum {
			return &event{psn: copyPosition(psn), poSlc: copyPawnOpt(po)}
		}
	}

	return &event{psn: copyPosition(psn)} // TODO: determine if this should panic here
}

// insert a position into an auto player and returns the position it is found in
// after sorting.
func (ap *autoPlayer) insert(brd board, st state) int {
	ap.psns = append(ap.psns, &position{brd: copyBoard(brd), st: st, pos: availPawnOpts(brd, st)})
	sort.Slice(ap.psns, ap.less)
	return ap.index(brd, st)
}

// remove a position from an auto player's experience.
func (ap *autoPlayer) remove(i int) *position {
	psn := ap.psns[i]
	ap.psns = append(ap.psns[:i], ap.psns[i+1:]...)
	return psn
}

// index returns the index a position is found in an auto player. If the position
// is not found, len(ap.psns) is returned.
func (ap *autoPlayer) index(brd board, st state) int {
	return sort.Search(len(ap.psns), func(i int) bool { return equalBoards(ap.psns[i].brd, brd) && ap.psns[i].st == st })
}

// Less returns true if each less-than pawn comparison in two boards is true and
// false if otherwise.
func (ap *autoPlayer) less(i, j int) bool {
	for a := range ap.psns[i].brd {
		for b := range ap.psns[i].brd[a] {
			if ap.psns[j].brd[a][b] <= ap.psns[i].brd[a][b] {
				return false
			}
		}
	}

	return true
}
