package cashshop

import (
	"atlas-cashshop/cashshop/character"
	"atlas-cashshop/cashshop/item"
	character2 "atlas-cashshop/character"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
)

// UseSupplyRateCoupons TODO implement rate coupon configuration
const UseSupplyRateCoupons = true

func EnterCashShop(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) error {
	return func(worldId byte, channelId byte, characterId uint32) error {
		l.Infof("Character %d has attempted to enter the cash shop.", characterId)
		// TODO verify character is allowed to enter cash shop
		//   not using vegas spell
		//   not registered for an event
		//   not in mini dungeon
		//   does not have cash shop already open

		emitEnterCashShop(l, span)(worldId, channelId, characterId)
		return nil
	}
}

func PurchaseCashShopItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, cashIndex uint32, serialNumber uint32) error {
	return func(characterId uint32, cashIndex uint32, serialNumber uint32) error {
		cc, err := character.GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %d purchasing %d.", characterId, serialNumber)
			return err
		}

		i, err := item.GetById(l)(serialNumber)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate item being purchased by character %d, via serial number %d.", characterId, serialNumber)
			return err
		}

		if !i.OnSale() {
			l.WithError(err).Errorf("Character %d attempting to purchase item %d, which is not on sale.", characterId, i.ItemId())
			return errors.New("item not on sale")
		}

		if cc.Cash(cashIndex) < i.Price() {
			l.Debugf("Character %d attempted to purchase item %d without enough NX.", characterId, i.ItemId())
			return errors.New("cannot afford")
		}

		c, err := character2.GetById(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve attributes for character %d.", characterId)
			return err
		}

		if isCashStore(i.ItemId()) && c.Level() < 16 {
			l.Debugf("Character %d denied purchase of a FM market item, because  they are less than level 16.", characterId)
			return errors.New("cannot purchase store")
		}

		if isRateCoupon(i.ItemId()) && !UseSupplyRateCoupons {
			l.Debugf("Character %d denied rate coupon, because they are not currently enabled to purchase.", characterId)
			return errors.New("cannot purchase rate coupon")
		}

		if isMapleLife(i.ItemId()) && c.Level() < 30 {
			l.Debugf("Character %d denied maple life coupon, because they are less than level 30.", characterId)
			return errors.New("cannot purchase maple life")
		}
	}
}

func isMapleLife(itemId uint32) bool {
	itemType := uint32(math.Floor(float64(itemId / 10000)))
	return itemType == 543 && itemId != 5430000
}

func isRateCoupon(itemId uint32) bool {
	itemType := uint32(math.Floor(float64(itemId / 1000)))
	return itemType == 5211 || itemType == 5360
}

func isCashStore(itemId uint32) bool {
	itemType := uint32(math.Floor(float64(itemId / 10000)))
	return itemType == 503 || itemType == 514
}
