package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
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
	// HoldeNums
	HoldeNums int
	// Размер мира по горизонтали
	WorldSizeX int
	// Размер мира по вертикали
	WorldSizeY int
	// Базовый доход на уровень в час
	MoneyPerHour float64
	// Коэфифициент синергии
	SynergyCoeff float64
	// Time degradation coeffs
	TimeDegradation float64
	// Locations
	Locations []string
	// Holde uprade cost byt level 1->2, 2->3, etc
	HoldeLevelUpgradeCost []Money
	// Holde maximum level
	HoldeMaxLevel int
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

func (hr *HoldeResponce) Show() string {
	return fmt.Sprintf("Выдай этому игроку %v монет", hr.Amount)
}

// we not use decimal type for money. Money round to ceil, for players fun

var WorldSizeX = gameSettings.WorldSizeX
var WorldSizeY = gameSettings.WorldSizeY
var HoldesNumber = WorldSizeX * WorldSizeY

// var Settings = HoldeGameSettings{
// 	MoneyPerHour:    2,
// 	SynergyCoeff:    1,
// 	TimeDegradation: 0.2,
// }

func (s HoldeGameSettings) GetSynergy(num int) float64 {
	return float64(num) * s.SynergyCoeff
}

type HoldeStorage struct {
	holdes map[int]*Holde
}

func NewHoldeStorage() *HoldeStorage {
	//TODO: load from json or db
	hs := HoldeStorage{
		holdes: make(map[int]*Holde),
	}
	//init_time := time.Now()
	holde_list := gspread.ReadHoldes()
	for i := 0; i < HoldesNumber; i++ {
		hs.holdes[i] = &holde_list[i]
	}

	// TODO: Add save to db!

	// for debug purpose only
	r := rand.Perm(HoldesNumber)
	for i := 0; i < HoldesNumber; i++ {
		hs.holdes[i].Amount = Money(r[i])
	}

	return &hs
}

func (hs *HoldeStorage) Get(id int) (*Holde, error) {
	if id < 0 || id >= HoldesNumber {
		return nil, errors.New("Holde not found in storage")
	}
	return hs.holdes[id], nil
}

func (hs *HoldeStorage) Update(holde *Holde) {
	hs.holdes[holde.ID] = holde
}

func (hs *HoldeStorage) String() string {
	res := "\n"
	for _, h := range hs.holdes {
		res += fmt.Sprintf("Holde : %d \t %s \t %d \n", h.ID, h.Name, h.Level)
	}
	return res
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
	money := h.Amount
	for h := 0; h < hours; h++ {
		money += Money(gameSettings.MoneyPerHour * (1 - float64(h)*gameSettings.TimeDegradation))
	}
	h.LastVisit = time.Now()

	return money * Money(dice/5) // D10
}

func (hs HoldeStorage) CalculateHoldes(req HoldeRequest) (Money, error) {

	// Get holdes from storage
	money := Money(0)

	for _, rh := range req.Holdes {
		hold, err := hs.Get(rh.HoldeID)
		if err != nil {
			return 0, err
		}

		// 4 random number ;)
		money += hold.Visit(rh.Dice)
	}
	//
	holdeNums := make([]int, len(req.Holdes))
	for i, rh := range req.Holdes {
		holdeNums[i] = rh.HoldeID
	}

	clusters := calculateClusters(holdeNums)
	for _, c := range clusters {
		if len(c) > 1 {
			// add cluster bonus example
			money += Money(len(c) * int(gameSettings.SynergyCoeff))
		}
	}
	return money, nil
}

func (h *Holde) Upgrade() bool {
	if h.Level == gameSettings.HoldeMaxLevel {
		return false
	}
	h.Level++
	return true
}
