package attributes

type CharacterDamageDataContainer struct {
	Data CharacterDamageData `json:"data"`
}

type CharacterDamageDataListContainer struct {
	Data []CharacterDamageData `json:"data"`
}

type CharacterDamageData struct {
	Id         string                    `json:"id"`
	Type       string                    `json:"type"`
	Attributes CharacterDamageAttributes `json:"attributes"`
}

type CharacterDamageAttributes struct {
	Type    string `json:"type"`
	Maximum uint32 `json:"maximum"`
}
