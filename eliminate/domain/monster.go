package domain

type Monster struct {
	refID string
	x     int
	y     int
	order int
}

func (monster Monster) getX() int {
	return monster.x
}

func (monster Monster) getY() int {
	return monster.y
}
