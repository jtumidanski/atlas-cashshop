package character

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&entity{})
}

type entity struct {
	CharacterId uint32 `gorm:"primaryKey;not null"`
	Credit      uint32 `gorm:"not null;default=0"`
	Points      uint32 `gorm:"not null;default=0"`
	Prepaid     uint32 `gorm:"not null;default=0"`
}

func (e entity) TableName() string {
	return "characters"
}
