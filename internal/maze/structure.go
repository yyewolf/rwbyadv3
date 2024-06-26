package maze

const (
	Up = 1 << iota
	Down
	Left
	Right
)

type Node struct {
	V, H  int
	Links int8
}

func (n *Node) Goes(dir int8) bool {
	return dir&n.Links == dir
}

func (n *Node) SetLink(dir int8) {
	n.Links = dir
}

func (n *Node) AddLink(dir int8) {
	n.Links |= dir
}

func (n *Node) DelLink(dir int8) {
	n.Links &= ^dir
}

type gridArray [][]*Node

type Grid struct {
	Width, Height int

	pointerAt *Node
	gridArray
}

func (n *Node) Neighbors(g *Grid) []*Node {
	var neighbors []*Node
	if n, exist := g.At(n.V, n.H-1); exist {
		neighbors = append(neighbors, n)
	}
	if n, exist := g.At(n.V-1, n.H); exist {
		neighbors = append(neighbors, n)
	}
	if n, exist := g.At(n.V, n.H+1); exist {
		neighbors = append(neighbors, n)
	}
	if n, exist := g.At(n.V+1, n.H); exist {
		neighbors = append(neighbors, n)
	}
	return neighbors
}

func (n *Node) DirTo(n2 *Node) int8 {
	if n2.H == n.H-1 && n2.V == n.V {
		return Left
	}
	if n2.H == n.H+1 && n2.V == n.V {
		return Right
	}
	if n2.H == n.H && n2.V == n.V-1 {
		return Up
	}
	if n2.H == n.H && n2.V == n.V+1 {
		return Down
	}
	return 0
}
