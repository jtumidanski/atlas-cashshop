package wishlist

import (
	"atlas-cashshop/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerNameModifyWishlist = "modify_wishlist_command"
	topicTokenModifyWishlist   = "TOPIC_MODIFY_WISHLIST_COMMAND"
)

func ModifyWishlistConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[modifyWishlistCommand](consumerNameModifyWishlist, topicTokenModifyWishlist, groupId, handleModifyWishlist(db))
	}
}

func handleModifyWishlist(db *gorm.DB) kafka.HandlerFunc[modifyWishlistCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command modifyWishlistCommand) {
		err := ModifyWishlist(l, db, span)(command.CharacterId, command.SerialNumbers)
		if err != nil {
			l.WithError(err).Errorf("Error processing award credit command.")
		}
	}
}

type modifyWishlistCommand struct {
	CharacterId   uint32   `json:"character_id"`
	SerialNumbers []uint32 `json:"serial_numbers"`
}
