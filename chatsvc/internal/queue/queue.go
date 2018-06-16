package queue

import (
	"github.com/streadway/amqp"
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type ChatMQ struct {
	MqConnections []*amqp.Connection
	roundRobinIdx int
	roundRobinLock *sync.Mutex
}

type ReceiveFunc func(deliveries <-chan amqp.Delivery, disconnect <-chan int, stream chatpb.Chat_StartChatServer, channel *amqp.Channel)

func NewChatMQ(urls []string) (chatMQ *ChatMQ, err error) {

	var mqConnections []*amqp.Connection

	for _, url := range urls {

		mqConnection, err := amqp.Dial(url)

		if err != nil {

			log.Printf("Fail to dial RabbitMQ %s, %s", url, err)

		} else {

			log.Printf("Dial RabbitMQ %s Successful.", url)
			mqConnections = append(mqConnections, mqConnection)

		}
	}

	if len(mqConnections) == 0 {
		return nil, fmt.Errorf("failed to dial to rabbitmq cluster for all URLS")
	}

	return &ChatMQ{
		MqConnections: mqConnections,
		roundRobinIdx: 0,
		roundRobinLock: &sync.Mutex{},
	}, nil
}

func (mq *ChatMQ) GetConnection() *amqp.Connection {

	mq.roundRobinLock.Lock()

	log.Printf("Using RabbitMQ conn #%d.", mq.roundRobinIdx)
	mqConn := mq.MqConnections[mq.roundRobinIdx]

	if mq.roundRobinIdx == (len(mq.MqConnections) - 1) {
		mq.roundRobinIdx = 0
	} else {
		mq.roundRobinIdx++
	}

	mq.roundRobinLock.Unlock()

	return mqConn
}

func (mq *ChatMQ) Send(chatConnectionID string, msg *chatpb.Message) (err error) {

	channel, err := mq.GetConnection().Channel()

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

	channel, err = mq.GetConnection().Channel()

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
