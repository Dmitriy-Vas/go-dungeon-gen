package main

//#region Imports
import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

//#endregion

//#region Structures
type Rectangle struct {
	X, Y, Width, Height int
}

type Dungeon struct {
	Grid                                         [][]int
	Capacity, NumRooms, Attempts, Min, Max, Seed int
	Rooms                                        []Rectangle
	Regions                                      []int
	Border                                       Rectangle
}

//#endregion

func logError(error error) {
	if error != nil {
		log.Panic(error)
	}
}

func (dung *Dungeon) Initiate() {

}

func main() {
	arguments := os.Args[1:]

	// Initiate default values
	m := make(map[string]int)
	m["size"] = 70
	m["rooms"] = 30
	m["attempts"] = 50
	m["min"] = 5
	m["max"] = 15
	m["seed"] = time.Now().Nanosecond()

	// Parse arguments
	for _, argument := range arguments {
		value := regexp.MustCompile("(\\w+)=(\\d+)").FindAllStringSubmatch(argument, -1)
		if value != nil {
			m[value[0][1]], _ = strconv.Atoi(value[0][2])
		}
	}

	// Get values to generate dungeon
	size := m["size"]
	rooms := m["rooms"]
	attempts := m["attempts"]
	min := m["min"]
	max := m["max"]
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
		Seed:     seed,
		Rooms:    []Rectangle{},
		Regions:  []int{},
		Border:   Rectangle{X: 1, Y: 1, Width: size - 2, Height: size - 2},
	}

	dungeon.Initiate()

	fmt.Println(size, rooms, attempts, min, max, seed)
}
