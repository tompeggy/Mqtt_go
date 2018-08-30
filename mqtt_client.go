package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	//import the Paho Go MQTT library

	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

/*
var msgRcvd := func(client *mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
}
*/

func main() {
	//loop
	looptest := make(chan os.Signal, 1)
	signal.Notify(looptest, os.Interrupt, syscall.SIGTERM)
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://140.124.182.66:1883")
	opts.SetClientID("teco-mqtt-0001")
	opts.SetDefaultPublishHandler(f)
	topic := "GIOT-GW/UL/1C497BE1FD99"
	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		//	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
		token.Wait()
		//	os.Exit(1)
	}

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 10; i++ {
		text := fmt.Sprintf("this is msg #%d!\n", i)
		token := c.Publish(topic, 0, false, text)
		token.Wait()
	}

	time.Sleep(10 * time.Second)
	//unsubscribe from /go-mqtt/sample
	/*
		if token := c.Unsubscribe("GIOT-GW/UL/1C497BE1FD99"); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	*/
	<-looptest
	//	c.Disconnect(11250)
}
