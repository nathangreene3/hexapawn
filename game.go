package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// action represents what a player can do at a position (m,n) on a board.
type action byte

// mode indicates player vs player, player vs an auto player, or auto player vs
// auto player.
type mode byte

// pawn is a playable piece or a blank space.
type pawn byte

// side is either white or black.
type side byte

// state indicates how the game will procede.
type state byte

// history is a set of positions that occur in a single game.
type history []*event

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
	forward      action = iota // Move forward from side's perspective
	captureLeft                // Capture left from side's perspective
	captureRight               // Capture right from side's perspective

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
	illegal   state = iota // Game halts on illegal state
	whiteTurn              // Game continues
	blackTurn              // Game continues
	whiteWin               // Game is over when white wins
	blackWin               // Game is over when black wins
	stalemate              // Game is over when stalemate occurs
)

// String returns a string representing the current state of a game.
func (gm *game) String() string {
	n := len(gm.brd[0])
	bldr := strings.Builder{}
	bldr.Grow((2*len(gm.brd)+1)*(2*n+1) + 32)

	bldr.Write(gm.brd.toBytes())
	switch gm.st {
	case whiteTurn:
		bldr.Write([]byte("\nwhite to move\n"))
	case blackTurn:
		bldr.Write([]byte("\nblack to move\n"))
	case whiteWin:
		bldr.Write([]byte("\nwhite wins\n"))
	case blackWin:
		bldr.Write([]byte("\nblack wins\n"))
	case stalemate:
		bldr.Write([]byte("\nstalemate\n"))
	case illegal:
		bldr.Write([]byte("\nillegal position\n"))
	default:
		bldr.Write([]byte("\nunknown state\n"))
	}

	return bldr.String()
}

// newGame returns a game to be played.
func newGame(m, n int, md mode) *game {
	return &game{
		brd: newBoard(m, n),
		st:  whiteTurn,
		md:  md,
		hst: make(history, 0, 32), // TODO: Find an upper bound on the total number of legal states depending on input m and n
	}
}

// play
func play(m, n int, md mode) {
	gm := newGame(m, n, md)
	trainSessions := 100000
	learningRate := weight(0.1)
	var psn *position

	switch md {
	case cvc:
		white := newAutoPlayer(whiteSide, m, n)
		black := newAutoPlayer(blackSide, m, n)
		white.train(trainSessions, learningRate)
		black.train(trainSessions, learningRate)
		var gameOver bool

		for !gameOver {
			psn = &position{brd: gm.brd, st: gm.st, pos: availPawnOpts(gm.brd, gm.st)}

			switch gm.st {
			case whiteTurn:
				// fmt.Printf("%s\nwhite to move\n\n", gm.String())
				gm.move(white.chooseEvent(psn))
			case blackTurn:
				// fmt.Printf("%s\nblack to move\n\n", gm.String())
				gm.move(black.chooseEvent(psn))
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

	// fmt.Println(gm.String())
	// fmt.Println()
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

func playNGames(numGames, numTrainSessions int, learningRate weight, m, n int, md mode) string {
	var (
		gm         *game     // Game to be played
		psn        *position // Position at each turn
		whiteWins  int       // Number of white wins
		blackWins  int       // Number of black wins
		stalemates int       // Number of stalemates reached
	)

	switch md {
	case cvc:
		white := newAutoPlayer(whiteSide, m, n)
		black := newAutoPlayer(blackSide, m, n)
		white.train(numTrainSessions, learningRate)
		black.train(numTrainSessions, learningRate)

		for ; 0 < numGames; numGames-- {
			gm = newGame(m, n, md)
			fmt.Println(gm)
			for {
				psn = &position{
					brd: gm.brd,
					st:  gm.st,
					pos: availPawnOpts(gm.brd, gm.st),
				}

				if gm.st == whiteTurn {
					gm.move(white.chooseEvent(psn))
					fmt.Println(gm)
					continue
				}

				if gm.st == blackTurn {
					gm.move(black.chooseEvent(psn))
					fmt.Println(gm)
					continue
				}

				break
			}

			switch gm.st {
			case whiteWin:
				whiteWins++
			case blackWin:
				blackWins++
			case stalemate:
				stalemates++
			default:
				log.Fatal("playNGames: invalid endgame state")
			}
		}

		fmt.Printf("%s\nwhite boards: %d\nblack boards: %d\n\n", gm, len(white.psns), len(black.psns))
	case cvp: // TODO
	case pvc: // TODO
	case pvp: // TODO
	default:
		log.Fatal("playNGames: invalid mode")
	}

	return fmt.Sprintf("white wins:  %d\nblack wins:  %d\nstalemates:  %d\n---------------\n     total: %d", whiteWins, blackWins, stalemates, whiteWins+blackWins+stalemates)
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
				if m+1 < len(gm.brd) && gm.brd[m+1][n] == space {
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
		case space:
			fallthrough
		default:
			panic("move: cannot move space")
		}
	} else {
		gm.st = stalemate // No pawn option selected is stalemate
	}

	gm.hst = append(gm.hst, evnt)
}

// availActions returns a set of actions that can be taken at a position (m,n).
// Actions are available if the state is either white or black turn.
func availActions(m, n int, brd board, st state) []action {
	acts := make([]action, 0, 4) // Actions to return
	lenB := len(brd)             // Number of rows
	lenB0m1 := len(brd[0]) - 1   // Number of columns minus one

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

	sort.Slice(acts, func(i, j int) bool { return acts[i] < acts[j] })
	return acts
}

// checkWin checks the board for a win condition given a state. If the state is
// neither white nor black turn, then false is returned.
func checkWin(brd board, st state) bool {
	switch st {
	case whiteTurn:
		// Check if any white pawns reached top row
		for i := range brd[0] {
			if brd[0][i] == whitePawn {
				return true
			}
		}

		// Check if any black pieces remain
		for i := 0; i < len(brd)-1; i++ {
			for _, p := range brd[i] {
				if p == blackPawn {
					return false
				}
			}
		}

		return true // White hasn't reached top row, but no pieces left for black to move
	case blackTurn:
		// Check if any black pawns reached bottom row
		n := len(brd) - 1 // Index of bottom row
		for i := range brd[n] {
			if brd[n][i] == blackPawn {
				return true
			}
		}

		// Check if any white pieces remain
		for i := 1; i < len(brd); i++ {
			for _, p := range brd[i] {
				if p == whitePawn {
					return false
				}
			}
		}

		return true // Black has not reached bottom row, but no pieces left for white to move
	default:
		return false // Neither white nor black's turn; the game is not over
	}
}
