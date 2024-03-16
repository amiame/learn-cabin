package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press Enter to simulate casbin policy update")
		reader.ReadString('\n')
		produce()
		fmt.Print("A message has been sent to all MS telling them that casbin policy has been updated\n")
		time.Sleep(time.Second * 2)
		fmt.Print("\n")
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

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("please update your casbin enforcers!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
