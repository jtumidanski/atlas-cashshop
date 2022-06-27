package wishlist

import (
	"atlas-cashshop/cashshop/item"
	"atlas-cashshop/database"
	"atlas-cashshop/model"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ByIdModelProvider(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) model.SliceProvider[Model] {
	return func(characterId uint32) model.SliceProvider[Model] {
		return database.ModelSliceProvider[Model, entity](db)(getById(characterId), makeModel)
	}
}

func GetById(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		return ByIdModelProvider(l, db)(characterId)()
	}
}

func DeleteForCharacter(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32) error {
	return func(characterId uint32) error {
		err := deleteForCharacter(db, characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to clear wishlist for character %d.", characterId)
			return err
		}
		return nil
	}
}

func ModifyWishlist(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, serialNumbers []uint32) error {
	return func(characterId uint32, serialNumbers []uint32) error {
		txError := db.Transaction(func(tx *gorm.DB) error {
			err := DeleteForCharacter(l, db, span)(characterId)
			if err != nil {
				return err
			}

			for _, sn := range serialNumbers {
				if sn == 0 {
					continue
				}

				var i item.Model
				i, err = item.GetById(l)(sn)
				if err != nil {
					l.WithError(err).Warnf("Serial number %d provided to character %d wishlist modification is invalid.", sn, characterId)
					continue
				}
				if !i.OnSale() {
					l.Warnf("Serial number %d added to character %d wishlist is not on sale.", sn, characterId)
				}

				_, err = create(db, characterId, sn)
				if err != nil {
					l.WithError(err).Errorf("Unable to create wishlist item %d for character %d.", sn, characterId)
					return err
				}
			}
			return nil
		})

		if txError == nil {
			emitWishlistUpdateEvent(l, span)(characterId)
		}
		return txError
	}
}
