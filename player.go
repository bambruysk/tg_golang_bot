package main

import "log"

type PlayerStorager interface {
	Get(name string) Player
	Add(player Player)
	Create(name string) Player
	GetOrCreate(name string) (Player, bool)
}

type Player struct {
	Name    string
	Holdes  []*Holde
	Request HoldeRequest
	World   *HoldeStorage
}

func NewPlayer(name string, world *HoldeStorage) Player {
	return Player{
		Name:    name,
		Holdes:  []*Holde{},
		Request: HoldeRequest{},
		World:   world,
	}
}

func (p *Player) HandleReq() (HoldeResponce, error) {
	if len(p.Request.Holdes) == 0 {
		return HoldeResponce{
			Amount: 0,
		}, nil
	}
	log.Println(p)
	world := p.World
	req := p.Request
	resp, err := world.CalculateHoldes(req)
	
	for _, h := range req.Holdes {
		holde, err := world.Get(h.HoldeID)
		if err != nil {
			panic(err)
		}
		holde.Owner = p.Name
		holde.Amount = 0
	}

	if err != nil {
		return HoldeResponce{}, err
	}
	return HoldeResponce{
		Amount: resp,
	}, nil
}

// Player storage impleneted as in memory storage  In future it will db connection.
type PlayerStorage struct {
	Players map[string]Player
	World   *HoldeStorage
}

// NewPlayerStorage conctructors
func NewPlayerStorage(world *HoldeStorage) PlayerStorage {
	return PlayerStorage{
		Players: map[string]Player{},
		World:   world,
	}
}

// Get player from storage by name
func (ps PlayerStorage) Get(name string) Player {
	return ps.Players[name]
}

// Add player to storage.  If player exist in storage replace it.
func (ps *PlayerStorage) Add(player Player) {
	ps.Players[player.Name] = player
}

func (ps *PlayerStorage) Create(name string) Player {
	ps.Players[name] = NewPlayer(name, ps.World)
	return ps.Players[name]
}

func (ps *PlayerStorage) GetOrCreate(name string) (Player, bool) {
	player, exist := ps.Players[name]
	if !exist {
		return ps.Create(name), true
	}
	return player, false
}

// Update player. If player not exist in storage in will be created
func (ps *PlayerStorage) Update(player Player) {
	ps.Players[player.Name] = player
}
