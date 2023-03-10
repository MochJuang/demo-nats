package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

var (
	Sc stan.Conn
)

func main() {
	//quit := make(chan struct{})
	ConnectStan("nats-publisher")

	num := 1
	for {
		data := []byte(fmt.Sprintf("{\"Message\" : \"This is message in %v\", \"Num\" : %v }", num, num))
		fmt.Println("Send : ", string(data))
		PublishNats(data, "test-01")
		num++
		time.Sleep(time.Second * 1)
	}

	//PrintMessage("test", "test", "test-1")
	//<-quit
}

func ConnectStan(clientID string) {
	clusterID := "test-cluster"     // nats cluster id
	url := "nats://127.0.0.1:10021" // nats url

	// you can set client id anything
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url),
		stan.Pings(1, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}),
		//stan.MaxPubAcksInflight(1000),
	)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, url)
	}

	log.Println("Connected Nats")

	Sc = sc
}

func PublishNats(data []byte, channel string) {
	ach := func(s string, err2 error) {
		fmt.Println("")
	}
	_, err := Sc.PublishAsync(channel, data, ach)
	if err != nil {
		log.Fatalf("Error during async publish: %v\n", err)
	}
}

func PrintMessage(subject, qgroup, durable string) {
	mcb := func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ACK msg:%v", err)
		}

		fmt.Println(string(msg.Data), msg.Subject)

	}

	_, err := Sc.QueueSubscribe(subject,
		qgroup, mcb,
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.DurableName(durable))
	if err != nil {
		log.Println(err)
	}
}

func PrintMessageSubscriber(subject, qgroup, durable string) {
	mcb := func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ACK msg:%v", err)
		}

		fmt.Println(string(msg.Data))
	}

	_, err := Sc.QueueSubscribe(subject,
		qgroup, mcb,
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.DurableName(durable))
	if err != nil {
		log.Println(err)
	}
}
