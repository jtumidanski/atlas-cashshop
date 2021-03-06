package cashshop

import (
	"atlas-cashshop/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type cashShopEnterEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func emitEnterCashShop(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ENTER_CASH_SHOP_EVENT")
	return func(worldId byte, channelId byte, characterId uint32) {
		e := &cashShopEnterEvent{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type cashShopEntryPoll struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func pollCashShopEntry(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CASH_SHOP_ENTRY_POLL")
	return func(worldId byte, channelId byte, characterId uint32) {
		e := &cashShopEntryPoll{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type cashShopEntryRejectionEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}

func emitCashShopEntryRejection(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, messageType string, message string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CASH_SHOP_ENTRY_REJECTION_EVENT")
	return func(worldId byte, channelId byte, characterId uint32, messageType string, message string) {
		e := &cashShopEntryRejectionEvent{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
			MessageType: messageType,
			Message:     message,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}
