package main

import "fmt"

type action byte
type mode byte
type pawn byte
type state byte
type board [][]pawn
type history []board

type game struct {
	b board
	s state
	m mode
	h history
}

const (
	noPawn    = pawn(' ')
	whitePawn = pawn('w')
	blackPawn = pawn('b')

	illegal state = iota
	stalemate
	whiteTurn
	blackTurn
	whiteWin
	blackWin

	forward action = iota
	captureLeft
	captureRight

	pvp mode = iota
	pvc
	cvp
	cvc
)

func newGame(m, n int, md mode) *game {
	return &game{
		b: newBoard(m, n),
		s: whiteTurn,
		m: md,
		h: make(history, 0, 32), // TODO: Find the total number of states (24?)
	}
}

func (g *game) turn() bool {
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
	return false
}

func (g *game) play() {
	g.h = append(g.h, copyBoard(g.b))

	for g.s == whiteTurn || g.s == blackTurn {
		if g.turn() {
			g.h = append(g.h, copyBoard(g.b))
			g.updateState()
		}

		g.s = stalemate
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

func newBoard(m, n int) board {
	if m < 3 || n < 3 {
		panic("newBoard: diminsions cannot be less than three")
	}

	b := make(board, 0, m)

	b = append(b, make([]pawn, 0, n))
	for i := 0; i < n; i++ {
		b[0] = append(b[0], blackPawn)
	}

	for i := 1; i < m-1; i++ {
		b = append(b, make([]pawn, 0, n))
		for j := 0; j < n; j++ {
			b[i] = append(b[i], noPawn)
		}
	}

	b = append(b, make([]pawn, 0, n))
	for i := 0; i < n; i++ {
		b[m-1] = append(b[m-1], whitePawn)
	}

	return b
}

func copyBoard(b board) board {
	c := make(board, 0, len(b))
	n := len(b[0])

	for i := range b {
		c = append(c, make([]pawn, n))
		copy(c[i], b[i])
	}

	return c
}

// move returns true if the game board is altered (a player has selected a valid move) and false if otherwise.
func (g *game) move(m, n int, a action) bool {
	switch g.b[m][n] {
	case whitePawn:
		switch a {
		case forward:
			if 0 < m && g.b[m-1][n] == noPawn {
				g.b[m-1][n] = whitePawn
				g.b[m][n] = noPawn
				return true
			}
		case captureLeft:
			if 0 < m && 0 < n && g.b[m-1][n-1] == blackPawn {
				g.b[m-1][n-1] = whitePawn
				g.b[m][n] = noPawn
				return true
			}
		case captureRight:
			if 0 < m && n+1 < len(g.b[0]) && g.b[m-1][n+1] == blackPawn {
				g.b[m-1][n+1] = whitePawn
				g.b[m][n] = noPawn
				return true
			}
		}

		return false
	case blackPawn:
		switch a {
		case forward:
			if m+1 < len(g.b) && g.b[m][n] == noPawn {
				g.b[m+1][n] = blackPawn
				g.b[m][n] = noPawn
				return true
			}
		case captureLeft:
			if m+1 < len(g.b) && n+1 < len(g.b[0]) && g.b[m+1][n+1] == whitePawn {
				g.b[m+1][n+1] = blackPawn
				g.b[m][n] = noPawn
				return true
			}
		case captureRight:
			if m+1 < len(g.b) && n-1 < len(g.b[0]) && g.b[m+1][n-1] == whitePawn {
				g.b[m+1][n-1] = blackPawn
				g.b[m][n] = noPawn
				return true
			}
		}

		return false
	}
	return false
}

func (g *game) updateState() {
	switch g.s {
	case whiteTurn:
		for _, b := range g.b[0] {
			if b == whitePawn {
				g.s = whiteWin
				return
			}
		}

		g.s = blackTurn
	case blackTurn:
		for _, b := range g.b[len(g.b)-1] {
			if b == blackPawn {
				g.s = blackWin
				return
			}
		}

		g.s = whiteTurn
	default:
		g.s = illegal
	}
}

// equalBoards returns true if two boards are equal in dimension and position and false if otherwise.
func equalBoards(b, c board) bool {
	m := len(b)
	if m != len(c) {
		return false
	}

	n := len(b[0])
	if n != len(c[0]) {
		return false
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if b[i][j] != c[i][j] {
				return false
			}
		}
	}

	return true
}
