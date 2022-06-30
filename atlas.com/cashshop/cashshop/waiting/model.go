package waiting

type Model struct {
	worldId     byte
	channelId   byte
	characterId uint32
	approvals   uint32
}

func (m Model) AddApproval() Model {
	return Model{
		worldId:     m.worldId,
		channelId:   m.channelId,
		characterId: m.characterId,
		approvals:   m.approvals + 1,
	}
}

func (m Model) Approvals() uint32 {
	return m.approvals
}

func (m Model) WorldId() byte {
	return m.worldId
}

func (m Model) ChannelId() byte {
	return m.channelId
}

func (m Model) CharacterId() uint32 {
	return m.characterId
}
