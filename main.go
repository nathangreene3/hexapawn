package main

import "fmt"

func main() {
	g := newGame(3, 3, pvp)
	fmt.Println(g.String())
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
