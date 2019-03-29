package main

import (
	"bytes"
	"fmt"
	"strings"
)

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
	space     = pawn(' ')
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

	pvp mode = iota // Player vs player
	pvc             // Player vs computer
	cvp             // Computer vs player
	cvc             // Computer vs computer
)

func (g *game) String() string {
	n := len(g.b[0])
	bldr := strings.Builder{}
	bldr.Grow((2*len(g.b) + 1) * (2*n + 1))

	dashPlus := []byte{'-', '+'}
	barnl := []byte{'|', '\n'}
	line := make([]byte, 0, n+2)
	line = append(line, '+')
	line = append(line, bytes.Repeat(dashPlus, n)...)
	line = append(line, '\n')

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

func newGame(m, n int, md mode) *game {
	return &game{
		b: newBoard(m, n),
		s: whiteTurn,
		m: md,
		h: make(history, 0, 32), // TODO: Find the total number of legal states (24?)
	}
}

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

func (g *game) play() {
	g.h = append(g.h, copyBoard(g.b))

	for g.s == whiteTurn || g.s == blackTurn {
		g.turn()
		g.updateState()
		g.h = append(g.h, copyBoard(g.b))
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

func (g *game) availActions(m, n int) []action {
	a := make([]action, 0, 4)

	lenB := len(g.b)
	lenB0 := len(g.b[0])
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
			case lenB0 - 1:
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
			case lenB0 - 1:
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
			b[i] = append(b[i], space)
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
