package main

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

type NeighbourTestCase struct {
	id  int
	res []int
}

// comapre two slice withoun order
func compareSlices(A, B []int) bool {
	if len(A) != len(B) {
		return false
	}

	sort.Ints(A)
	sort.Ints(B)

	for i, a := range A {
		if a != B[i] {
			return false
		}
	}

	return true
}

// comapre two slice withoun order
func compare2DSlices(A, B []int) bool {
	if len(A) != len(B) {
		return false
	}

	if len(A) == 0 {
		return true // empty slices  is equal
	}

	sort.Ints(A)
	sort.Ints(B)

	for i, a := range A {
		if a != B[i] {
			return false
		}
	}

	return true
}

func TestGetNeighbour(t *testing.T) {
	testCases := []NeighbourTestCase{
		{
			id:  0,
			res: []int{1, 10},
		},
		{
			id:  13,
			res: []int{3, 12, 14, 23},
		},
		{
			id:  9,
			res: []int{8, 19},
		},
		{
			id:  93,
			res: []int{92, 94, 83},
		},
		{
			id:  69,
			res: []int{68, 59, 79},
		},
		{
			id:  99,
			res: []int{98, 89},
		},
	}

	for i, c := range testCases {
		res := getNeighbour(c.id)
		t.Log("Test №", i, c, res)

		if !compareSlices(res, c.res) {
			t.Log("result not compares", c, res)
			t.Fail()
		}

	}
}

type ClusterTestCase struct {
	id       int
	holdes   map[int]bool
	expected []int
}

func TestFindCluster(t *testing.T) {
	testCases := []ClusterTestCase{
		{
			id: 0,
			holdes: map[int]bool{
				0:  true,
				1:  true,
				10: true,
			},
			expected: []int{0, 1, 10},
		},
		{
			id: 13,
			holdes: map[int]bool{
				13: true,
				14: true,
				15: true,
				24: true,
				4:  true,
				65: true,
				99: true,
			},
			expected: []int{13, 14, 15, 24, 4},
		},
		{
			id: 99,
			holdes: map[int]bool{
				99: true,
				1:  true,
				10: true,
			},
			expected: []int{99},
		},
	}
	visited := make(map[int]bool)

	for _, c := range testCases {
		res := make([]int, 0)
		findCluster(c.id, &res, c.holdes, visited)
		fmt.Println(res, c.holdes, visited)
		t.Log("Test ", res)
		if !compareSlices(res, c.expected) {
			t.Log("result not compares", c, res)
			t.Fail()
		}
	}
}

type ManyClustersTestCase struct {
	holdes   []int
	expected [][]int
}

func TestCalculateClusters(t *testing.T) {
	testCases := []ManyClustersTestCase{
		{
			holdes: []int{0, 1, 3, 5, 6, 8},
			expected: [][]int{
				{0},
				{1},
				{3},
				{5, 6},
				{8},
			},
		},
		{
			holdes: []int{0, 13, 1, 3, 5, 6, 8, 14, 24, 25, 35, 34, 33, 23},
			expected: [][]int{
				{0},
				{1},
				{3},
				{5, 6},
				{8},
			},
		},
	}

	for _, c := range testCases {
		res := calculateClusters(c.holdes)
		fmt.Println(res, c.holdes)
		t.Log("Test ", res)
		// if !compareSlices(res, c.expected) {
		// 	t.Log("result not compares", c, res)
		// 	t.Fail()
		// }
	}

	t.Fail()

}

func TestResponseText(t *testing.T) {
	holde := Holde{
		Name:      "My Holde",
		ID:        15,
		Amount:    12,
		Level:     5,
		Owner:     "Bambr",
		LastVisit: time.Now(),
	}
	t.Log(holde.LastVisitString())

	t.Log(holde.ResponseText())
	t.Fail()
}
