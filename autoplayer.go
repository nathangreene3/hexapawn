package main

import (
	"log"
	"math/rand"
	"sort"
)

// autoPlayer is a set of positions that can be trained. The pawn option selection
// will be set for each field to indicate the optimal pawn option and the set of
// pawn options will be optimized and ranked.
type autoPlayer []*position

// train returns an auto player that is capable of playing hexapawn.
func train(m, n, numGames int, t state) autoPlayer {
	ap := make(autoPlayer, 0, 32)
	reward := weight(0.01)
	var (
		gm         *game
		poSlc      *pawnOpt
		pawnOpts   []*pawnOpt
		psn        *position
		choice     weight
		sumWeights weight
		index      int
		evnt       *event
		punishment weight
		apPosLen   int
	)

	for ; 0 < numGames; numGames-- {
		gm = newGame(m, n, cvc)

		// Alternate turns until neither side can move (that is, win, illegal, or stalemate state is reached)
		for {
			psn = &position{brd: gm.brd, st: gm.st, pos: availPawnOpts(gm.brd, gm.st)}
			switch gm.st {
			case whiteTurn:
				switch t {
				case whiteTurn:
					index = ap.index(gm.brd, gm.st)
					if index == len(ap) {
						index = ap.insert(gm.brd, gm.st)
					}

					choice = weight(rand.Float64())
					for _, po := range ap[index].pos {
						sumWeights += po.wght
						if choice <= sumWeights {
							gm.move(po.m, po.n, po.act)
							evnt = &event{psn: psn, poSlc: copyPawnOpt(po)}
							break
						}
					}

					sumWeights = 0.0
				case blackTurn:
					pawnOpts = availPawnOpts(gm.brd, gm.st)
					poSlc = pawnOpts[rand.Intn(len(pawnOpts))]
					gm.move(poSlc.m, poSlc.n, poSlc.act)
				default:
					log.Fatal("turn: invalid state entered")
				}
			case blackTurn:
				switch t {
				case whiteTurn:
					pawnOpts = availPawnOpts(gm.brd, gm.st)
					poSlc = pawnOpts[rand.Intn(len(pawnOpts))]
					gm.move(poSlc.m, poSlc.n, poSlc.act)
				case blackTurn:
					index = ap.index(gm.brd, gm.st)
					if index == len(ap) {
						index = ap.insert(gm.brd, gm.st)
					}

					choice = weight(rand.Float64())
					for _, po := range ap[index].pos {
						sumWeights += po.wght
						if choice <= sumWeights {
							gm.move(po.m, po.n, po.act)
							evnt = &event{psn: psn, poSlc: copyPawnOpt(po)}
							break
						}
					}

					sumWeights = 0.0
				default:
					log.Fatal("turn: invalid state entered")
				}
			case whiteWin:
				switch t {
				case whiteTurn:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						apPosLen = len(ap[index].pos)
						switch apPosLen {
						case 0:
							log.Fatal("train: cannot train on zero pawn options")
						case 1:
							continue // Only one pawn option to select
						default:
							punishment = reward / weight(apPosLen-1)
							for i, appo := range ap[index].pos {
								if appo.act == hstpsn.poSlc.act {
									ap[index].pos[i].wght += reward
									continue
								}

								ap[index].pos[i].wght -= punishment
							}
						}
					}
				case blackTurn:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						apPosLen = len(ap[index].pos)
						switch apPosLen {
						case 0:
							log.Fatal("train: cannot train on zero pawn options")
						case 1:
							continue // Only one pawn option to select
						default:
							punishment = reward / weight(apPosLen-1)
							for i, appo := range ap[index].pos {
								if appo.act == hstpsn.poSlc.act {
									ap[index].pos[i].wght -= reward
									continue
								}

								ap[index].pos[i].wght += punishment
							}
						}
					}
				}
				break
			case blackWin:
				switch t {
				case whiteTurn:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						apPosLen = len(ap[index].pos)
						switch apPosLen {
						case 0:
							log.Fatal("train: cannot train on zero pawn options")
						case 1:
							continue // Only one pawn option to select
						default:
							punishment = reward / weight(apPosLen-1)
							for i, appo := range ap[index].pos {
								if appo.act == hstpsn.poSlc.act {
									ap[index].pos[i].wght -= reward
									continue
								}

								ap[index].pos[i].wght += punishment
							}
						}
					}
				case blackTurn:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						apPosLen = len(ap[index].pos)
						switch apPosLen {
						case 0:
							log.Fatal("train: cannot train on zero pawn options")
						case 1:
							continue // Only one pawn option to select
						default:
							punishment = reward / weight(apPosLen-1)
							for i, appo := range ap[index].pos {
								if appo.act == hstpsn.poSlc.act {
									ap[index].pos[i].wght += reward
									continue
								}

								ap[index].pos[i].wght -= punishment
							}
						}
					}
				}
				break
			case stalemate:
				for _, hstpsn := range gm.hst {
					index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
					apPosLen = len(ap[index].pos)
					switch apPosLen {
					case 0:
						log.Fatal("train: cannot train on zero pawn options")
					case 1:
						continue // Only one pawn option to select
					default:
						punishment = reward / weight(apPosLen-1)
						for i, appo := range ap[index].pos {
							if appo.act == hstpsn.poSlc.act {
								ap[index].pos[i].wght -= reward
								continue
							}

							ap[index].pos[i].wght += punishment
						}
					}
				}
				break
			case illegal:
				log.Fatal("train: reached illegal state")
				break
			default:
				log.Fatal("train: reached unknown state")
				break
			}

			gm.hst = append(gm.hst, evnt)
			gm.checkWin()
		}
	}

	return ap
}

// insert a position into an auto player and returns the position it is found in after sorting.
func (ap autoPlayer) insert(brd board, st state) int {
	ap = append(ap, &position{brd: copyBoard(brd), st: st, pos: availPawnOpts(brd, st)})
	sort.Slice(ap, ap.less)
	return ap.index(brd, st)
}

// remove a position from an auto player's experience.
func (ap autoPlayer) remove(i int) *position {
	psn := ap[i]
	ap = append(ap[:i], ap[i+1:]...)
	return psn
}

// index returns the index a position is found in an auto player. If the position
// is not found, len(ap) is returned.
func (ap autoPlayer) index(brd board, st state) int {
	return sort.Search(len(ap), func(i int) bool { return equalBoards(ap[i].brd, brd) && ap[i].st == st })
}

// Less returns true if each less-than pawn comparison in two boards is true and
// false if otherwise.
func (ap autoPlayer) less(i, j int) bool {
	for a := range ap[i].brd {
		for b := range ap[i].brd[a] {
			if ap[j].brd[a][b] < ap[i].brd[a][b] {
				return false
			}
		}
	}

	return true
}
