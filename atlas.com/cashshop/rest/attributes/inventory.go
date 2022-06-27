package attributes

type InventoryListDataContainer struct {
	Data []InventoryData `json:"data"`
}

type InventoryDataContainer struct {
	Data     InventoryData `json:"data"`
	Included []interface{} `json:"included"`
}

type InventoryData struct {
	Id            string              `json:"id"`
	Type          string              `json:"type"`
	Attributes    InventoryAttributes `json:"attributes"`
	Relationships Relationships       `json:"relationships"`
}

type InventoryAttributes struct {
	Type     string `json:"type"`
	Capacity uint32 `json:"capacity"`
}

type Relationships struct {
	InventoryItems      []Relationship `json:"inventoryItems"`
	EquipmentStatistics []Relationship `json:"equipmentStatistics"`
}

type Relationship struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type ItemDataContainer struct {
	Data InventoryItemData `json:"data"`
}

type InventoryItemData struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes InventoryItemAttributes `json:"attributes"`
}

type InventoryItemAttributes struct {
	ItemId   uint32 `json:"itemId"`
	Quantity int32  `json:"quantity"`
	Slot     int16  `json:"slot"`
}

type InventoryEquipmentData struct {
	Id            string                       `json:"id"`
	Type          string                       `json:"type"`
	Attributes    InventoryEquipmentAttributes `json:"attributes"`
	Relationships Relationships                `json:"relationships"`
}

type InventoryEquipmentAttributes struct {
	EquipmentId uint32 `json:"equipmentId"`
	Slot        int16  `json:"slot"`
}

type InventoryEquipmentStatisticsData struct {
	Id         string                                 `json:"id"`
	Type       string                                 `json:"type"`
	Attributes InventoryEquipmentStatisticsAttributes `json:"attributes"`
}

type InventoryEquipmentStatisticsAttributes struct {
	ItemId        uint32 `json:"itemId"`
	Strength      uint16 `json:"strength"`
	Dexterity     uint16 `json:"dexterity"`
	Intelligence  uint16 `json:"intelligence"`
	Luck          uint16 `json:"luck"`
	HP            uint16 `json:"hp"`
	MP            uint16 `json:"mp"`
	WeaponAttack  uint16 `json:"weaponAttack"`
	MagicAttack   uint16 `json:"magicAttack"`
	WeaponDefense uint16 `json:"weaponDefense"`
	MagicDefense  uint16 `json:"magicDefense"`
	Accuracy      uint16 `json:"accuracy"`
	Avoidability  uint16 `json:"avoidability"`
	Hands         uint16 `json:"hands"`
	Speed         uint16 `json:"speed"`
	Jump          uint16 `json:"jump"`
	Slots         uint16 `json:"slots"`
}
