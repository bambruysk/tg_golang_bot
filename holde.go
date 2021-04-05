package main

type Holde struct {
	Name  string
	ID    int
	Money float64
}

// we not use decimal type for money. Money round to ceil, for players fun

var WorldSizeX = 10
var WorldSizeY = 10
var HoldesNumber = WorldSizeX * WorldSizeY

type HoldeStorage map[int]Holde

func NewHoldeStorage() HoldeStorage {
	//TODO: load from json or db
	holdes := make(HoldeStorage)
	return holdes
}

func getNeighbour(id int) []int {
	if id < 0 || id >= HoldesNumber {
		return nil
	}
	res := make([]int, 0)

	// 0  1  2  3  4  5  6  7  8  9
	if id >= WorldSizeX {
		res = append(res, id-WorldSizeX)
	}

	// 90 91 92 93 94 95 96 97 98 99
	if id < HoldesNumber-WorldSizeX { // 100 - 10
		res = append(res, id+WorldSizeX)
	}

	//  00  10  20  30  40  50  60  70  80  90
	if id%WorldSizeX != 0 {
		res = append(res, id-1)
	}
	// 09 19 29 39 49 59 69 79 89 99
	if id%WorldSizeX != WorldSizeX-1 {
		res = append(res, id+1)
	}

	return res
}

func findCluster(id int, res *[]int, holdes, visited map[int]bool) {
	_, ok := holdes[id]
	if !ok {
		return
	}
	_, ok = visited[id]
	if ok {
		return
	}
	visited[id] = true

	*res = append(*res, id)
	for _, n := range getNeighbour(id) {
		findCluster(n, res, holdes, visited)
	}

}

func calculateClusters(ids []int) [][]int {

	holdeMap := make(map[int]bool)
	visited := make(map[int]bool)
	for _, id := range ids {
		holdeMap[id] = true
	}

	res := make([][]int, 0)

	for holde := range holdeMap {
		r := make([]int, 0)
		findCluster(holde, &r, holdeMap, visited)
		if len(r) > 0 {
			res = append(res, r)
		}
	}

	return res

}
