package character

type Model struct {
	id    uint32
	level byte
}

func (a Model) Id() uint32 {
	return a.id
}

func (a Model) Level() byte {
	return a.level
}
