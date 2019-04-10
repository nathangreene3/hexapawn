package main

import (
	"bytes"
)

// board is an m-by-n array of pieces.
type board [][]pawn

func (brd board) String() string {
	return string(brd.toBytes())
}

func (brd board) toBytes() []byte {
	n := len(brd[0])
	buf := bytes.Buffer{}
	buf.Grow((2*len(brd) + 1) * (2*n + 1))

	dashPlus := []byte{'-', '+'}
	barnl := []byte{'|', '\n'}
	line := bytes.Join([][]byte{[]byte{'+'}, bytes.Repeat(dashPlus, n), []byte{'\n'}}, []byte{})

	for i := range brd {
		buf.Write(line)

		for j := range brd[i] {
			buf.WriteByte('|')
			buf.WriteByte(byte(brd[i][j]))
		}

		buf.Write(barnl)
	}

	buf.WriteByte('+')
	buf.Write(bytes.Repeat(dashPlus, n))

	return buf.Bytes()
}

// newBoard returns a new board with black on top, white on bottom. Panics if m or
// n are less than three.
func newBoard(m, n int) board {
	if m < 3 || n < 3 {
		panic("newBoard: diminsions cannot be less than three")
	}

	brd := make(board, 0, m)

	// Add black pawns to first row
	brd = append(brd, make([]pawn, 0, n))
	for i := 0; i < n; i++ {
		brd[0] = append(brd[0], blackPawn)
	}

	// Add spaces to middle rows
	for i := 1; i < m-1; i++ {
		brd = append(brd, make([]pawn, 0, n))
		for j := 0; j < n; j++ {
			brd[i] = append(brd[i], space)
		}
	}

	// Add white pawns to last row
	brd = append(brd, make([]pawn, 0, n))
	for i := 0; i < n; i++ {
		brd[m-1] = append(brd[m-1], whitePawn)
	}

	return brd
}

// copyBoard returns a new copy of a board.
func copyBoard(brd board) board {
	cpy := make(board, 0, len(brd))
	n := len(brd[0])

	for i := range brd {
		cpy = append(cpy, make([]pawn, n))
		copy(cpy[i], brd[i])
	}

	return cpy
}

// equalBoards returns true if two boards are equal in dimension and position and
// false if otherwise.
func equalBoards(brd0, brd1 board) bool {
	m := len(brd0)
	if m != len(brd1) {
		panic("lessBoards: cannot compare boards of differing dimensions")
	}

	n := len(brd0[0])
	if n != len(brd1[0]) {
		panic("lessBoards: cannot compare boards of differing dimensions")
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if brd0[i][j] != brd1[i][j] {
				return false
			}
		}
	}

	return true
}

// symmetricEqualBoards returns true if two boards are equal reflection across the
// vertical axis and false if otherwise. That is, b == reflect(c) is returned.
func symmetricEqualBoards(brd0, brd1 board) bool {
	m := len(brd0)
	if m != len(brd1) {
		return false
	}

	n := len(brd0[0])
	if n != len(brd1[0]) {
		return false
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if brd0[i][j] != brd1[i][n-j-1] {
				return false
			}
		}
	}

	return true
}

func lessBoards(brd0, brd1 board) bool {
	if compareBoards(brd0, brd1) < 0 {
		return true
	}

	return false
}

func lessEqBoards(brd0, brd1 board) bool {
	if 0 < compareBoards(brd0, brd1) {
		return false
	}

	return true
}

func compareBoards(brd0, brd1 board) int {
	m := len(brd0)
	if m != len(brd1) {
		panic("lessBoards: cannot compare boards of differing dimensions")
	}

	n := len(brd0[0])
	if n != len(brd1[0]) {
		panic("lessBoards: cannot compare boards of differing dimensions")
	}

	var x, y pawn
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			x, y = brd0[i][j], brd1[i][j]
			if x < y {
				return -1
			}

			if y < x {
				return 1
			}
		}
	}

	return 0
}
