package main

import "fmt"

type thing struct {
	a int
	b *int
}

type things []*thing

func main() {
	v := 1
	x := &thing{a: 1, b: &v}
	arr := things{}
	arr = append(arr, x)
	x.a++
	*x.b++
	fmt.Println(arr[0], x)
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
