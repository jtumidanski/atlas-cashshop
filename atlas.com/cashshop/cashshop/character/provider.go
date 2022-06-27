package character

import (
	"atlas-cashshop/database"
	"atlas-cashshop/model"
	"gorm.io/gorm"
)

func getById(characterId uint32) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{CharacterId: characterId})
	}
}
