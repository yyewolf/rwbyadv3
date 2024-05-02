package rbmq

import "github.com/rabbitmq/amqp091-go"

func (j *JobHandler) Init() error {
	ch, err := j.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		j.config.Rbmq.Jobs.Exchange,
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		j.config.Rbmq.Jobs.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		j.config.Rbmq.Jobs.Queue,
		"job",
		j.config.Rbmq.Jobs.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
