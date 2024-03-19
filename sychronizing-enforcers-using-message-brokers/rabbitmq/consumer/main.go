package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _enforcer *casbin.Enforcer

func main() {
	setupCasbin()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal("failed to dial rabbitmq:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("failed to create channel:", err)
	}

	exchangeName := "update_reminder"
	err = ch.ExchangeDeclare(
		exchangeName, // name string
		"fanout",     // kind string
		true,         // durable bool
		false,        // autoDelete bool
		false,        // internal bool
		false,        // noWait bool
		nil,          // args amqp.Table
	)
	if err != nil {
		log.Fatal("failed to declare an exchange:", err)
	}

	q, err := ch.QueueDeclare(
		"",    // name string
		false, // durable bool
		false, // autoDelete bool
		true,  // exclusive bool
		false, // noWait bool
		nil,   // args amqp.Table
	)
	if err != nil {
		log.Fatal("failed to declare a queue:", err)
	}

	err = ch.QueueBind(
		q.Name,       // name string
		"",           // key string
		exchangeName, // exchange string
		false,        // noWait bool
		nil,          // args amqp.Table
	)
	if err != nil {
		log.Fatal("failed to bind to an queue:", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue string
		"",     // consumer string
		true,   // autoAck bool
		false,  // exclusive bool
		false,  // noLocal bool
		false,  // noWait bool
		nil,    // args amqp.Table
	)
	if err != nil {
		log.Fatal("failed to register a consumer:", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			fmt.Printf("Received message: %s\n", string(d.Body))
			renewEnforcer()
			writeToFile()
		}
	}()

	fmt.Println("Connected to RabbitMQ. Waiting for messages")
	<-forever
}

func setupCasbin() {
	e, err := casbin.NewEnforcer("./config/rbac_model.conf", "./config/initial_policy.csv")
	if err != nil {
		log.Fatal("failed to create enforcer:", err)
	}

	_enforcer = e
}

func renewEnforcer() {
	if err := _enforcer.LoadPolicy(); err != nil {
		log.Fatal("failed to load policy:", err)
	}

	fmt.Println("renewed enforcer with latest policy")
}

func writeToFile() {
	f, err := os.Create("./consumer/consumer_policy.csv")
	if err != nil {
		log.Fatal("error creating file:", err)
	}

	defer f.Close()

	policyTypes := map[string]func() [][]string{
		"p": _enforcer.GetPolicy,
		"g": _enforcer.GetGroupingPolicy,
	}

	for policyType, policyFunction := range policyTypes {
		rules := policyFunction()
		for _, rule := range rules {
			str := fmt.Sprintf("%s, %s\n", policyType, strings.Join(rule, ", "))
			_, err := f.WriteString(str)
			if err != nil {
				log.Fatal("error writing string:", err)
			}
		}
	}

	if err := f.Sync(); err != nil {
		log.Fatal("failed to sync file:", err)
	}
}
