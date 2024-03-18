package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/casbin/casbin/v2"
	kafka "github.com/segmentio/kafka-go"
)

var _enforcer *casbin.Enforcer

func main() {
	setupCasbin()

	// make a new reader that consumes from topic-A
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Partition: 0,
		Topic:     "messages",
		MaxBytes:  10e6, // 10MB
	})
	// Set last offset so we don't get old messages
	if err := r.SetOffset(kafka.LastOffset); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	go func() {
		// Set up channel on which to send signal notifications.
		// We must use a buffered channel or risk missing the signal
		// if we're not ready to receive when the signal is sent.
		trap := make(chan os.Signal, 1)
		signal.Notify(trap, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

		// Block until a signal is received
		s := <-trap
		fmt.Println("Got a signal:", s)
		fmt.Println("Closing kafka Reader...")

		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		//fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		fmt.Printf("Received message: %s\n", string(m.Value))
		renewEnforcer()
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func setupCasbin() {
	e, err := casbin.NewEnforcer("./config/rbac_model.conf", "./config/initial_policy.csv")
	if err != nil {
		log.Fatal("failed to create enforcer:", err)
	}

	_enforcer = e
}

func renewEnforcer() {
	_enforcer.LoadPolicy()
	fmt.Println("renewed enforcer with latest policy")
}
