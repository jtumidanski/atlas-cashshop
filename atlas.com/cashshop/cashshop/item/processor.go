package item

import (
	"atlas-cashshop/model"
	"github.com/sirupsen/logrus"
)

func byIdModelProvider(_ logrus.FieldLogger) func(serialNumber uint32) model.Provider[Model] {
	return func(serialNumber uint32) model.Provider[Model] {
		return func() (Model, error) {
			return GetCache().GetItem(serialNumber)
		}
	}
}

func GetById(l logrus.FieldLogger) func(serialNumber uint32) (Model, error) {
	return func(serialNumber uint32) (Model, error) {
		return byIdModelProvider(l)(serialNumber)()
	}
}
