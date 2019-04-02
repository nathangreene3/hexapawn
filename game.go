package main

import (
	"bytes"
	"fmt"
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

// board is an m-by-n array of pieces.
type board [][]pawn

// game joins a board, state, mode, and a history of board positions reached in alternating turns.
type game struct {
	b board   // Current board
	s state   // Current state
	m mode    // Type to play
	h history // Ordered set of board positions
}

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
func (g *game) String() string {
	n := len(g.b[0])
	bldr := strings.Builder{}
	bldr.Grow((2*len(g.b) + 1) * (2*n + 1))

	dashPlus := []byte{'-', '+'}
	barnl := []byte{'|', '\n'}
	line := bytes.Join([][]byte{[]byte{'+'}, bytes.Repeat(dashPlus, n), []byte{'\n'}}, []byte{})

	for i := range g.b {
		bldr.Write(line)

		for j := range g.b[i] {
			bldr.WriteByte('|')
			bldr.WriteByte(byte(g.b[i][j]))
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
		b: newBoard(m, n),
		s: whiteTurn,
		m: md,
		h: make(history, 0, 32), // TODO: Find the total number of legal states (24?)
	}
}

// play
func (g *game) play() {
	// g.h = append(g.h, copyBoard(g.b))

	for g.s == whiteTurn || g.s == blackTurn {
		g.turn()
		g.updateState()
		// g.h = append(g.h, copyBoard(g.b))
	}

	switch g.s {
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
func (g *game) turn() {
	switch g.s {
	case whiteTurn:
		switch g.m {
		case pvp, pvc:
			// player move
			// return move result
		case cvp, cvc:
			// random move by npc
			// return move result
		}
	case blackTurn:
		switch g.m {
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
func (g *game) move(m, n int, a action) {
	switch g.b[m][n] {
	case whitePawn:
		switch a {
		case forward:
			if 0 < m && g.b[m-1][n] == space {
				g.b[m-1][n] = whitePawn
				g.b[m][n] = space
			}
		case captureLeft:
			if 0 < m && 0 < n && g.b[m-1][n-1] == blackPawn {
				g.b[m-1][n-1] = whitePawn
				g.b[m][n] = space
			}
		case captureRight:
			if 0 < m && n+1 < len(g.b[0]) && g.b[m-1][n+1] == blackPawn {
				g.b[m-1][n+1] = whitePawn
				g.b[m][n] = space
			}
		}
	case blackPawn:
		switch a {
		case forward:
			if m+1 < len(g.b) && g.b[m][n] == space {
				g.b[m+1][n] = blackPawn
				g.b[m][n] = space
			}
		case captureLeft:
			if m+1 < len(g.b) && n+1 < len(g.b[0]) && g.b[m+1][n+1] == whitePawn {
				g.b[m+1][n+1] = blackPawn
				g.b[m][n] = space
			}
		case captureRight:
			if m+1 < len(g.b) && n-1 < len(g.b[0]) && g.b[m+1][n-1] == whitePawn {
				g.b[m+1][n-1] = blackPawn
				g.b[m][n] = space
			}
		}
	}
}

// availActions returns a set of actions that can be taken at a position (m,n).
func (g *game) availActions(m, n int) []action {
	a := make([]action, 0, 4)
	lenB := len(g.b)
	lenB0m1 := len(g.b[0]) - 1

	switch g.s {
	case whiteTurn:
		if g.b[m][n] == whitePawn && 0 < m {
			if g.b[m-1][n] == space {
				a = append(a, forward)
			}

			switch n {
			case 0:
				if g.b[m-1][n+1] == blackPawn {
					a = append(a, captureRight)
				}
			case lenB0m1:
				if g.b[m-1][n-1] == blackPawn {
					a = append(a, captureLeft)
				}
			default:
				if g.b[m-1][n-1] == blackPawn {
					a = append(a, captureLeft)
				}

				if g.b[m-1][n+1] == blackPawn {
					a = append(a, captureRight)
				}
			}
		}
	case blackTurn:
		if g.b[m][n] == blackPawn && m+1 < lenB {
			if g.b[m+1][n] == space {
				a = append(a, forward)
			}

			switch n {
			case 0:
				if g.b[m+1][n+1] == whitePawn {
					a = append(a, captureLeft)
				}
			case lenB0m1:
				if g.b[m+1][n-1] == whitePawn {
					a = append(a, captureRight)
				}
			default:
				if g.b[m+1][n-1] == whitePawn {
					a = append(a, captureRight)
				}

				if g.b[m+1][n+1] == whitePawn {
					a = append(a, captureLeft)
				}
			}
		}
	}

	return a
}

// availPawnOpts returns a set of pawn options available at a position (m,n).
func (g *game) availPawnOpts() []*pawnOpt {
	po := make([]*pawnOpt, 0, 4)
	var a []action // Set of actions for each position (i,j)
	var w weight   // Weight to apply to each action
	var d weight   // Difference in each action weight
	var n int      // Number of available actions

	for i := range g.b {
		for j := range g.b[i] {
			a = g.availActions(i, j)
			n = len(a)
			if n == 0 {
				continue
			}

			d = weight(1.0 / float64(n))
			for k := range a {
				w += d
				po = append(po, &pawnOpt{m: i, n: j, a: a[k], p: w})
			}

			w = 0
		}
	}

	return po
}

// TODO: rename to checkWin or something
// updateState checks the board for a win condition and sets the game state to the winning state. If the game has not been won, then it swaps the turn
func (g *game) updateState() {
	switch g.s {
	case whiteTurn:
		g.s = blackTurn
		for i := range g.b[0] {
			if g.b[0][i] == whitePawn {
				g.s = whiteWin
				break
			}
		}
	case blackTurn:
		g.s = whiteTurn
		n := len(g.b[0]) - 1
		for i := range g.b[n] {
			if g.b[n][i] == blackPawn {
				g.s = blackWin
				break
			}
		}
	default:
		g.s = illegal
	}
}
