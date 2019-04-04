package main

import (
	"log"
	"math/rand"
	"sort"
)

// autoPlayer is a set of positions that can be trained. The pawn option selection
// will be set for each field to indicate the optimal pawn option and the set of
// pawn options will be optimized and ranked.
type autoPlayer struct {
	sd   side
	psns []*position
}

// train returns an auto player that is capable of playing hexapawn.
func train(m, n, numGames int, sd side) *autoPlayer {
	if sd != whiteSide && sd != blackSide {
		panic("train: side must be white or black")
	}

	ap := &autoPlayer{sd: sd, psns: make([]*position, 0, 32)}
	reward := weight(0.01)
	var (
		gm         *game
		poSlc      *pawnOpt
		pawnOpts   []*pawnOpt
		psn        *position
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
				switch sd {
				case whiteSide:
					evnt = ap.move(psn)
					gm.move(evnt.poSlc.m, evnt.poSlc.n, evnt.poSlc.act)
				case blackSide:
					pawnOpts = availPawnOpts(gm.brd, gm.st)
					poSlc = pawnOpts[rand.Intn(len(pawnOpts))]
					gm.move(poSlc.m, poSlc.n, poSlc.act)
				default:
					log.Fatal("turn: invalid state entered")
				}
			case blackTurn:
				switch sd {
				case whiteSide:
					pawnOpts = availPawnOpts(gm.brd, gm.st)
					poSlc = pawnOpts[rand.Intn(len(pawnOpts))]
					gm.move(poSlc.m, poSlc.n, poSlc.act)
				case blackSide:
					evnt = ap.move(psn)
					gm.move(evnt.poSlc.m, evnt.poSlc.n, evnt.poSlc.act)
				default:
					panic("turn: invalid state entered")
				}
			case whiteWin:
				switch sd {
				case whiteSide:
					for _, hstpsn := range gm.hst {
						index = ap.index(hstpsn.psn.brd, hstpsn.psn.st)
						if index == len(ap.psns) {
							continue
						}

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

	return &ap
}

func (ap *autoPlayer) move(psn *position) *event {
	index := ap.index(psn.brd, psn.st)
	if index == len(ap.psns) {
		index = ap.insert(psn.brd, psn.st)
	}

	choice := weight(rand.Float64())
	var sum weight
	// fmt.Println(index, ap, len(ap.psns))
	for _, po := range ap.psns[index].pos {
		sum += po.wght
		if choice <= sum {
			return &event{psn: copyPosition(psn), poSlc: copyPawnOpt(po)}
		}
	}

	return &event{psn: copyPosition(psn)}
}

// insert a position into an auto player and returns the position it is found in after sorting.
func (ap *autoPlayer) insert(brd board, st state) int {
	*ap = append(*ap, &position{brd: copyBoard(brd), st: st, pos: availPawnOpts(brd, st)})
	sort.Slice(*ap, ap.less)
	return ap.index(brd, st)
}

// remove a position from an auto player's experience.
func (ap *autoPlayer) remove(i int) *position {
	psn := ap.psns[i]
	*ap = append(ap.psns[:i], ap.psns[i+1:]...)
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
			if ap.psns[j].brd[a][b] < ap.psns[i].brd[a][b] {
				return false
			}
		}
	}

	return true
}
