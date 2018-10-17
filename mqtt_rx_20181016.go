/***********************************************************************************
PROJECT NAME  Mqtt Client...
DESCRIPTION   Mqtt Client封包處埋.....
HardWare      ***
Copyright   : 2018 TECO Ltd.
All Rights Reserved
//influxdb使用...
influx
> create user teco with password 'teco1133'
> grant all privileges to teco  //授權全部的資料庫(等於管理者帳號)...
> CREATE DATABASE chiler
> USE chiler
> SELECT * FROM "MAC_013157800087586"


***********************************************************************************/
/***********************************************************************************
Revision History
DD.MM.YYYY OSO-UID Description..
30.08.2018 Thomas start work..
25.09.2018 Thomas 修改成NB-IoT方案,import github用指令 go get 加入libery...
16.10.2018 Thomas 修改收NB-IoT,打到influx...
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
	MyDB     = "chiler"
	username = "teco"
	password = "teco1133"

//	MyMeasurement = "cpu_usage"
)

//設定MQTT接收的Buffer...
const (
	MaxClientIdLen = 10
)

var Knt int
var awsKnt int
var Rxdata [100]uint8
var testData float32
var data [13]float32

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
	//if strings.Contains(s, "0000000005010e9e") {
	if strings.Contains(s, "013157800087586") {
		indexMessage = strings.Index(s, "data")
		fmt.Printf("index %d\n", indexMessage)
		//for i := 0; i < 58; i++ {
		for i := 0; i < 68; i++ {
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

var awsfun MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	//	mQTTMessage := &MQTTMessage{msg, m}
	//	fmt.Printf("TOPIC: %s\n", msg.Topic())
	//var indexAwsMessage int
	awsKnt++
	if awsKnt > 65530 {
		awsKnt = 0
	}
	fmt.Printf("awsMSG: %d\n", awsKnt)
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
/*
func connInflux() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("influx pass\n")
	return cli

}

*/
func main() {
	Knt = 0

	//influxdb測試...
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Printf("pass http\n")

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pass creat db\n")

	data[0] = 27.8
	data[1] = 26.8
	data[2] = 14.6
	data[3] = 16.7
	data[4] = 33.8
	data[5] = 54.6
	data[6] = 4.15
	data[7] = 4.18
	data[8] = 7.75
	data[9] = 2.8
	data[10] = 380.0
	data[11] = 0.0
	data[12] = 80.8

	// Create a point and add to batch
	tags := map[string]string{"location": "bank"}
	fields := map[string]interface{}{
		"water_temp_in":      data[0],
		"water_temp_out":     data[1],
		"ice_temp_in":        data[2],
		"ice_temp_out":       data[3],
		"left_temp_out":      data[4],
		"right_temp_out":     data[5],
		"left_pressor_high":  data[6],
		"left_pressor_low":   data[7],
		"right_pressor_high": data[8],
		"right_pressor_low":  data[9],
		"voltage":            data[10],
		"left_amp":           data[11],
		"right_amp":          data[12],
	}

	pt, err := client.NewPoint("MAC_013157800087586", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pass write\n")
	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}

	//建立getRandomClientId returns randomized ClientId.
	clientId := getRandomClientId()
	fmt.Printf("clientId: %s\n", clientId)
	//建立Go的監聴封包...
	sigc := make(chan os.Signal, 1)

	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://140.124.182.66:1883")
	//opts := MQTT.NewClientOptions().AddBroker("tcp://13.114.3.126:1883")

	//opts.SetClientID("u89-0001")
	opts.SetClientID(clientId)

	opts.SetDefaultPublishHandler(f)
	//

	//topic := "GIOT-GW/UL/1C497BE1FD99"
	//topic := "iaq"
	topic := "NB-IoT"
	//create and start a client using the above ClientOptions
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to 140.124.182.66 server\n")
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

	//建立第二個MQTT Cliend...
	awsTopic := "iaq-test"
	awsMQTTBroker := MQTT.NewClientOptions().AddBroker("tcp://13.114.3.126:1883")
	awsMQTTBroker.SetClientID(clientId)
	awsMQTTBroker.SetDefaultPublishHandler(awsfun)

	//create and start a client using the above ClientOptions
	awsClient := MQTT.NewClient(awsMQTTBroker)
	if awsToken := awsClient.Connect(); awsToken.Wait() && awsToken.Error() != nil {
		panic(awsToken.Error())
	} else {
		fmt.Printf("Connected to 13.114.3.126 server\n")
	}

	if awsToken := awsClient.Subscribe(awsTopic, 0, nil); awsToken.Wait() && awsToken.Error() != nil {
		fmt.Println(awsToken.Error())
		awsToken.Wait()
		//	os.Exit(1)
	}
	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 1; i++ {
		text := fmt.Sprintf("this is msg #%d!\n", i)
		awsToken := awsClient.Publish(awsTopic, 0, false, text)
		awsToken.Wait()
	}

	time.Sleep(5 * time.Second)

	// Wait for receiving a signal.
	//<-sigc
	s := <-sigc
	fmt.Println("Got signal:", s) //Got signal: terminated
	//	c.Disconnect(11250)
}
