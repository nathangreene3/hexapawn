package main

import (
	"bytes"
	"strings"
)

// action represents what a player can do at a position (m,n) on a board.
type action byte

// mode indicates player vs player, player vs an auto player, or auto player vs auto player.
type mode byte

// pawn is a playable piece or a blank space.
type pawn byte

// state indicates how the game will procede.
type state byte

// game joins a board, state, mode, and a history of board positions reached in alternating turns.
type game struct {
	brd board   // Current board
	st  state   // Current state
	md  mode    // Type to play
	hst history // Ordered set of board positions
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

	// States
	illegal state = iota
	stalemate
	whiteTurn
	blackTurn
	whiteWin
	blackWin
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
func (gm *game) play() {
	// for gm.st == whiteTurn || gm.st == blackTurn {
	// 	gm.turn()
	// 	gm.checkWin()
	// 	// g.h = append(g.h, copyBoard(g.b))
	// }

	// switch gm.st {
	// case whiteWin:
	// 	fmt.Println("WHITE WINS")
	// case blackWin:
	// 	fmt.Println("BLACK WINS")
	// case illegal:
	// 	fmt.Println("ILLEGAL STATE")
	// case stalemate:
	// 	fmt.Println("STALEMATE")
	// }
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

// move performs an action altering the position of the board. State is NOT altered.
func (gm *game) move(m, n int, a action) {
	switch gm.brd[m][n] {
	case whitePawn:
		switch a {
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
	case blackPawn:
		switch a {
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
	}
}

// availActions returns a set of actions that can be taken at a position (m,n).
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

// checkWin checks the board for a win condition and sets the game state to the winning state. If the game has not been won, then it swaps the turn
func (gm *game) checkWin() {
	switch gm.st {
	case whiteTurn:
		gm.st = blackTurn
		for i := range gm.brd[0] {
			if gm.brd[0][i] == whitePawn {
				gm.st = whiteWin
				break
			}
		}
	case blackTurn:
		gm.st = whiteTurn
		n := len(gm.brd[0]) - 1
		for i := range gm.brd[n] {
			if gm.brd[n][i] == blackPawn {
				gm.st = blackWin
				break
			}
		}
	}
}
