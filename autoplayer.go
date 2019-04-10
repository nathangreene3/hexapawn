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

	bldr.Write([]byte{byte(ap.sd), '\n'})
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

// train returns an auto player that is capable of playing hexapawn.
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
						if index == len(ap.psns) {
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
						if index == len(ap.psns) {
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
						if index == len(ap.psns) {
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
						if index == len(ap.psns) {
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
					if index == len(ap.psns) {
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

// randPawnOpt returns a random pawn option (nil if none available).
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
	if index == len(ap.psns) {
		index = ap.insert(psn)
	}
	fmt.Printf("psn.pos:\n%s\n", psn)
	fmt.Printf("ap.psns[%d]:\n%s\n", index, ap.psns[index])

	choice := weight(rand.Float64())
	var sum weight
	for _, po := range ap.psns[index].pos {
		sum += po.wght
		if choice <= sum {
			fmt.Printf("po: %v\n", po)
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
// is not found, len(ap.psns) is returned.
func (ap *autoPlayer) index(psn *position) int {
	// for i := range ap.psns {
	// 	if equalPositions(ap.psns[i], psn) {
	// 		return i
	// 	}
	// }

	// return len(ap.psns)

	// return ap.search(psn, 0, len(ap.psns)-1)

	return sort.Search(len(ap.psns), func(i int) bool { return lessEqPositions(psn, ap.psns[i]) })
}

// Less returns true if each less-than pawn comparison in two boards is true and
// false if otherwise.
func (ap *autoPlayer) less(i, j int) bool {
	if comparePositions(ap.psns[i], ap.psns[j]) < 0 {
		return true
	}

	return false
}

func (ap *autoPlayer) lessEq(i, j int) bool {
	if 0 < comparePositions(ap.psns[i], ap.psns[j]) {
		return false
	}

	return true
}

func (ap *autoPlayer) isSorted() bool {
	n := len(ap.psns) - 1
	for i := 0; i < n; i++ {
		if 0 < comparePositions(ap.psns[i], ap.psns[i+1]) {
			return false
		}
	}

	return true
}

func (ap *autoPlayer) search(psn *position, i, j int) int {
	var k int

	// THIS WORKS
	// for i <= j {
	// 	k = i + int(uint(j-i)>>1)
	// 	switch comparePositions(psn, ap.psns[k]) {
	// 	case -1:
	// 		j = k - 1
	// 	case 1:
	// 		i = k + 1
	// 	default:
	// 		return k
	// 	}
	// }

	// THIS DOESN'T WORK
	for i < j {
		k = i + int(uint(j-i)>>1)
		if 0 < comparePositions(psn, ap.psns[k]) {
			i = k + 1 // ap.psns[k] < psn
			continue
		}

		j = k // psn <= ap.psns[k]
	}

	return i
}
