package character

import "gorm.io/gorm"

type EntityUpdateFunction func() ([]string, func(e *entity))

func create(db *gorm.DB, characterId uint32) (Model, error) {
	e := &entity{CharacterId: characterId}
	err := db.Create(e).Error
	if err != nil {
		return Model{}, err
	}
	return makeModel(*e)
}

func makeModel(e entity) (Model, error) {
	return Model{
		characterId: e.CharacterId,
		credit:      e.Credit,
		points:      e.Points,
		prepaid:     e.Prepaid,
	}, nil
}

func update(db *gorm.DB, characterId uint32, modifiers ...EntityUpdateFunction) error {
	e := &entity{}

	var columns []string
	for _, modifier := range modifiers {
		c, u := modifier()
		columns = append(columns, c...)
		u(e)
	}
	return db.Model(&entity{CharacterId: characterId}).Select(columns).Updates(e).Error
}

func SetCredit(amount uint32) EntityUpdateFunction {
	return func() ([]string, func(e *entity)) {
		return []string{"Credit"}, func(e *entity) {
			e.Credit = amount
		}
	}
}

func SetPoints(amount uint32) EntityUpdateFunction {
	return func() ([]string, func(e *entity)) {
		return []string{"Points"}, func(e *entity) {
			e.Points = amount
		}
	}
}

func SetPrepaid(amount uint32) EntityUpdateFunction {
	return func() ([]string, func(e *entity)) {
		return []string{"Prepaid"}, func(e *entity) {
			e.Prepaid = amount
		}
	}
}
