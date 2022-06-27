package wishlist

import "gorm.io/gorm"

type EntityUpdateFunction func() ([]string, func(e *entity))

func create(db *gorm.DB, characterId uint32, serialNumber uint32) (Model, error) {
	e := &entity{CharacterId: characterId, SerialNumber: serialNumber}
	err := db.Create(e).Error
	if err != nil {
		return Model{}, err
	}
	return makeModel(*e)
}

func makeModel(e entity) (Model, error) {
	return Model{
		id:           e.Id,
		characterId:  e.CharacterId,
		serialNumber: e.SerialNumber,
	}, nil
}

func deleteForCharacter(db *gorm.DB, characterId uint32) error {
	return db.Where(&entity{CharacterId: characterId}).Delete(&entity{}).Error
}
