package wishlist

import (
	"atlas-cashshop/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type wishListStatusEvent struct {
	CharacterId uint32 `json:"character_id"`
	Status      string `json:"status"`
}

func emitWishlistUpdateEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		emitWishlistStatusEvent(l, span)(characterId, "UPDATE")
	}
}

func emitWishlistStatusEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, status string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_WISHLIST_STATUS_EVENT")
	return func(characterId uint32, status string) {
		e := &wishListStatusEvent{
			CharacterId: characterId,
			Status:      status,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}
