package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sort"
)

var paths [][]string

type room struct {
	id      int
	name    string
	x       int
	y       int
	friends []string
}

type traffic struct {
	len int
}

type ant struct {
	id       int
	path     []string
	position int
	finished bool
}

func main() {
	lines := parseFile()
	antNum, rooms, start, end := getData(lines)
	getPaths(rooms, start, end)
	Ants := trafficAnts(antNum, paths, start)
	printAnts(Ants, start, end)
}

func printAnts(Ants map[int]ant, start, end string) {
	isBusy := make(map[string]bool)
	finishedCount := len(Ants)
	keys := make([]int, len(Ants))
	n := 0
	for k := range Ants {
		keys[n] = k
		n++ 
	}
	sort.Ints(keys)

	//fmt.Println(Ants)
	for i := 0; i < finishedCount; {
		for _, j := range keys {
			var v ant
			v = Ants[j]
			if !Ants[j].finished {
				if !isBusy[Ants[j].path[Ants[j].position + 1]] {
					if Ants[j].path[Ants[j].position + 1] == end {
						v.finished = true
						Ants[v.id] = v
						i++
						fmt.Printf("L%v-%s ", Ants[j].id, end)
						isBusy[Ants[j].path[Ants[j].position]] = false 
						continue
					}
					v.position++
					Ants[Ants[j].id] = v
					fmt.Printf("L%v-%s ", Ants[j].id, Ants[j].path[Ants[j].position])
					isBusy[Ants[j].path[Ants[j].position]] = true
					isBusy[Ants[j].path[Ants[j].position - 1]] = false
				}
			}
		}
		fmt.Println()
	}
	
	

}

func trafficAnts(antNum int, paths [][]string, start string) map[int]ant {
	Ants := make(map[int]ant)
	for i := 0; i < antNum; i++ {
		var currentAnt ant
		currentAnt.id = i + 1
		currentAnt.position = 0
		currentAnt.finished = false
		Ants[currentAnt.id] = currentAnt
	}
	index := 1
	if len(paths) == 1 {
		for range Ants {
			var currentAnt ant
			currentAnt.id = index
			currentAnt.position = 0
			currentAnt.finished = false
			currentAnt.path = paths[0]
			Ants[currentAnt.id] = currentAnt
			index++
		}
		return Ants
	}
	var copyPaths [][]string
	for _, v := range paths {
		copyPaths = append(copyPaths, v)
	}
	index = 1
	for i := 0; i < antNum; {
		for i := range copyPaths {
			for j := range copyPaths {
				if i != j && antNum > 0 {
					if len(copyPaths[i]) < len(copyPaths[j]) || len(copyPaths[i]) == len(copyPaths[j]) {
						copyPaths[i] = append(copyPaths[i], "*")
						var currentAnt ant
						currentAnt.id = index
						currentAnt.path = paths[i]
						Ants[currentAnt.id] = currentAnt
						antNum--
						index++
						continue
					}
					if len(copyPaths[i]) > len(copyPaths[j]) {
						copyPaths[j] = append(copyPaths[j], "*")
						var currentAnt ant
						currentAnt.id = index
						currentAnt.path = paths[j]
						Ants[currentAnt.id] = currentAnt
						index++
						antNum--
					}
				}
			}
		}
	}
	return Ants

}

func getPaths(rooms map[string]room, start, end string) {
	var path []string
	path = append(path, start)
	visited := make(map[string]bool)
	DFS(rooms, visited, path, start, end)
	paths = getUniquePaths(paths)
}

