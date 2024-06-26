package maze

import (
	"math/rand"
)

func NewGrid(height, width int) *Grid {
	var grid [][]*Node = make([][]*Node, height)

	for v := range height {
		grid[v] = make([]*Node, width)

		for h := range width {
			grid[v][h] = &Node{
				V: v,
				H: h,
			}
		}
	}

	return &Grid{
		Width:     width,
		Height:    height,
		gridArray: grid,
	}
}

func (g *Grid) At(v, h int) (*Node, bool) {
	if v < 0 || v > g.Height-1 {
		return nil, false
	}
	if h < 0 || h > g.Width-1 {
		return nil, false
	}
	return g.gridArray[v][h], true
}

func (g *Grid) ForEach(callback func(n *Node) bool) {
	for v := range g.gridArray {

		var b bool
		for h := range g.gridArray[v] {
			n, _ := g.At(v, h)
			if callback(n) {
				b = true
				break
			}
		}
		if b {
			break
		}
	}
}

func DefaultMaze(height, width int) *Grid {
	grid := NewGrid(height, width)

	grid.ForEach(func(n *Node) bool {
		if n.H == width-1 {
			// last coloum all goes down
			n.SetLink(Down)
			return false // continue
		}
		// everything else goes right
		n.SetLink(Right)
		return false // continue
	})

	return grid
}

func (g *Grid) Init() {
	v := rand.Intn(g.Height)
	h := rand.Intn(g.Width)

	n, _ := g.At(v, h)
	g.pointerAt = n
}

func (g *Grid) Step() {
	neighbors := g.pointerAt.Neighbors(g)
	r := rand.Intn(len(neighbors))
	newPointer := neighbors[r]
	g.pointerAt.SetLink(g.pointerAt.DirTo(newPointer))
	g.pointerAt = newPointer
}

func (g *Grid) Generate() {
	var rounds = rand.Intn(1000) + 1000
	g.Init()
	for range rounds {
		g.Step()
	}
}
