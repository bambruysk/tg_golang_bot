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

func calculateClusters(ids []int) [][]int {
	// holde_map := make([][]bool, WorldSizeY)
	// for i, _ := range holde_map {
	// 	holde_map[i] = make([]bool, WorldSizeX)
	// }

	// for id := range ids {
	// 	holde_map[id/WorldSizeX][id/WorldSizeY] = true
	// }

	holde_map := make([]bool, HoldesNumber)
	visited := make(map[int]bool)
	for id := range ids {
		holde_map[id] = true
	}

	res := make([][]int, 0)

	cnt := len(ids)

	dfs := func(id int) {

		_, ok := visited[id]
		if ok {
			return
		}
		visited[id] = true

		if id > WorldSizeX && holde_map[id-WorldSizeX] {

		}

	}

}
