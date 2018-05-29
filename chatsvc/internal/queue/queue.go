package queue

import (
	"github.com/streadway/amqp"
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"encoding/json"
)

type ChatMQ struct {
	MqConnection *amqp.Connection
}

type ReceiveFunc func(deliveries <-chan amqp.Delivery, disconnect <-chan int, stream chatpb.Chat_StartChatServer, channel *amqp.Channel)

func NewChatMQ(mqConnection *amqp.Connection) (chatMQ *ChatMQ) {

	return &ChatMQ{
		MqConnection: mqConnection,
	}
}

func (mq *ChatMQ) Send(chatConnectionID string, msg *chatpb.Message) (err error) {

	channel, err := mq.MqConnection.Channel()

	defer channel.Close()

	if err != nil {
		return err
	}

	msgJSON, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	err = channel.ExchangeDeclare(
		chatConnectionID,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = channel.Publish(
		chatConnectionID,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgJSON,
		})

	if err != nil {
		return err
	}

	return nil
}

func (mq *ChatMQ) ReceiveFromQueue(chatConnectionID string) (deliveries <-chan amqp.Delivery, channel *amqp.Channel, err error) {

	channel, err = mq.MqConnection.Channel()

	if err != nil {
		return deliveries, channel, err
	}

	err = channel.ExchangeDeclare(
		chatConnectionID,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return deliveries, channel, err
	}

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		return deliveries, channel, err
	}

	err = channel.QueueBind(
		queue.Name,
		"",
		chatConnectionID,
		false,
		nil,
	)

	if err != nil {
		return deliveries, channel, err
	}

	deliveries, err = channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return deliveries, channel, err
	}

	return deliveries, channel, nil
}
