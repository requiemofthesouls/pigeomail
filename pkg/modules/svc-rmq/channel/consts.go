package channel

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeDirect = amqp.ExchangeDirect
	ExchangeTopic  = amqp.ExchangeTopic
	ExchangeFanout = amqp.ExchangeFanout

	DeliveryModeTransient  = amqp.Transient
	DeliveryModePersistent = amqp.Persistent
)
