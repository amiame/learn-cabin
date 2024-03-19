package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _enforcer *casbin.Enforcer
var _watcher *watcher

func main() {
	setupCasbin()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press Enter to simulate casbin policy update")
		if _, err := reader.ReadString('\n'); err != nil {
			log.Fatal("failed to read string:", err)
		}

		editPolicy()
	}
}

func setupCasbin() {
	e, err := casbin.NewEnforcer("./config/rbac_model.conf", "./config/initial_policy.csv")
	if err != nil {
		log.Fatal("failed to create enforcer:", err)
	}

	_enforcer = e
	_watcher = &watcher{}

	if err := _enforcer.SetWatcher(_watcher); err != nil {
		log.Fatal("failed to set watcher:", err)
	}
}

func editPolicy() {
	subject := time.Now().String()
	object := "xxx-xxxApi"
	action := "use"

	added, err := _enforcer.AddPolicy(subject, object, action)
	if err != nil {
		log.Fatal("failed to add authorization rule:", err)
	}

	if !added {
		log.Fatal("new authorization rule was not added")
	}

	if err := _enforcer.SavePolicy(); err != nil {
		log.Fatal("failed to save policy:", err)
	}
}

func produce() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal("failed to dial RabbitMQ:", err)
	}

	defer conn.Close()

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,          // ctx context.Context
		exchangeName, // exchange string
		"",           // routing_key string
		false,        // mandatory bool
		false,        // immediate bool
		amqp.Publishing{ // msg amqp.Publishing
			ContentType: "text/plain",
			Body:        []byte("please update your casbin enforcers!"),
		},
	)
	if err != nil {
		log.Fatal("failed to publish messages:", err)
	}

	fmt.Println("A message has been sent to Kafka to tell consumer MS them that they need to update their casbin policy")
	time.Sleep(time.Second * 2)
	fmt.Print("\n")
}
