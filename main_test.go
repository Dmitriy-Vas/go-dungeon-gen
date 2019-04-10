package main

import (
	"math/rand"
	"testing"
)

func TestMax(t *testing.T) {
	a := 1
	b := 2
	if res := Max(a, b); res == a {
		t.Fatal(res)
	}
}

func TestMin(t *testing.T) {
	a := 1
	b := 2
	if res := Min(a, b); res == b {
		t.Fatal(res)
	}
}

func TestRectangle_Contains(t *testing.T) {
	r := &Rectangle{
		X:      1,
		Y:      1,
		Width:  50,
		Height: 50,
	}
	p := Point{
		X: 5,
		Y: 5,
	}
	if res := r.Contains(p); res == false {
		t.Fatal(res)
	}
}

func TestRectangle_Intersects(t *testing.T) {
	r := &Rectangle{
		X:      1,
		Y:      1,
		Width:  50,
		Height: 50,
	}
	r2 := &Rectangle{
		X:      49,
		Y:      49,
		Width:  5,
		Height: 5,
	}
	if res := r.Intersects(r2); res == false {
		t.Fatal(res)
	}
}

func TestGenerator(t *testing.T) {
	size := 30

	grid := make([][]int, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]int, size)
	}

	dungeon := &Dungeon{
		Grid:     grid,
		Capacity: size,
		NumRooms: 5,
		Attempts: 50,
		Min:      1,
		Max:      30,
		Rooms:    []Rectangle{},
		Regions:  []int{},
		Border:   Rectangle{X: 1, Y: 1, Width: size - 2, Height: size - 2},
		Seed:     rand.New(rand.NewSource(15)),
	}

	dung := dungeon.Initiate()
	if dung.Grid[1][1] != 1 || dung.Grid[7][6] != 0 {
		t.Fatal(dung.Grid)
	}
}
