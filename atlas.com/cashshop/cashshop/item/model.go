package item

type Model struct {
	serialNumber uint32
	itemId       uint32
	price        uint32
	period       uint64
	count        uint16
	onSale       bool
}

func (m Model) SerialNumber() uint32 {
	return m.serialNumber
}

func (m Model) OnSale() bool {
	return m.onSale
}

func (m Model) ItemId() uint32 {
	return m.itemId
}

func (m Model) Price() uint32 {
	return m.price
}
