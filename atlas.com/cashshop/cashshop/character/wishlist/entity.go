package wishlist

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&entity{})
}

type entity struct {
	Id           uint32 `gorm:"primaryKey;autoIncrement;not null"`
	CharacterId  uint32 `gorm:"not null"`
	SerialNumber uint32 `gorm:"not null"`
}

func (e entity) TableName() string {
	return "wishlist_items"
}
