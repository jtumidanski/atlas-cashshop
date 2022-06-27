package attributes

type CharacterSkillDataContainer struct {
	Data CharacterSkillData `json:"data"`
}

type CharacterSkillDataListContainer struct {
	Data []CharacterSkillData `json:"data"`
}

type CharacterSkillData struct {
	Id         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes CharacterSkillAttributes `json:"attributes"`
}

type CharacterSkillAttributes struct {
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  int64  `json:"expiration"`
}
