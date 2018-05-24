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

func (mq *ChatMQ) Send(msg *chatpb.Message) (err error) {

	channel, err := mq.MqConnection.Channel()

	if err != nil {
		return err
	}

	msgJSON, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	err = channel.Publish(
		"",
		msg.ToContactId,
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

func (mq *ChatMQ) ReceiveFromQueue(ID string, receiveFunc ReceiveFunc, disconnect <-chan int, stream chatpb.Chat_StartChatServer) (err error) {

	channel, err := mq.MqConnection.Channel()

	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(
		ID,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	deliveries, err := channel.Consume(
		ID,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go receiveFunc(deliveries, disconnect, stream, channel)

	return nil
}
