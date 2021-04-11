package main

type PLayerStorager interface {
	Get(name string) Player
	Add(player Player)
	Create(name string) Player
	GetOrCreate(name string) (Player, bool)
}

type Player struct {
	Name    string
	Holdes  []*Holde
	Request HoldeRequest
}

func NewPlayer(name string) Player {
	return Player{
		Name:   name,
		Holdes: []*Holde{},
	}
}

// Player storage mpleneted as in mememory storage/  In futer it will db connection.
type PlayerStorage struct {
	Players map[string]Player
}

// NewPlayerStorage conctructor
func NewPlayerStorage(name string) PlayerStorage {
	return PlayerStorage{
		Players: map[string]Player{},
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
	ps.Players[name] = NewPlayer(name)
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
