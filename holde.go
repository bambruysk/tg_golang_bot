package main

import (
	"bytes"
	"errors"
	"math"
	"strconv"
	"text/template"
	"time"
)

// type for money represnt
type Money float64

// Round money to for palyer payment
func (m Money) Round() int {
	return int(math.Round(float64(m)))
}

func (m Money) String() string {
	return strconv.Itoa(m.Round())
}

type HoldeGameSettings struct {
	// Базовый доход на уровень в час
	MoneyPerHour float64

	// Коэфифициент синергии
	SynergyCoeff float64

	// Time degradation coeffs
	TimeDegradation float64
}

// Holde - basic structures  for holde
type Holde struct {
	// Имя поместья
	Name string
	// ID -  номер поместья
	ID int
	// Money - текущее накопелния
	Amount Money

	// Уровень
	Level int

	// Владелец
	Owner string

	//
	LastVisit time.Time
}

func (h *Holde) LastVisitString() string {
	return h.LastVisit.Format(time.Stamp)
}

const holdeDescription = `
Поместье №{{.ID}} {{.Name}}

Количество денег {{.Amount}}
Уровень {{.Level}}
Владелец {{.Owner}}
Время последнего посещения {{.LastVisit}}
`

func (h Holde) ResponseText() string {
	var buf bytes.Buffer

	tmpl := template.Must(template.New("HoldeDescription").Parse(holdeDescription))
	tmpl.Execute(&buf, h)
	return buf.String()
}

// Request

type HoldeRequestItem struct {
	HoldeID int
	Dice    int
}

type HoldeRequest struct {
	Holdes []HoldeRequestItem
	Owner  string
	User   string
}

type HoldeResponce struct {
	Amount Money
}

func (r *HoldeRequest) Calculate() HoldeResponce {

}

// we not use decimal type for money. Money round to ceil, for players fun

var WorldSizeX = 10
var WorldSizeY = 10
var HoldesNumber = WorldSizeX * WorldSizeY

var Settings = HoldeGameSettings{
	MoneyPerHour:    2,
	SynergyCoeff:    1,
	TimeDegradation: 0.2,
}

func (s HoldeGameSettings) GetSynergy(num int) float64 {
	return float64(num) * s.SynergyCoeff
}

type HoldeStorage map[int]*Holde

func NewHoldeStorage() HoldeStorage {
	//TODO: load from json or db
	holdes := make(HoldeStorage)
	init_time := time.Now()
	for i := 0; i < HoldesNumber; i++ {
		holdes[i] = &Holde{
			Name:      "",
			ID:        i,
			Amount:    0,
			Level:     1,
			Owner:     "",
			LastVisit: init_time,
		}
	}
	return holdes
}

func (hs HoldeStorage) Get(id int) (*Holde, error) {
	if id < 0 || id >= HoldesNumber {
		return nil, errors.New("Holde not found in storage")
	}
	return hs[id], nil
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

func (h *Holde) Visit(dice int) Money {

	//  time from last visit
	hours := int(h.LastVisit.Sub(time.Now()).Hours())
	// calculate money
	money := Money(0)
	for h := 0; h < hours; h++ {
		money += Money(Settings.MoneyPerHour * (1 - float64(h)*Settings.TimeDegradation))
	}
	h.LastVisit = time.Now()

	return money * Money(dice/5) // D10
}

func (hs HoldeStorage) CalculateHoldes(holdes []int) (Money, error) {
	// Get holdes from storage
	money := Money(0)
	for _, h := range holdes {
		hold, err := hs.Get(h)
		if err != nil {
			return 0, err
		}
		// 4 random number ;)
		money += hold.Visit(4)
	}
	//
	clusters := calculateClusters(holdes)
	for _, c := range clusters {
		if len(c) > 1 {
			// add cluster bonus example
			money += Money(len(c) * int(Settings.SynergyCoeff))
		}
	}
	return money, nil
}
