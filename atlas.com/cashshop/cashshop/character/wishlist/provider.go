package wishlist

import (
	"atlas-cashshop/database"
	"atlas-cashshop/model"
	"gorm.io/gorm"
)

func getById(characterId uint32) database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		return database.SliceQuery[entity](db, &entity{CharacterId: characterId})
	}
}
