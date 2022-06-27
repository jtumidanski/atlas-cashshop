package character

import (
	"atlas-cashshop/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerNameCharacterCreated = "character_created_event"
	consumerNameAwardCredit      = "award_credit_command"
	consumerNameAwardPoints      = "award_points_command"
	consumerNameAwardPrepaid     = "award_prepaid_command"
	topicTokenCharacterCreated   = "TOPIC_CHARACTER_CREATED_EVENT"
	topicTokenAwardCredit        = "TOPIC_AWARD_CREDIT_COMMAND"
	topicTokenAwardPoints        = "TOPIC_AWARD_POINTS_COMMAND"
	topicTokenAwardPrepaid       = "TOPIC_AWARD_PREPAID_COMMAND"
)

func CreatedConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[createdEvent](consumerNameCharacterCreated, topicTokenCharacterCreated, groupId, handleCreated(db))
	}
}

type createdEvent struct {
	CharacterId uint32 `json:"characterId"`
	WorldId     byte   `json:"worldId"`
	Name        string `json:"name"`
}

func handleCreated(db *gorm.DB) kafka.HandlerFunc[createdEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event createdEvent) {
		err := Created(l, db, span)(event.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to process enter cash shop command.")
		}
	}
}

func AwardCreditConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[awardCreditCommand](consumerNameAwardCredit, topicTokenAwardCredit, groupId, handleAwardCredit(db))
	}
}

func handleAwardCredit(db *gorm.DB) kafka.HandlerFunc[awardCreditCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command awardCreditCommand) {
		err := AwardCredit(l, db, span)(command.CharacterId, command.Amount)
		if err != nil {
			l.WithError(err).Errorf("Error processing award credit command.")
		}
	}
}

type awardCreditCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint32 `json:"amount"`
}

func AwardPointsConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[awardPointsCommand](consumerNameAwardPoints, topicTokenAwardPoints, groupId, handleAwardPoints(db))
	}
}

func handleAwardPoints(db *gorm.DB) kafka.HandlerFunc[awardPointsCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command awardPointsCommand) {
		err := AwardPoints(l, db, span)(command.CharacterId, command.Amount)
		if err != nil {
			l.WithError(err).Errorf("Error processing award credit command.")
		}
	}
}

type awardPointsCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint32 `json:"amount"`
}

func AwardPrepaidConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[awardPrepaidCommand](consumerNameAwardPrepaid, topicTokenAwardPrepaid, groupId, handleAwardPrepaid(db))
	}
}

func handleAwardPrepaid(db *gorm.DB) kafka.HandlerFunc[awardPrepaidCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command awardPrepaidCommand) {
		err := AwardPrepaid(l, db, span)(command.CharacterId, command.Amount)
		if err != nil {
			l.WithError(err).Errorf("Error processing award credit command.")
		}
	}
}

type awardPrepaidCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint32 `json:"amount"`
}
