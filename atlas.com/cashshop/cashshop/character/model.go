package character

import (
	"atlas-cashshop/cashshop/character/wishlist"
)

type Model struct {
	characterId uint32
	credit      uint32
	points      uint32
	prepaid     uint32
	wishlist    []wishlist.Model
}

func (m *Model) Credit() uint32 {
	return m.credit
}

func (m *Model) Points() uint32 {
	return m.points
}

func (m *Model) Prepaid() uint32 {
	return m.prepaid
}

func (m *Model) CharacterId() uint32 {
	return m.characterId
}

func (m *Model) SetWishlist(wishlist []wishlist.Model) Model {
	return Model{
		characterId: m.characterId,
		credit:      m.credit,
		points:      m.points,
		prepaid:     m.prepaid,
		wishlist:    wishlist,
	}
}

func (m *Model) Wishlist() []wishlist.Model {
	return m.wishlist
}

func (m *Model) Cash(cashIndex uint32) uint32 {
	if cashIndex == 1 {
		return m.credit
	}
	if cashIndex == 2 {
		return m.points
	}
	if cashIndex == 4 {
		return m.prepaid
	}
	return 0
}
