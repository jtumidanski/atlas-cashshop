package cashshop

import (
	"atlas-cashshop/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerNameEnter             = "enter_cash_shop_command"
	consumerNamePurchaseItem      = "purchase_cash_shop_item_command"
	consumerNameGatekeeperCommand = "cash_shop_gatekeeper_command"
	consumerNameEntryPollResponse = "cash_shop_entry_poll_response"
	topicTokenEnter               = "TOPIC_ENTER_CASH_SHOP_COMMAND"
	topicTokenPurchaseItem        = "TOPIC_PURCHASE_CASH_SHOP_ITEM_COMMAND"
	topicTokenGatekeeperCommand   = "TOPIC_CASH_SHOP_GATEKEEPER_COMMAND"
	topicTokenEntryPollResponse   = "TOPIC_CASH_SHOP_ENTRY_POLL_RESPONSE"
)

func EnterCommandConsumer() func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[enterCashShopCommand](consumerNameEnter, topicTokenEnter, groupId, handleEnterCashShopCommand())
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

func PurchaseItemCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[purchaseItemCommand](consumerNamePurchaseItem, topicTokenPurchaseItem, groupId, handlePurchaseCashShopItemCommand(db))
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

func GatekeeperCommandConsumer() func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[gatekeeperCommand](consumerNameGatekeeperCommand, topicTokenGatekeeperCommand, groupId, handleGatekeeperCommand())
	}
}

type gatekeeperCommand struct {
	Service string `json:"service"`
	Type    string `json:"type"`
}

func handleGatekeeperCommand() kafka.HandlerFunc[gatekeeperCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command gatekeeperCommand) {
		if command.Type == "REGISTER" {
			err := RegisterGatekeeper(l, span)(command.Service)
			if err != nil {
				l.WithError(err).Errorf("Unable to register service %s interested in gatekeeping cash shop entry.", command.Type)
			}
		} else if command.Type == "UNREGISTER" {
			err := UnregisterGatekeeper(l, span)(command.Service)
			if err != nil {
				l.WithError(err).Errorf("Unable to unregister service %s interested in gatekeeping cash shop entry.", command.Type)
			}
		} else {
			l.Warnf("Unhandled command type %s.", command.Type)
		}
	}
}

func EntryPollResponseConsumer() func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[entryPollResponse](consumerNameEntryPollResponse, topicTokenEntryPollResponse, groupId, handleEntryPollResponse())
	}
}

type entryPollResponse struct {
	CharacterId uint32 `json:"character_id"`
	Service     string `json:"service"`
	Type        string `json:"type"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}

func handleEntryPollResponse() kafka.HandlerFunc[entryPollResponse] {
	return func(l logrus.FieldLogger, span opentracing.Span, response entryPollResponse) {
		if response.Type == "APPROVE" {
			err := GatekeeperApproval(l, span)(response.Service, response.CharacterId, response.MessageType, response.Message)
			if err != nil {
				l.WithError(err).Errorf("Unable to process entry poll response for character %d from service %s.", response.CharacterId, response.Service)
			}
		} else if response.Type == "DENY" {
			err := GatekeeperDenial(l, span)(response.Service, response.CharacterId, response.MessageType, response.Message)
			if err != nil {
				l.WithError(err).Errorf("Unable to process entry poll response for character %d from service %s.", response.CharacterId, response.Service)
			}
		} else {
			l.Warnf("Unhandled entry poll response %s from %s.", response.Type, response.Service)
		}
	}
}
