package plugins

import (
    "encoding/json"
    "fmt"
    "github.com/streadway/amqp"
    "log"
)

type MQExecutor struct {
    connection *amqp.Connection
    channel    *amqp.Channel
}

func NewMQExecutor(amqpURL string) (*MQExecutor, error) {
    conn, err := amqp.Dial(amqpURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open a channel: %w", err)
    }

    return &MQExecutor{
        connection: conn,
        channel:    ch,
    }, nil
}

func (m *MQExecutor) Send(queueName string, message interface{}) error {
    body, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %w", err)
    }

    _, err = m.channel.QueueDeclare(
        queueName,
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %w", err)
    }

    err = m.channel.Publish(
        "",
        queueName,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        return fmt.Errorf("failed to publish message: %w", err)
    }

    return nil
}

func (m *MQExecutor) Close() {
    if err := m.channel.Close(); err != nil {
        log.Printf("failed to close channel: %v", err)
    }
    if err := m.connection.Close(); err != nil {
        log.Printf("failed to close connection: %v", err)
    }
}