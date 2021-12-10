package workflow

import (
	"encoding/json"

	"go-micro.dev/v4/broker"
)

type producer struct {
}

func NewProducer() Producer {
	return &producer{}
}

type Producer interface {
	Init() error
	Connect() error
	SendMessage(topic string, message interface{}) error
}

func (p *producer) Init() error {
	return broker.Init()
}

func (p *producer) Connect() error {
	return broker.Connect()
}

func (p *producer) SendMessage(topic string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &broker.Message{
		Body: body,
	}

	return broker.Publish(topic, msg)
}
