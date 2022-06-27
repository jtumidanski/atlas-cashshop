package cashshop

import (
	"atlas-cashshop/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerNameEnterCashShop        = "enter_cash_shop_command"
	consumerNamePurchaseCashShopItem = "purchase_cash_shop_item_command"
	topicTokenEnterCashShop          = "TOPIC_ENTER_CASH_SHOP_COMMAND"
	topicTokenPurchaseCashShopItem   = "TOPIC_PURCHASE_CASH_SHOP_ITEM_COMMAND"
)

func EnterCashShopCommandConsumer() func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[enterCashShopCommand](consumerNameEnterCashShop, topicTokenEnterCashShop, groupId, handleEnterCashShopCommand())
	}
}

type enterCashShopCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func handleEnterCashShopCommand() kafka.HandlerFunc[enterCashShopCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command enterCashShopCommand) {
		err := EnterCashShop(l, span)(command.WorldId, command.ChannelId, command.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to process enter cash shop command.")
		}
	}
}

func PurchaseCashShopItemCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[purchaseItemCommand](consumerNamePurchaseCashShopItem, topicTokenPurchaseCashShopItem, groupId, handlePurchaseCashShopItemCommand(db))
	}
}

type purchaseItemCommand struct {
	CharacterId  uint32 `json:"character_id"`
	CashIndex    uint32 `json:"cash_index"`
	SerialNumber uint32 `json:"serial_number"`
}

func handlePurchaseCashShopItemCommand(db *gorm.DB) kafka.HandlerFunc[purchaseItemCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command purchaseItemCommand) {
		err := PurchaseCashShopItem(l, db, span)(command.CharacterId, command.CashIndex, command.SerialNumber)
		if err != nil {
			l.WithError(err).Errorf("Unable to process item purchase command.")
		}
	}
}
