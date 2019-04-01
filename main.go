package main

import "fmt"

type thing struct {
	a int
	b *int
}

type things []*thing

func main() {
	// g := newGame(3, 3, pvp)
	// fmt.Println(g.String())

	v := 1
	x := thing{a: 1, b: &v}
	y := &thing{a: 1, b: &v}
	ts := things{}
	fmt.Println(ts, x, y)
	insert(ts, &x)
	x.a++
	*x.b++
	y.a++
	*y.b++
	fmt.Println(ts, x, y)
}

func insert(ts things, t *thing) {
	ts = append(ts, t)
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
