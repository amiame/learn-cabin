package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/casbin/casbin/v3"
	kafka "github.com/segmentio/kafka-go"
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
	e, err := casbin.NewEnforcer("../config/rbac_model.conf", "../config/initial_policy.csv")
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
	// to produce messages
	topic := "messages"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Fatal("failed to set write deadline:", err)
	}
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("please update your casbin enforcers!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	fmt.Print("A message has been sent to all MS telling them that casbin policy has been updated\n")
	time.Sleep(time.Second * 2)
	fmt.Print("\n")
}
