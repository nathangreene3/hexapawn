package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// rand.Seed(int64(time.Now().Nanosecond()))
	// fmt.Println(playNGames(100, 3, 3, cvc))

	A := []int{1, 2, 5, 6, 8, 8, 9}
	x := 2 // Value to find in A

	i := 0
	j := len(A) - 1
	var k int
	for i < j {
		k = i + int(uint(j-i)>>1)
		if x <= A[k] {
			j = k
		} else {
			i = k + 1
		}
	}
	fmt.Println(A, x, i)

	// n := 5
	// brds := make([]board, 0, n)
	// for i := 0; i < n; i++ {
	// 	brds = append(brds, randBoard(2, 2))
	// }

	// sort.Slice(
	// 	brds,
	// 	func(i, j int) bool {
	// 		if compareBoards(brds[i], brds[j]) < 0 {
	// 			return true
	// 		}
	// 		return false
	// 	},
	// )

	// for i := 0; i < n; i++ {
	// 	fmt.Printf("%s\n\n", brds[i])
	// }

	// brd := copyBoard(brds[rand.Intn(n)])
	// i, j := 0, n-1
	// var k int
	// for i < j {
	// 	k = i + (j-i)/2
	// 	if 0 < compareBoards(brd, brds[k]) {
	// 		i = k + 1
	// 	} else {
	// 		j = k
	// 	}
	// }
	// fmt.Println(i, brd)
}

func readMove(s state) (int, int, action) {
	// EXAMPLE
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter text: ")
	// text, _ := reader.ReadString('\n')
	// fmt.Println(text)
	var m, n int
	var a action
	// r := bufio.NewReader(os.Stdin)
	// switch s {
	// case whiteTurn:
	// 	fmt.Print("White to move")
	// 	input, _ := r.ReadString('\n')
	// }
	return m, n, a
}

func randBoard(m, n int) board {
	brd := make(board, 0, m)
	for i := 0; i < m; i++ {
		brd = append(brd, make([]pawn, 0, n))
		for j := 0; j < n; j++ {
			brd[i] = append(brd[i], randPawn())
		}
	}

	return brd
}

func randPawn() pawn {
	switch rand.Intn(3) {
	case 0:
		return whitePawn
	case 1:
		return blackPawn
	default:
		return space
	}
}
