package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
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

func (ap *autoPlayer) String() string {
	bldr := strings.Builder{}

	bldr.WriteString(fmt.Sprintf("side: %q\n", byte(ap.sd)))
	for i := range ap.psns {
		bldr.WriteString(ap.psns[i].String())
	}

	return bldr.String()
}

// newAutoPlayer returns an autoPlayer associated with a side.
func newAutoPlayer(sd side) *autoPlayer {
	if sd != whiteSide && sd != blackSide {
		panic("newAutoPlayer: invalid side")
	}

	return &autoPlayer{sd: sd, psns: make([]*position, 0, 32)}
}

// train an auto player on a number of random games.
func (ap *autoPlayer) train(m, n, numGames int) {
	var (
		gameOver   bool          // Indicates game was won/lost
		index      int           // Index of position in auto player
		apPosLen   int           // Number of pawn options in the indexed position of auto player
		punishment weight        // Amount to alter non-selected pawn options' weights
		reward     = weight(0.1) // Amount to alter selected pawn option weight
		gm         *game         // Game to be played for a given number of games
		psn        *position     // Current position of the game
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
						index = ap.index(hstpsn.psn)
						if index < 0 {
							continue
						}

						apPosLen = len(ap.psns[index].pos)
						if apPosLen < 2 {
							continue // Either zero or one pawn option to select; nothing to train on
						}

						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap.psns[index].pos {
							if equalPawnOpts(appo, hstpsn.poSlc) {
								ap.psns[index].pos[i].wght += reward
								continue
							}

							ap.psns[index].pos[i].wght -= punishment
						}
					}
				case blackSide:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn)
						if index < 0 {
							continue
						}

						apPosLen = len(ap.psns[index].pos)
						if apPosLen < 2 {
							continue // Either zero or one pawn option to select; nothing to train on
						}

						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap.psns[index].pos {
							if equalPawnOpts(appo, hstpsn.poSlc) {
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
						index = ap.index(hstpsn.psn)
						if index < 0 {
							continue
						}

						apPosLen = len(ap.psns[index].pos)
						if apPosLen < 2 {
							continue // Either zero or one pawn option to select; nothing to train on
						}

						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap.psns[index].pos {
							if equalPawnOpts(appo, hstpsn.poSlc) {
								ap.psns[index].pos[i].wght -= reward
								continue
							}

							ap.psns[index].pos[i].wght += punishment
						}
					}
				case blackSide:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn)
						if index < 0 {
							continue
						}

						apPosLen = len(ap.psns[index].pos)
						if apPosLen < 2 {
							continue // Either zero or one pawn option to select; nothing to train on
						}

						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap.psns[index].pos {
							if equalPawnOpts(appo, hstpsn.poSlc) {
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
					index = ap.index(hstpsn.psn)
					if index < 0 {
						continue
					}

					apPosLen = len(ap.psns[index].pos)
					if hstpsn.poSlc == nil || apPosLen < 2 {
						continue // Either zero (stalemate) or one pawn option to select; nothing to train on
					}

					punishment = reward / weight(apPosLen-1)
					for i, appo := range ap.psns[index].pos {
						if equalPawnOpts(appo, hstpsn.poSlc) {
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

// randPawnOpt returns a random pawn option at a given position (nil if none available).
func randPawnOpt(psn *position) *pawnOpt {
	n := len(psn.pos)
	if 0 < n {
		return psn.pos[rand.Intn(n)]
	}

	return nil
}

// move returns an event representing an action taken on a given position. An event
// with no pawn option selected is returned if a position has no available pawn
// options.
func (ap *autoPlayer) move(psn *position) *event {
	index := ap.index(psn)
	if index < 0 {
		index = ap.insert(psn)
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
func (ap *autoPlayer) insert(psn *position) int {
	ap.psns = append(ap.psns, copyPosition(psn))
	sort.SliceStable(ap.psns, ap.less)
	return ap.index(psn)
}

// remove a position from an auto player's experience.
func (ap *autoPlayer) remove(i int) *position {
	psn := ap.psns[i]
	ap.psns = append(ap.psns[:i], ap.psns[i+1:]...)
	return psn
}

// index returns the index a position is found in an auto player. If the position
// is not found, -1 is returned.
func (ap *autoPlayer) index(psn *position) int {
	n := len(ap.psns)
	index := sort.Search(n, func(i int) bool { return lessEqPositions(psn, ap.psns[i]) })
	if index < n && equalPositions(psn, ap.psns[index]) {
		return index
	}

	return -1
}

// less compares two indexed positions.
func (ap *autoPlayer) less(i, j int) bool {
	if comparePositions(ap.psns[i], ap.psns[j]) < 0 {
		return true
	}

	return false
}

// lessEq compares two indexed positions.
func (ap *autoPlayer) lessEq(i, j int) bool {
	if 0 < comparePositions(ap.psns[i], ap.psns[j]) {
		return false
	}

	return true
}
