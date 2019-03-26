package main

var (
	// White to move from start
	p00 = board{
		[]pawn{'b', 'b', 'b'},
		[]pawn{' ', ' ', ' '},
		[]pawn{'w', 'w', 'w'},
	}

	// FIRST MOVE
	// White to move from p00
	p01 = board{ // forward (2,0)
		[]pawn{'b', 'b', 'b'},
		[]pawn{'w', ' ', ' '},
		[]pawn{' ', 'w', 'w'},
	}
	p02 = board{ // forward (2,1)
		[]pawn{'b', 'b', 'b'},
		[]pawn{' ', 'w', ' '},
		[]pawn{'w', ' ', 'w'},
	}
	p03 = board{ // forward (2,2)
		[]pawn{'b', 'b', 'b'},
		[]pawn{' ', ' ', 'w'},
		[]pawn{'w', 'w', ' '},
	}

	// SECOND MOVE
	// Black to move from p01
	p04 = board{ // captureRight (0,1)
		[]pawn{'b', ' ', 'b'},
		[]pawn{'b', ' ', ' '},
		[]pawn{' ', 'w', 'w'},
	}
	p05 = board{ // forward (0,1)
		[]pawn{'b', ' ', 'b'},
		[]pawn{'w', 'b', ' '},
		[]pawn{' ', 'w', 'w'},
	}
	p06 = board{ // forward (1,2)
		[]pawn{'b', 'b', ' '},
		[]pawn{'w', ' ', 'b'},
		[]pawn{' ', 'w', 'w'},
	}

	// Black to move from p02
	p07 = board{ // captureLeft (0,0)
		[]pawn{' ', 'b', 'b'},
		[]pawn{' ', 'b', ' '},
		[]pawn{'w', ' ', 'w'},
	}
	p08 = board{ // forward (0,0)
		[]pawn{' ', 'b', 'b'},
		[]pawn{'b', 'w', ' '},
		[]pawn{'w', ' ', 'w'},
	}
	p09 = board{ // captureRight (0,2)
		[]pawn{'b', 'b', ' '},
		[]pawn{' ', 'w', ' '},
		[]pawn{'w', ' ', 'w'},
	}
	p10 = board{ // forward (0,2)
		[]pawn{'b', 'b', ' '},
		[]pawn{' ', 'w', 'b'},
		[]pawn{'w', ' ', 'w'},
	}

	// Black to move from p03
	p11 = board{ // forward (0,0)
		[]pawn{' ', 'b', 'b'},
		[]pawn{'b', ' ', 'w'},
		[]pawn{'w', 'w', ' '},
	}
	p12 = board{ // captureLeft (0,1)
		[]pawn{'b', ' ', 'b'},
		[]pawn{' ', ' ', 'b'},
		[]pawn{'w', 'w', ' '},
	}
	p13 = board{ // forward (0,1)
		[]pawn{'b', ' ', 'b'},
		[]pawn{' ', 'b', 'w'},
		[]pawn{'w', 'w', ' '},
	}

	// THIRD MOVE
	// White to move from p04
	p14 = board{ // captureLeft (2,1)
		[]pawn{'b', ' ', 'b'},
		[]pawn{'w', ' ', ' '},
		[]pawn{' ', ' ', 'w'},
	}
	p15 = board{ // forward (2,1)
		[]pawn{'b', ' ', 'b'},
		[]pawn{'b', 'w', ' '},
		[]pawn{' ', ' ', 'w'},
	}
	p16 = board{ // forward (2,2)
		[]pawn{'b', ' ', 'b'},
		[]pawn{'b', ' ', 'w'},
		[]pawn{' ', 'w', ' '},
	}

	// White to move from p05
	p17 = board{ // captureLeft (2,0)
		[]pawn{'b', ' ', 'b'},
		[]pawn{'w', 'w', ' '},
		[]pawn{' ', 'w', ' '},
	}
	p18 = board{ // forward (2,0)
		[]pawn{'b', ' ', 'b'},
		[]pawn{'w', 'b', 'w'},
		[]pawn{' ', 'w', ' '},
	}
)

type weight float64
type autoPlayer []position

type pawnOpt struct {
	m int
	n int
	a action
	p weight
}

type position struct {
	b  board
	s  state
	po []*pawnOpt
}

// func (ap autoPlayer) insertPosition(p *position) {
// 	sort.Sort(ap)
// }

// func (ap autoPlayer)Less(i,j int)bool{
// 	return
// }
