package main

//#region Imports
import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

//#endregion

//#region Structures
type Point struct {
	X, Y int
}

type Node struct {
	X, Y, Depth int
	Dungeon     *Dungeon
	Parent      *Node
}

type Rectangle struct {
	X, Y, Width, Height int
}

type Dungeon struct {
	Grid                                   [][]int
	Capacity, NumRooms, Attempts, Min, Max int
	Rooms                                  []Rectangle
	Regions                                []int
	Border                                 Rectangle
	Seed                                   Randomizer
}

//#endregion

type Randomizer interface {
	Int() int
}

func (r *Rectangle) Contains(p Point) bool {
	return p.X >= r.X &&
		p.Y >= r.Y &&
		p.X <= r.X+r.Width &&
		p.Y <= r.Y+r.Height
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (r *Rectangle) Intersects(b *Rectangle) bool {
	return Max(r.X, b.X) < Min(r.X+r.Width, b.X+b.Width) &&
		Max(r.Y, b.Y) < Min(r.Y+r.Height, b.Y+b.Height)
}

func (r *Rectangle) Intersection(b *Rectangle) *Rectangle {
	x := Max(r.X, b.X)
	y := Max(r.Y, b.Y)
	w := Min(r.X+r.Width, b.X+b.Width) - x
	h := Min(r.Y+r.Height, b.Y+b.Height) - y
	return &Rectangle{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
	}
}

func (d *Dungeon) CarvePoint(p Point, cellType int) {
	if d.Border.Contains(p) {
		d.Grid[p.X][p.Y] = cellType
	}
}

func (d *Dungeon) CarveRect(rect *Rectangle, cellType int) {
	for x := 0; x < rect.Width; x++ {
		for y := 0; y < rect.Height; y++ {
			d.CarvePoint(Point{X: x + rect.X, Y: y + rect.Y}, cellType)
		}
	}
}

func (d *Dungeon) AddWalls() {
	for x := 1; x < d.Capacity-1; x++ {
		for y := 1; y < d.Capacity-1; y++ {
			if d.Grid[x][y] == 1 {
				for j := 0; j < 4; j++ {
					gX := (j/2%2*2-1)*(j%2) + x
					gY := ((1-j/2%2)*2-1)*(1-j%2) + y
					if d.Grid[gX][gY] == 0 {
						d.Grid[gX][gY] = 2
					}
				}
			}
		}
	}
}

func (d *Dungeon) AddRoom() {
	for i := 0; i < d.Attempts; i++ {
		w := RandInt(d.Min, d.Max, d.Seed)
		h := RandInt(d.Min, d.Max, d.Seed)
		x := RandInt(1, d.Capacity-w-1, d.Seed)
		y := RandInt(1, d.Capacity-w-1, d.Seed)

		rect := Rectangle{X: x, Y: y, Width: w, Height: h}
		bounds := Rectangle{X: x - 1, Y: y - 1, Width: w + 2, Height: h + 2}

		intersect := false
		for _, room := range d.Rooms {
			if (&bounds).Intersects(&room) {
				intersect = true
				break
			}
		}
		if !intersect {
			d.Rooms = append(d.Rooms, rect)
			d.CarveRect(&rect, 1)
			return
		}

	}
}

func (d *Dungeon) SetRegion(x, y, region int) {
	if d.Border.Contains(Point{X: x, Y: y}) && d.Grid[x][y] != region && d.Grid[x][y] > 0 {
		d.Grid[x][y] = region
		d.Regions[region]++
		d.SetRegion(x-1, y, region)
		d.SetRegion(x+1, y, region)
		d.SetRegion(x, y-1, region)
		d.SetRegion(x, y+1, region)
	}
}

type NodeList []*Node

func (c NodeList) Len() int           { return len(c) }
func (c NodeList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c NodeList) Less(i, j int) bool { return c[i].Depth < c[j].Depth }

func (d *Dungeon) NewNode(x, y int, parent *Node) *Node {
	depth := 0
	if parent != nil {
		depth = parent.Depth + 1
	}

	return &Node{
		X:       x,
		Y:       y,
		Parent:  parent,
		Depth:   depth,
		Dungeon: d,
	}
}

func (n *Node) Carve(cellType int) {
	n.Dungeon.CarvePoint(Point{X: n.X, Y: n.Y}, cellType)
}

func (d *Dungeon) PrintNodes(nodes [][]*Node) {
	for x, row := range d.Grid {
		for y, col := range row {
			if nodes[x][y] == nil || nodes[x][y].Depth != 0 {
				if col == 0 {
					fmt.Print(" ")
				} else {
					fmt.Print(col)
				}
			} else {
				fmt.Print(nodes[x][y].Depth)
			}
		}
		fmt.Println()
	}
}

func (d *Dungeon) Extend(region int) {
	nodes := make([][]*Node, d.Capacity)
	for i := 0; i < d.Capacity; i++ {
		nodes[i] = make([]*Node, d.Capacity)
	}
	nodeList := NodeList{}
	for x := 0; x < d.Capacity; x++ {
		for y := 0; y < d.Capacity; y++ {
			if d.Grid[x][y] == region {
				nodes[x][y] = d.NewNode(x, y, nil)
				nodeList = append(nodeList, nodes[x][y])
			}
		}
	}

	for len(nodeList) > 0 {
		node := nodeList[0]
		nodeList = nodeList[1:]
		sort.Sort(nodeList)
		for j := 0; j < 4; j++ {
			x := (j/2%2*2-1)*(j%2) + node.X
			y := ((1-j/2%2)*2-1)*(1-j%2) + node.Y
			if x >= 0 && x < d.Capacity-1 && y >= 0 && y < d.Capacity-1 {
				if nodes[x][y] == nil {
					nodes[x][y] = d.NewNode(x, y, node)
					nodeList = append(nodeList, nodes[x][y])
					if d.Grid[x][y] != region && d.Grid[x][y] != 0 {
						d.SetRegion(x, y, region)
						p := nodes[x][y]
						for p.Parent != nil {
							p.Carve(region)
							p.Depth = 0
							p = p.Parent
						}
						return
					}
				}
			}
		}
	}
}

func (d *Dungeon) Initiate() *Dungeon {
	for i := 0; i < d.NumRooms; i++ {
		d.AddRoom()
	}

	used := map[Point]bool{}
	var connected []int

	for i, roomA := range d.Rooms {
		for j, roomB := range d.Rooms {
			if i == j {
				continue
			}
			if used[Point{X: i, Y: j}] || used[Point{X: j, Y: i}] {
				continue
			}
			used[Point{X: i, Y: j}] = true

			boundsA := &Rectangle{
				X:      roomA.X - 1,
				Y:      roomA.Y - 1,
				Width:  roomA.Width + 2,
				Height: roomA.Height + 2,
			}
			boundsB := &Rectangle{
				X:      roomB.X - 1,
				Y:      roomB.Y - 1,
				Width:  roomB.Width + 2,
				Height: roomB.Height + 2,
			}

			if boundsA.Intersects(boundsB) {
				if !Contains(connected, i) {
					connected = append(connected, i)
				}
				if !Contains(connected, j) {
					connected = append(connected, j)
				}

				intersect := boundsA.Intersection(boundsB)
				if intersect.Width > 2 {
					intersect.X += 1
					intersect.Width -= 2
				} else if intersect.Height > 2 {
					intersect.Y += 1
					intersect.Height -= 2
				} else if intersect.Width*intersect.Height == 2 {
					d.CarveRect(intersect, 1)
					continue
				} else {
					continue
				}

				d.CarvePoint(Point{
					X: RandInt(intersect.X, intersect.X+intersect.Width, d.Seed),
					Y: RandInt(intersect.Y, intersect.Y+intersect.Height, d.Seed),
				}, 1)
			}

		}
	}

	var region int
	for {
		region = 2
		d.Regions = []int{0, 0}
		for _, room := range d.Rooms {
			if room.X < d.Capacity && room.Y < d.Capacity && d.Grid[room.X][room.Y] == 1 {
				d.Regions = append(d.Regions, 0)
				d.SetRegion(room.X, room.Y, region)
				region++
			}
		}

		// max int value
		max := int(^uint(0) >> 1)
		regionNum := -1
		for i := 2; i < region; i++ {
			if d.Regions[i] < max {
				max = d.Regions[i]
				regionNum = i
			}
		}

		d.Extend(regionNum)
		for x := 0; x < d.Capacity; x++ {
			for y := 0; y < d.Capacity; y++ {
				if d.Grid[x][y] != 0 {
					d.Grid[x][y] = 1
				}
			}
		}

		if region <= 3 {
			break
		}
	}
	return d
}

func (d *Dungeon) Print() {
	for _, row := range d.Grid {
		for _, col := range row {
			if col == 0 {
				fmt.Print("-")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

func RandInt(min int, max int, rnd Randomizer) int {
	if max == min {
		return min
	}
	return rnd.Int()%(max-min) + min
}

func Contains(arr []int, item int) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func Arguments() map[string]int64 {
	m := make(map[string]int64)
	m["size"] = 50
	m["rooms"] = 30
	m["attempts"] = 50
	m["min"] = 5
	m["max"] = 15
	m["seed"] = time.Now().Unix()

	// Parse arguments
	for _, argument := range os.Args[1:] {
		value := regexp.MustCompile("(\\w+)=(\\d+)").FindAllStringSubmatch(argument, -1)
		if value != nil {
			m[value[0][1]], _ = strconv.ParseInt(value[0][2], 10, 64)
		}
	}

	if m["min"] > m["max"] {
		m["min"] = m["max"]
	}
	return m
}

func main() {
	// Initiate default values
	m := Arguments()

	// Get values to generate dungeon
	size := int(m["size"])
	rooms := int(m["rooms"])
	attempts := int(m["attempts"])
	min := int(m["min"])
	max := int(m["max"])
	seed := m["seed"]

	// Initiate grid with provided size
	grid := make([][]int, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]int, size)
	}

	// Initiate new Dungeon struct
	dungeon := &Dungeon{
		Grid:     grid,
		Capacity: size,
		NumRooms: rooms,
		Attempts: attempts,
		Min:      min,
		Max:      max,
		Rooms:    []Rectangle{},
		Regions:  []int{},
		Border:   Rectangle{X: 1, Y: 1, Width: size - 2, Height: size - 2},
		Seed:     rand.New(rand.NewSource(seed)),
	}

	fmt.Println("Size:", size)
	fmt.Println("Rooms:", rooms)
	fmt.Println("Attempts:", attempts)
	fmt.Println("Min:", min)
	fmt.Println("Max:", max)
	fmt.Println("Seed:", seed)

	dungeon.Initiate().Print()
}
