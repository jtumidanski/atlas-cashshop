package character

import (
	"atlas-cashshop/cashshop/character/wishlist"
	"atlas-cashshop/database"
	"atlas-cashshop/model"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func byIdModelProvider(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) model.Provider[Model] {
	return func(characterId uint32) model.Provider[Model] {
		mp := database.ModelProvider[Model, entity](db)(getById(characterId), makeModel)
		wp := wishlist.ByIdModelProvider(l, db)(characterId)
		return func() (Model, error) {
			m, err := mp()
			if err != nil {
				return Model{}, err
			}

			wm, err := wp()
			if err != nil {
				return m, err
			}

			return m.SetWishlist(wm), nil
		}
	}
}

func GetById(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		return byIdModelProvider(l, db)(characterId)()
	}
}

func Created(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32) error {
	return func(characterId uint32) error {
		l.Debugf("Processing a character created event for %d. Creating in local database.", characterId)
		_, err := create(db, characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to create character %d.", characterId)
		}
		return err
	}
}

func AwardCredit(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32, amount uint32) error {
	return func(characterId uint32, amount uint32) error {
		c, err := GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Cannot award credit to character %d who does not exist.", characterId)
			return err
		}
		return update(db, characterId, SetCredit(c.Credit()+amount))
	}
}

func AwardPoints(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32, amount uint32) error {
	return func(characterId uint32, amount uint32) error {
		c, err := GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Cannot award points to character %d who does not exist.", characterId)
			return err
		}
		return update(db, characterId, SetPoints(c.Points()+amount))
	}
}

func AwardPrepaid(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32, amount uint32) error {
	return func(characterId uint32, amount uint32) error {
		c, err := GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Cannot award prepaid nx to character %d who does not exist.", characterId)
			return err
		}
		return update(db, characterId, SetPrepaid(c.Points()+amount))
	}
}
