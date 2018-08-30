package main

import (
	"fmt"
	// via: https://eclipse.org/paho/clients/golang/
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	//MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

// Sangoを使ってるならダッシュボード、自分で立てたならそのIPです
const MQTT_BROKER = "140.124.182.66:1883"

func createMQTTClient(brokerAddr, clientId, username, password string) *MQTT.Client {
	opts := MQTT.NewClientOptions().AddBroker(brokerAddr)
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	client := MQTT.NewClient(opts)
	return client
}

func subscribe(client *MQTT.Client, sub chan<- MQTT.Message) {
	fmt.Println("start subscribing...")
	// forループしなくても勝手に中でループしててくれます。なんかそういうのってあんまりGolangっぽくない気がしますけど。
	subToken := client.Subscribe(
		"path/to/topic", // Sangoのダッシュボードからコピペしましょう
		0,
		func(client *MQTT.Client, msg MQTT.Message) {
			sub <- msg
		})
	if subToken.Wait() && subToken.Error() != nil {
		fmt.Println(subToken.Error())
		os.Exit(1)
	}
}

func publish(client *MQTT.Client, input string) {
	token := client.Publish("path/to/topic/chat", 0, true, input)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func input(pub chan<- string) {
	for {
		var input string
		fmt.Scanln(&input)
		pub <- input
	}
}

func main() {
	fmt.Print("your id: ")
	var id string
	fmt.Scanln(&id)
	client := createMQTTClient(MQTT_BROKER, id, "YourID", "YouPassword")
	defer client.Disconnect(250)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	sub := make(chan MQTT.Message)
	go subscribe(client, sub)
	pub := make(chan string)
	go input(pub)
	for {
		select {
		case s := <-sub:
			msg := string(s.Payload())
			fmt.Printf("\nmsg: %s\n", msg)
		case p := <-pub:
			publish(client, p)
		}
	}
}
