package cashshop

import (
	"atlas-cashshop/cashshop/character"
	"atlas-cashshop/cashshop/gatekeeper"
	"atlas-cashshop/cashshop/item"
	"atlas-cashshop/cashshop/waiting"
	character2 "atlas-cashshop/character"
	"atlas-cashshop/model"
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

		if gatekeeper.GetRegistry().Count() == 0 {
			// Shortcut gatekeeper process, as no one is participating.
			emitEnterCashShop(l, span)(worldId, channelId, characterId)
			return nil
		}

		// Poll cash shop gatekeepers to see if anyone is interested in denying entry.
		err := waiting.GetRegistry().Add(worldId, channelId, characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to add character %d to waitlist for cash shop entry.", characterId)
			return err
		}
		pollCashShopEntry(l, span)(worldId, channelId, characterId)

		// TODO verify character is allowed to enter cash shop
		//   not using vegas spell
		//   not registered for an event
		//   does not have cash shop already open

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
		return nil
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

func RegisterGatekeeper(l logrus.FieldLogger, _ opentracing.Span) func(service string) error {
	return func(service string) error {
		l.Debugf("Service %s has registered to be a gatekeeper for cash shop entry.", service)
		gatekeeper.GetRegistry().Register(service)
		return nil
	}
}

func UnregisterGatekeeper(l logrus.FieldLogger, span opentracing.Span) func(service string) error {
	return func(service string) error {
		l.Debugf("Service %s has unregistered to be a gatekeeper for cash shop entry.", service)
		gatekeeper.GetRegistry().Unregister(service)

		err := waiting.GetRegistry().ProcessAllApproved(gatekeeper.GetRegistry().Count(), emitEnters(l, span))
		if err != nil {
			l.WithError(err).Errorf("Error issuing cash shop entry notifications.")
			return err
		}
		return nil
	}
}

func emitEnters(l logrus.FieldLogger, span opentracing.Span) func(m []waiting.Model) error {
	return model.ExecuteForEach[waiting.Model](emitEnter(l, span))
}

func emitEnter(l logrus.FieldLogger, span opentracing.Span) model.Operator[waiting.Model] {
	return func(m waiting.Model) error {
		emitEnterCashShop(l, span)(m.WorldId(), m.ChannelId(), m.CharacterId())
		return nil
	}
}

func GatekeeperApproval(l logrus.FieldLogger, span opentracing.Span) func(service string, characterId uint32, messageType string, message string) error {
	return func(service string, characterId uint32, messageType string, message string) error {
		l.Debugf("Service %s approved entry to cash shop for character %d.", service, characterId)
		err := waiting.GetRegistry().AddApproval(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to track character %d entry to cash shop via service %s.", characterId, service)
			return err
		}

		// Only error condition here is if the character has been rejected by another gatekeeper.
		_ = waiting.GetRegistry().ProcessIfApproved(characterId, gatekeeper.GetRegistry().Count(), emitEnter(l, span))
		return nil
	}
}

func GatekeeperDenial(l logrus.FieldLogger, span opentracing.Span) func(service string, characterId uint32, messageType string, message string) error {
	return func(service string, characterId uint32, messageType string, message string) error {
		l.Debugf("Service %s denied entry to cash shop for character %d.", service, characterId)

		// Only error condition here is if the character has been rejected by another gatekeeper.
		m, _ := waiting.GetRegistry().Remove(characterId)

		emitCashShopEntryRejection(l, span)(m.WorldId(), m.ChannelId(), m.CharacterId(), messageType, message)
		return nil
	}
}
