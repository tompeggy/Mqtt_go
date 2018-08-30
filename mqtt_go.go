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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var knt int

//
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("JSON: %s\n", msg.Payload())
	text := fmt.Sprintf("this is result msg #%d!", knt)
	knt++
	fmt.Printf("MSG: %s\n", text)
	//修改...
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	//push上到 broker..
	//token := client.Publish("GIOT-GW/UL/1C497BE1FD99", 0, false, text)
	//token.Wait()
}

func main() {
	knt = 0
	//建立loop...
	loopset := make(chan os.Signal, 1)
	signal.Notify(loopset, os.Interrupt, syscall.SIGTERM)
	//
	opts := MQTT.NewClientOptions().AddBroker("tcp://140.124.182.66:1883")
	opts.SetClientID("mac-go")
	opts.SetDefaultPublishHandler(f)
	topic := "GIOT-GW/UL/1C497BE1FD99"

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}
	//回到c 的loop...
	<-loopset
}
