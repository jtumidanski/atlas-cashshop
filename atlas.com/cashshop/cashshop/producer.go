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
