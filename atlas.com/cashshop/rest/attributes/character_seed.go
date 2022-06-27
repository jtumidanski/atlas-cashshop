package attributes

type CharacterSeedInputDataContainer struct {
	Data CharacterSeedData `json:"data"`
}

type CharacterSeedData struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes CharacterSeedAttributes `json:"attributes"`
}

type CharacterSeedAttributes struct {
	AccountId uint32 `json:"accountId"`
	WorldId   byte   `json:"worldId"`
	Name      string `json:"name"`
	JobIndex  uint32 `json:"jobIndex"`
	Face      uint32 `json:"face"`
	Hair      uint32 `json:"hair"`
	HairColor uint32 `json:"hairColor"`
	Skin      byte   `json:"skin"`
	Gender    byte   `json:"gender"`
	Top       uint32 `json:"top"`
	Bottom    uint32 `json:"bottom"`
	Shoes     uint32 `json:"shoes"`
	Weapon    uint32 `json:"weapon"`
}
