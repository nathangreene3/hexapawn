package main

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

// symmetricEqualBoards returns true if two boards are equal under row reflection and false if otherwise. That is, b == reflect(c) is returned.
func symmetricEqualBoards(b, c board) bool {
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
			if b[i][j] != c[i][n-j-1] {
				return false
			}
		}
	}

	return true
}
