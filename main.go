package main

import (
	//"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"time"
	//"io/ioutil"
)

var LIGHT_TOPIC string = "esp8266_arduino_out_esp8266_again"
var mqttClient mqtt.Client

var alexaResponse = "{\"version\": \"0.1\"," +
	"\"sessionAttributes\": {},\"response\": {\"outputSpeech\": {\"type\": \"PlainText\",\"text\":\"Okay\"},\"card\": {\"type\": \"Simple\",\"title\": \"Oberon\",\"content\": \"Oberon will fulfill your command.\"}," +
	"\"reprompt\": {\"outputSpeech\": {\"type\": \"PlainText\",\"text\": \"No!\"}},\"shouldEndSession\": true}}"

type AlexaMessage struct {
	Version string
	Session AlexaSession
	Request AlexaRequest
}

type AlexaSession struct {
	SessionId   string
	Application AlexaApplication
	User        AlexaUser
	New         bool
}

type AlexaRequest struct {
	Type      string
	RequestId string
	Locale    string
	TimeStamp string
	Intent    AlexaIntent
}

type AlexaApplication struct {
	ApplicationId string
}

type AlexaUser struct {
	UserId string
}

type AlexaIntent struct {
	Name  string
	Slots AlexaSlot
}

type AlexaSlot struct {
	State AlexaState
}

type AlexaState struct {
	Name  string
	Value string
}

func handlePing(c *gin.Context) {
	c.String(200, "pong")
}

func handlePost(c *gin.Context) {
	//x, _ := ioutil.ReadAll(c.Request.Body)
	var msg AlexaMessage
	//json.Unmarshal(x, &msg)

	//// var msg AlexaMessage
	c.BindJSON(&msg)
	if msg.Request.Type == "IntentRequest" {
		if msg.Request.Intent.Name == "Light" {
			handleLightIntent(msg, c)
		} else if msg.Request.Intent.Name == "CheckWindows" {
			handleCheckWindowsIntent(msg, c)
		} else if msg.Request.Intent.Name == "CheckTemperature" {
			handleCheckTemperatureIntent(msg, c)
		}
	}
	// fmt.Println(alexaResponse)

	c.String(200, alexaResponse)
}

func handleLightIntent(msg AlexaMessage, c *gin.Context) {
	if msg.Request.Intent.Slots.State.Name == "state" {
		if msg.Request.Intent.Slots.State.Value == "up" || msg.Request.Intent.Slots.State.Value == "on" {
			fmt.Println("Turn on!")
			mqttClient.Publish(LIGHT_TOPIC, 1, false, "Sometext")
		} else if msg.Request.Intent.Slots.State.Value == "down" || msg.Request.Intent.Slots.State.Value == "off" {
			fmt.Println("Turn off!")
			mqttClient.Publish(LIGHT_TOPIC, 1, false, "Sometext")
		}
	}
}

func handleCheckWindowsIntent(msg AlexaMessage, c *gin.Context) {
	fmt.Println("Checking the windows...")
}

func handleCheckTemperatureIntent(msg AlexaMessage, c *gin.Context) {
	fmt.Println("Checking the temperature...")
}

func main() {
	clientId := "brian-agent-" + fmt.Sprintf("%d", time.Now().Nanosecond())
	fmt.Println("Using client id: " + clientId)
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://iot.eclipse.org:1883")
	opts.SetClientID(clientId)
	mqttClient = mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Client connected!")

	r := gin.Default()
	r.GET("alexa/ping", handlePing)
	r.POST("alexa", handlePost)
	r.Run(":8080")
}
