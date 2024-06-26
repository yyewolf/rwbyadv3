package maze

import (
	"fmt"
	"testing"
)

func DrawMaze(g *Grid) {
	g.ForEach(func(n *Node) bool {
		if n.Goes(Right) {
			fmt.Print("→")
		}
		if n.Goes(Down) {
			fmt.Print("↓")
		}
		if n.Goes(Left) {
			fmt.Print("←")
		}
		if n.Goes(Up) {
			fmt.Print("↑")
		}

		if n.H == g.Width-1 {
			fmt.Print("\n")
		}
		return false
	})
}

func TestDefaultMaze(t *testing.T) {
	g := DefaultMaze(10, 10)
	DrawMaze(g)
}

func TestStep(t *testing.T) {
	g := DefaultMaze(10, 10)
	DrawMaze(g)

	fmt.Println()
	g.Init()
	g.Step()

	DrawMaze(g)
}

func TestGen(t *testing.T) {
	g := DefaultMaze(10, 10)
	DrawMaze(g)

	fmt.Println()
	g.Generate()

	DrawMaze(g)
}