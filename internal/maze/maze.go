package maze

import "math/rand"

func NewGrid(r *rand.Rand, height, width int) *Grid {
	var grid Grid
	grid.gridArray = make([][]*Node, height)

	for v := range height {
		grid.gridArray[v] = make([]*Node, width)

		for h := range width {
			grid.gridArray[v][h] = &Node{
				V: v,
				H: h,
			}
		}
	}

	grid.r = r
	grid.h = height
	grid.w = width

	grid.ForEach(func(n *Node) bool {
		if n.H == width-1 {
			// last column all goes down
			n.SetLink(Down)

			if n.V == height-1 {
				n.SetLink(Up)
			}

			return false // continue
		}
		// everything else goes right
		n.SetLink(Right)
		return false // continue
	})

	return &grid
}

func (g *Grid) At(v, h int) (*Node, bool) {
	if v < 0 || v > g.h-1 {
		return nil, false
	}
	if h < 0 || h > g.w-1 {
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

func (g *Grid) Init() {
	v := g.r.Intn(g.h)
	h := g.r.Intn(g.w)

	n, _ := g.At(v, h)
	g.pointerAt = n
}

func (g *Grid) Step() {
	neighbors := g.pointerAt.Neighbors(g)
	r := g.r.Intn(len(neighbors))
	newPointer := neighbors[r]
	g.pointerAt.SetLink(g.pointerAt.DirTo(newPointer))
	g.pointerAt = newPointer
}

func (g *Grid) Generate() {
	var rounds = g.r.Intn(1000) + 1000 + g.h*g.w*10
	g.Init()
	for range rounds {
		g.Step()
	}
}

func (g *Grid) Expand(pathSize int) [][]int {
	// return a grid of 1s and 0s, 1s are walls, 0s are paths
	// paths are pathSize wide
	var expandedHeight = g.h*2*pathSize + pathSize
	var expandedWidth = g.w*2*pathSize + pathSize

	var expanded [][]int = make([][]int, expandedHeight)
	for v := range expanded {
		expanded[v] = make([]int, expandedWidth)
	}

	for v := range expanded {
		for h := range expanded[v] {
			expanded[v][h] = 1
		}
	}

	g.ForEach(func(n *Node) bool {
		expandedV := n.V*2*pathSize + pathSize
		expandedH := n.H*2*pathSize + pathSize

		expanded[expandedV][expandedH] = 0
		for i := 0; i < pathSize; i++ {
			for j := 0; j < pathSize; j++ {
				expanded[expandedV+i][expandedH+j] = 0
			}
		}

		if n.Goes(Up) {
			for i := 0; i < pathSize; i++ {
				for j := 0; j < pathSize; j++ {
					expanded[expandedV-(i+1)][expandedH+j] = 0
				}
			}
		}

		if n.Goes(Down) {
			for i := 0; i < pathSize; i++ {
				for j := 0; j < pathSize; j++ {
					expanded[expandedV+(i+pathSize)][expandedH+j] = 0
				}
			}
		}

		if n.Goes(Left) {
			for i := 0; i < pathSize; i++ {
				for j := 0; j < pathSize; j++ {
					expanded[expandedV+i][expandedH-(j+1)] = 0
				}
			}
		}

		if n.Goes(Right) {
			for i := 0; i < pathSize; i++ {
				for j := 0; j < pathSize; j++ {
					expanded[expandedV+i][expandedH+(j+pathSize)] = 0
				}
			}
		}
		return false
	})

	return expanded
}
