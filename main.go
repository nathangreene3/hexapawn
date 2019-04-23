package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	fmt.Println(playNGames(1, 1000, 0.1, 3, 3, cvc))
}

func run() {

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
