/***********************************************************************************
PROJECT NAME  Mqtt Client...
DESCRIPTION   Mqtt Client封包處埋.....
HardWare      ***
Copyright   : 2018 TECO Ltd.
All Rights Reserved
***********************************************************************************/
/***********************************************************************************
Revision History
DD.MM.YYYY OSO-UID Description..
30.08.2018 Thomas start work..
***********************************************************************************/
package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/influxdata/influxdb/client/v2"

	"time"
	//import the Paho Go MQTT library
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//
const (
	MyDB          = "MAC_5010e9e"
	username      = "thomas"
	password      = "teco1133"
	MyMeasurement = "cpu_usage"
)

//設定MQTT接收的Buffer...
const (
	MaxClientIdLen = 10
)

var Knt int
var Rxdata [100]uint8

// var outgoing chan *MQTTMessage

//define a function for the default message handler
//set callback function
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	//	mQTTMessage := &MQTTMessage{msg, m}
	//	fmt.Printf("TOPIC: %s\n", msg.Topic())
	var indexMessage int
	Knt++
	if Knt > 65530 {
		Knt = 0
	}
	fmt.Printf("MSG: %d\n", Knt)
	//fmt.Printf("%s\n", msg.Payload())

	//轉換資料...
	s := string(msg.Payload()[:])
	//	fmt.Printf("%s\n", s)
	//fmt.Println(strings.Contains(s, "data"))
	//比較字串,並接收資料...
	if strings.Contains(s, "0000000005010e9e") {
		indexMessage = strings.Index(s, "data")
		fmt.Printf("index %d\n", indexMessage)
		for i := 0; i < 58; i++ {
			Rxdata[i] = s[indexMessage+i+7]
		}
		fmt.Printf("Rxdata: %s\n", Rxdata)
		//寫入influxdb...

	}

	/*
		//	fmt.Printf("check: %t\n", (s == "data"))
		if s == "on" {
			fmt.Println("on is received!")
			//	TurnAllOn()
		} else if s == "off" {
			fmt.Println("off is received!")
			//	GlowOff()
		}
	*/
	//	fmt.Println(strings.Contains(msg.Payload(), "data")) //true
}

/*
var msgRcvd := func(client *mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
}
*/

// getRandomClientId returns randomized ClientId.
func getRandomClientId() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, MaxClientIdLen)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return "ClientID-" + string(bytes)
}

//influx tets
func connInflux() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli

}

func main() {
	Knt = 0

	//influxdb測試...
	conn := connInflux()
	fmt.Println(conn)

	//建立getRandomClientId returns randomized ClientId.
	clientId := getRandomClientId()
	fmt.Printf("clientId: %s\n", clientId)
	//建立Go的監聴封包...
	sigc := make(chan os.Signal, 1)

	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	//opts := MQTT.NewClientOptions().AddBroker("tcp://140.124.182.66:1883")
	opts := MQTT.NewClientOptions().AddBroker("tcp://13.114.3.126:1883")

	//opts.SetClientID("u89-0001")
	opts.SetClientID(clientId)

	opts.SetDefaultPublishHandler(f)
	//

	//topic := "GIOT-GW/UL/1C497BE1FD99"
	topic := "iaq"

	//create and start a client using the above ClientOptions
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}

	//subscribe to the topic "iaq" and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		//	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
		token.Wait()
		//	os.Exit(1)
	}

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 1; i++ {
		text := fmt.Sprintf("this is msg #%d!\n", i)
		token := client.Publish(topic, 0, false, text)
		token.Wait()
	}

	time.Sleep(5 * time.Second)
	//unsubscribe from /go-mqtt/sample
	/*
		if token := c.Unsubscribe("GIOT-GW/UL/1C497BE1FD99"); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	*/

	// Wait for receiving a signal.
	//<-sigc
	s := <-sigc
	fmt.Println("Got signal:", s) //Got signal: terminated
	//	c.Disconnect(11250)
}