func getUniquePaths(paths [][]string) [][]string {
	if len(paths) == 0 {
		fmt.Println("ERROR: validPaths not found")
		os.Exit(1)
	}
	if len(paths) == 1 {
		return paths
	}
	var validPaths [][]string
	//BubbleSort for sorting paths in array
	for i := 0; i < len(paths); i++ {
		for j := i; j < len(paths); j++ {
			if len(paths[i]) > len(paths[j]) {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}

	var visited []string
	count := 0
	for i, path := range paths {
		count = 0
		path = path[1 : len(path)-1]
		for j := range path {
			for _, validPath := range validPaths {
				validPath = validPath[1 : len(validPath)-1]
				for l := range validPath {
					visited = append(visited, validPath[l])
				}
			}

			if !isExist(visited, path[j]) {
				count++
			}
		}
		if count == len(path) {
			validPaths = append(validPaths, paths[i])
		}
	}
	return validPaths
}

// DFS function
func DFS(rooms map[string]room, visited map[string]bool, path []string, currentRoom, endRoom string) {
	visited[currentRoom] = true
	if rooms[currentRoom].name == endRoom {
		newPath := make([]string, len(path))
		for i, v := range path {
			newPath[i] = v
		}
		paths = append(paths, newPath)
		visited[currentRoom] = false
		return
	}

	for _, v := range rooms[currentRoom].friends {
		if !visited[v] {
			path = append(path, v)
			DFS(rooms, visited, path, v, endRoom)
			path = path[:len(path)-1]
		}
	}
	visited[currentRoom] = false
}

//parseFile get command line arguments, read text file and return array
func parseFile() []string {
	var lines []string
	if len(os.Args) < 2 {
		log.Println("ERROR: Please input fileName")
		os.Exit(1)
	}
	fileName := os.Args[1]
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		fmt.Println(scanner.Text())
	}
	fmt.Println()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func getData(lines []string) (int, map[string]room, string, string) {

	rooms := make(map[string]room)

	var (
		index  int
		antNum int
		links  []string
		start  string
		end    string
	)

	antNum, err := strconv.Atoi(lines[0])
	if err != nil {
		log.Fatal(err)
	}
	if antNum <= 0 {
		log.Println("ERROR: Incorrect number of ants")
		os.Exit(1)
	}

	for i, line := range lines {
		var r room
		splitRoomsLine := strings.Split(line, " ")
		if len(splitRoomsLine) == 3 {
			if strings.HasPrefix(splitRoomsLine[0], "L") || strings.HasPrefix(splitRoomsLine[0], "#") {
				log.Println("ERROR: A room will never start with the letter L or with # and must have no spaces")
				os.Exit(1)
			}
			x, err := strconv.Atoi(splitRoomsLine[1])
			if err != nil {
				log.Printf("ERROR: The coordinates of the rooms will always be int.\n %v\n", err)
				os.Exit(1)
			}
			y, err := strconv.Atoi(splitRoomsLine[2])
			if err != nil {
				log.Printf("ERROR: The coordinates of the rooms will always be int.\n %v\n", err)
				os.Exit(1)
			}

			if _, found := rooms[splitRoomsLine[0]]; found {
				log.Println("ERROR: This room already exists")
				os.Exit(1)
			}

			r.id = index
			r.name = splitRoomsLine[0]
			r.x = x
			r.y = y
			rooms[r.name] = r
			index++

		}
		splitLinksLine := strings.Split(line, "-")
		if len(splitLinksLine) == 2 {
			if splitLinksLine[0] == splitLinksLine[1] {
				log.Println("ERROR: Wrong link")
				os.Exit(1)
			}
			if !isExist(links, line) {
				links = append(links, line)
			}
		}
		if line == "##start" && i+2 < len(lines) {
			splitStart := strings.Split(lines[i+1], " ")
			start = splitStart[0]
		}
		if line == "##end" && i+2 < len(lines) {
			splitEnd := strings.Split(lines[i+1], " ")
			end = splitEnd[0]
		}

	}
	if start == "" || end == "" {
		log.Println("ERROR: Start or end room not found")
		os.Exit(1)
	}
	if start == end {
		log.Println("ERROR: a room can't point to itself")
		os.Exit(1)
	}
	for _, v := range rooms {
		r := v
		for _, line := range links {
			splitLinks := strings.Split(line, "-")
			if splitLinks[0] == r.name {
				r.friends = append(r.friends, splitLinks[1])
			}
			if splitLinks[1] == r.name {
				r.friends = append(r.friends, splitLinks[0])
			}
		}
		rooms[v.name] = r
	}
	return antNum, rooms, start, end
}

func isExist(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
