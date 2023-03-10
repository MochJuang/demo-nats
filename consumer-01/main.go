package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
)

var (
	Sc stan.Conn
)

func main() {
	quit := make(chan struct{})
	ConnectStan("nats-consumer")

	//data := []byte("this is a test message")
	//PublishNats(data, "test")

	//PrintMessageSubscriber("test-01", "test-01")
	PrintMessageQueue("test-01", "test-01", "test-01")
	<-quit
}

func ConnectStan(clientID string) {
	clusterID := "test-cluster"     // nats cluster id
	url := "nats://127.0.0.1:10021" // nats url

	// you can set client id anything
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url),
		stan.Pings(1, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, url)
	}

	log.Println("Connected Nats")

	Sc = sc
}

func PublishNats(data []byte, channel string) {
	ach := func(s string, err2 error) {}
	_, err := Sc.PublishAsync(channel, data, ach)
	if err != nil {
		log.Fatalf("Error during async publish: %v\n", err)
	}
}

var num int

type Message struct {
	Message string
	Num     int
}

type History struct {
	Time string
	Num  int
}

func PrintMessageQueue(subject, qgroup, durable string) {
	var err error
	//f, err := os.OpenFile("./../data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//loc, _ := time.LoadLocation("Asia/Jakarta")

	mcb := func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ACK msg:%v", err)
		}
		fmt.Println("Receive : ", string(msg.Data), msg.Subject)

		//time.Sleep(2 * time.Second)
		//
		//
		//num += 1
		//jsonMsg := Message{}
		//err = json.Unmarshal(msg.Data, &jsonMsg)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//
		//if jsonMsg.Num == num {
		//	fmt.Println("Receive : ", string(msg.Data))
		//} else {
		//	data := History{
		//		Time: time.Now().In(loc).Format("15:04:05"),
		//		Num:  0,
		//	}
		//	byteArray, err := json.Marshal(data)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//	if _, err := fmt.Fprintf(f, "%s\n", byteArray); err != nil {
		//		fmt.Println(err)
		//	}
		//	n, err := io.WriteString(f, string(byteArray))
		//	if err != nil {
		//		fmt.Println(n, err)
		//	}
		//
		//}

	}

	_, err = Sc.QueueSubscribe(subject,
		qgroup, mcb,
		stan.SetManualAckMode(),
		stan.StartWithLastReceived(),
		stan.DeliverAllAvailable(),
		//stan.AckWait(1*time.Second),
		stan.MaxInflight(3),
		//stan.DurableName(durable),
	)
	if err != nil {
		log.Println(err)
	}

}

func PrintMessageSubscriber(subject, durable string) {
	var err error

	mcb := func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ACK msg:%v", err)
		}
		fmt.Println("Receive : ", string(msg.Data))
	}

	_, err = Sc.
		Subscribe(subject,
			mcb,
			stan.SetManualAckMode(),
			stan.StartWithLastReceived(),
			//stan.DeliverAllAvailable(),
			//stan.DurableName(durable),
		)
	if err != nil {
		log.Println(err)
	}

}
