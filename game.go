package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

// action represents what a player can do at a position (m,n) on a board.
type action byte

// mode indicates player vs player, player vs an auto player, or auto player vs auto player.
type mode byte

// pawn is a playable piece or a blank space.
type pawn byte

// side is either white or black.
type side byte

// state indicates how the game will procede.
type state byte

// game joins a board, state, mode, and a history of events reached in alternating
// turns.
type game struct {
	brd board   // Current board
	st  state   // Current state
	md  mode    // Type of game to play
	hst history // Ordered set of events
}

// Game constants
const (
	// Actions
	forward action = iota
	captureLeft
	captureRight

	// Modes
	pvp mode = iota // Player vs player
	pvc             // Player vs computer
	cvp             // Computer vs player
	cvc             // Computer vs computer

	// Pawns
	space     = pawn(' ')
	whitePawn = pawn('w')
	blackPawn = pawn('b')

	// Sides
	whiteSide = side('w')
	blackSide = side('b')

	// States
	illegal state = iota
	whiteTurn
	blackTurn
	whiteWin
	blackWin
	stalemate
)

// String returns a string representing a game.
func (gm *game) String() string {
	n := len(gm.brd[0])
	bldr := strings.Builder{}
	bldr.Grow((2*len(gm.brd) + 1) * (2*n + 1))

	dashPlus := []byte{'-', '+'}
	barnl := []byte{'|', '\n'}
	line := bytes.Join([][]byte{[]byte{'+'}, bytes.Repeat(dashPlus, n), []byte{'\n'}}, []byte{})

	for i := range gm.brd {
		bldr.Write(line)

		for j := range gm.brd[i] {
			bldr.WriteByte('|')
			bldr.WriteByte(byte(gm.brd[i][j]))
		}

		bldr.Write(barnl)
	}

	bldr.WriteByte('+')
	bldr.Write(bytes.Repeat(dashPlus, n))

	return bldr.String()
}

// newGame returns a game to be played.
func newGame(m, n int, md mode) *game {
	return &game{
		brd: newBoard(m, n),
		st:  whiteTurn,
		md:  md,
		hst: make(history, 0, 32), // TODO: Find the total number of legal states (24?)
	}
}

// play
func play(m, n int, md mode) {
	gm := newGame(m, n, md)

	var psn *position

	switch md {
	case cvc:
		white := newAutoPlayer(whiteSide)
		black := newAutoPlayer(blackSide)
		white.train(m, n, 10)
		black.train(m, n, 10)
		var gameOver bool

		for !gameOver {
			psn = &position{brd: gm.brd, st: gm.st, pos: availPawnOpts(gm.brd, gm.st)}
			fmt.Printf("%s\n\n", gm.String())

			switch gm.st {
			case whiteTurn:
				gm.move(white.move(psn))
			case blackTurn:
				gm.move(black.move(psn))
			default:
				gameOver = true
			}
		}
	case cvp: // TODO
	case pvc: // TODO
	case pvp: // TODO
	default:
		log.Fatal("play: invalid mode")
	}

	switch gm.st {
	case whiteWin:
		fmt.Println("WHITE WINS")
	case blackWin:
		fmt.Println("BLACK WINS")
	case illegal:
		fmt.Println("ILLEGAL STATE")
	case stalemate:
		fmt.Println("STALEMATE")
	}
}

// turn
func (gm *game) turn() {
	switch gm.st {
	case whiteTurn:
		switch gm.md {
		case pvp, pvc:
			// player move
			// return move result
		case cvp, cvc:
			// random move by npc
			// return move result
		}
	case blackTurn:
		switch gm.md {
		case pvp, cvp:
			// player move
			// return move result
		case pvc, cvc:
			// random move by npc
			// return move result
		}
	}
}

