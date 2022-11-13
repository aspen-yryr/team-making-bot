package player

type Player struct {
	ID   string
	Name string
}

func NewPlayer(id, name string) *Player {
	return &Player{
		ID:   id,
		Name: name,
	}
}

func (p *Player) Is(p_ *Player) bool {
	if p_.ID == p.ID && p_.Name == p.Name {
		return true
	}
	return false
}
