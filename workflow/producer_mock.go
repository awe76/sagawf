package workflow

import (
	"encoding/json"
	"fmt"
)

type ProducerMock struct {
	messages map[string][]string
}

func NewProducerMock() *ProducerMock {
	return &ProducerMock{
		messages: make(map[string][]string),
	}
}

func (p *ProducerMock) Init() error {
	return nil
}

func (p *ProducerMock) Connect() error {
	return nil
}

func (p *ProducerMock) SendMessage(topic string, message interface{}) error {
	raw, err := json.Marshal(message)
	if err != nil {
		return err
	}

	messages, found := p.messages[topic]
	if found {
		p.messages[topic] = append(messages, string(raw))
	} else {
		p.messages[topic] = append([]string{}, string(raw))
	}

	return nil
}

func (p *ProducerMock) Has(topic string, message interface{}) bool {
	raw, err := json.Marshal(message)
	if err != nil {
		return false
	}

	pattern := string(raw)
	if messages, found := p.messages[topic]; found {
		for _, message := range messages {
			if message == pattern {
				return true
			}
		}
	}

	return false
}

func (p *ProducerMock) Print() {
	for topic, messages := range p.messages {
		fmt.Printf("%s\n", topic)
		for _, message := range messages {
			fmt.Printf("%v\n", message)
		}
	}
}