// move performs an action altering the position of the board.
func (gm *game) move(evnt *event) {
	if evnt.poSlc != nil {
		m, n := evnt.poSlc.m, evnt.poSlc.n
		act := evnt.poSlc.act

		switch gm.brd[m][n] {
		case whitePawn:
			switch act {
			case forward:
				if 0 < m && gm.brd[m-1][n] == space {
					gm.brd[m-1][n] = whitePawn
					gm.brd[m][n] = space
				}
			case captureLeft:
				if 0 < m && 0 < n && gm.brd[m-1][n-1] == blackPawn {
					gm.brd[m-1][n-1] = whitePawn
					gm.brd[m][n] = space
				}
			case captureRight:
				if 0 < m && n+1 < len(gm.brd[0]) && gm.brd[m-1][n+1] == blackPawn {
					gm.brd[m-1][n+1] = whitePawn
					gm.brd[m][n] = space
				}
			}

			if checkWin(gm.brd, gm.st) {
				gm.st = whiteWin
			} else {
				gm.st = blackTurn
			}
		case blackPawn:
			switch act {
			case forward:
				if m+1 < len(gm.brd) && gm.brd[m][n] == space {
					gm.brd[m+1][n] = blackPawn
					gm.brd[m][n] = space
				}
			case captureLeft:
				if m+1 < len(gm.brd) && n+1 < len(gm.brd[0]) && gm.brd[m+1][n+1] == whitePawn {
					gm.brd[m+1][n+1] = blackPawn
					gm.brd[m][n] = space
				}
			case captureRight:
				if m+1 < len(gm.brd) && n-1 < len(gm.brd[0]) && gm.brd[m+1][n-1] == whitePawn {
					gm.brd[m+1][n-1] = blackPawn
					gm.brd[m][n] = space
				}
			}

			if checkWin(gm.brd, gm.st) {
				gm.st = blackWin
			} else {
				gm.st = whiteTurn
			}
		default:
			panic("move: cannot move space")
		}
	} else {
		gm.st = stalemate
	}

	gm.hst = append(gm.hst, evnt)
}

// availActions returns a set of actions that can be taken at a position (m,n).
// Actions are available if the state is either white or black turn.
func availActions(m, n int, brd board, st state) []action {
	acts := make([]action, 0, 4)
	lenB := len(brd)
	lenB0m1 := len(brd[0]) - 1

	switch st {
	case whiteTurn:
		if brd[m][n] == whitePawn && 0 < m {
			if brd[m-1][n] == space {
				acts = append(acts, forward)
			}

			switch n {
			case 0:
				if brd[m-1][n+1] == blackPawn {
					acts = append(acts, captureRight)
				}
			case lenB0m1:
				if brd[m-1][n-1] == blackPawn {
					acts = append(acts, captureLeft)
				}
			default:
				if brd[m-1][n-1] == blackPawn {
					acts = append(acts, captureLeft)
				}

				if brd[m-1][n+1] == blackPawn {
					acts = append(acts, captureRight)
				}
			}
		}
	case blackTurn:
		if brd[m][n] == blackPawn && m+1 < lenB {
			if brd[m+1][n] == space {
				acts = append(acts, forward)
			}

			switch n {
			case 0:
				if brd[m+1][n+1] == whitePawn {
					acts = append(acts, captureLeft)
				}
			case lenB0m1:
				if brd[m+1][n-1] == whitePawn {
					acts = append(acts, captureRight)
				}
			default:
				if brd[m+1][n-1] == whitePawn {
					acts = append(acts, captureRight)
				}

				if brd[m+1][n+1] == whitePawn {
					acts = append(acts, captureLeft)
				}
			}
		}
	}

	return acts
}

// checkWin checks the board for a win condition given a state. If the state is
// neither white nor black turn, then false is returned.
func checkWin(brd board, st state) bool {
	var count int
	switch st {
	case whiteTurn:
		for i := range brd[0] {
			if brd[0][i] == whitePawn {
				return true
			}

			count++
		}
		for i := 0; i < len(brd)-1; i++ {
			for _, p := range brd[i] {
				if p == blackPawn {
					count++
				}
			}
		}

		if count == 0 {
			return true // No pieces left for black to move
		}

		return false
	case blackTurn:
		n := len(brd[0]) - 1
		for i := range brd[n] {
			if brd[n][i] == blackPawn {
				return true
			}

			count++
		}

		for i := 1; i < len(brd); i++ {
			for _, p := range brd[i] {
				if p == whitePawn {
					count++
				}
			}
		}

		if count == 0 {
			return true // No pieces left for white to move
		}

		return false
	default:
		return false
	}
}
