package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func logError(error error) {
	if error != nil {
		log.Panic(error)
	}
}

func main() {
	arguments := os.Args[1:]

	m := make(map[string]int)
	m["size"] = 100
	m["rooms"] = 200

	for _, argument := range arguments {
		value := regexp.MustCompile("size=(\\d+)").FindAllStringSubmatch(argument, -1)
		if value != nil {
			m["size"], _ = strconv.Atoi(value[0][1])
		}

		value = regexp.MustCompile("rooms=(\\d+)").FindAllStringSubmatch(argument, -1)
		if value != nil {
			m["rooms"], _ = strconv.Atoi(value[0][1])
		}
	}

	size := m["size"]
	rooms := m["rooms"]

	fmt.Println(size, rooms)
}
