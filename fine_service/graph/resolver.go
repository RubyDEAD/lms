package graph

import (
	"database/sql"
	"fine_service/graph/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Resolver struct {
	DB     *sql.DB
	Rabbit *amqp.Channel
	RabbitChannel         *amqp.Channel
	FineCreatedSubscribers []chan *model.Fine
	ViolationSubscribers    []chan *model.ViolationRecord
}


