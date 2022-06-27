package wishlist

type Model struct {
	id           uint32
	characterId  uint32
	serialNumber uint32
}

func (m *Model) Id() uint32 {
	return m.id
}

func (m *Model) CharacterId() uint32 {
	return m.characterId
}

func (m *Model) SerialNumber() uint32 {
	return m.serialNumber
}
