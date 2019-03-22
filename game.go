package main

import "fmt"

type action byte
type mode byte
type pawn byte
type state byte
type board [][]pawn

type game struct {
	b board
	s state
	m mode
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
	}
}

func (g *game) turn() {
	switch g.s {
	case whiteTurn:
		switch g.m {
		case pvp, pvc:
			// player move
		case cvp, cvc:
			// random move by npc
		}

		g.updateState()
	case blackTurn:
		switch g.m {
		case pvp, cvp:
			// player move
		case pvc, cvc:
			// random move by npc
		}

		g.updateState()
	}
}

func (g *game) play() {
	for {
		g.turn()
		switch g.s {
		case whiteWin:
			fmt.Println("WHITE WINS")
			break
		case blackWin:
			fmt.Println("BLACK WINS")
			break
		case illegal:
			fmt.Println("ILLEGAL STATE")
			break
		case stalemate:
			fmt.Println("STALEMATE")
		}
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

func (b board) move(m, n int, a action) {
	switch b[m][n] {
	case whitePawn:
		switch a {
		case forward:
			if 0 < m && b[m-1][n] == noPawn {
				b[m-1][n] = whitePawn
				b[m][n] = noPawn
			}
		case captureLeft:
			if 0 < m && 0 < n && b[m-1][n-1] == blackPawn {
				b[m-1][n-1] = whitePawn
				b[m][n] = noPawn
			}
		case captureRight:
			if 0 < m && n+1 < len(b[0]) && b[m-1][n+1] == blackPawn {
				b[m-1][n+1] = whitePawn
				b[m][n] = noPawn
			}
		}
	case blackPawn:
		switch a {
		case forward:
			if m+1 < len(b) && b[m][n] == noPawn {
				b[m+1][n] = blackPawn
				b[m][n] = noPawn
			}
		case captureLeft:
			if m+1 < len(b) && n+1 < len(b[0]) && b[m+1][n+1] == whitePawn {
				b[m+1][n+1] = blackPawn
				b[m][n] = noPawn
			}
		case captureRight:
			if m+1 < len(b) && n-1 < len(b[0]) && b[m+1][n-1] == whitePawn {
				b[m+1][n-1] = blackPawn
				b[m][n] = noPawn
			}
		}
	}
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
