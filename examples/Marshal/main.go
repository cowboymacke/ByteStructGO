package main

import (
	"encoding/binary"
	"fmt"

	structGo "github.com/MarcusLing/ByteStructGO"
)

type Message struct {
	Length  uint16 `bytePayload:"Payload"`
	Payload string
}

type MessageWithStruct struct {
	Length  uint16 `bytePayload:"Payload"`
	Payload Message
}

type MessageWithArrayStruct struct {
	Length  uint16 `bytePayload:"Payload"`
	Payload []Message
}

type MixedStruct struct {
	MessageType    uint8
	UserNameLength uint8 `bytePayload:"UserName"`
	UserName       string
	Data           UserData
}

type UserData struct {
	Age           uint8
	CountryLength uint8 `bytePayload:"Country"`
	Country       string
}

func main() {

	message := Message{
		Payload: "Hello World",
	}

	data, err := structGo.Marshal(binary.BigEndian, message)
	if err != nil {
		fmt.Printf("Failed to Marshal: %s", err.Error())
		return
	}

	fmt.Printf("byteArray: %d", data) // -> byteArray: [0 11 72 101 108 108 111 32 87 111 114 108 100]

	messageWithStruct := MessageWithStruct{
		Payload: message,
	}

	data, err = structGo.Marshal(binary.BigEndian, messageWithStruct)
	if err != nil {
		fmt.Printf("Failed to Marshal: %s", err.Error())
		return
	}

	fmt.Printf("byteArray: %d", data) // -> byteArray: [0 13 0 11 72 101 108 108 111 32 87 111 114 108 100]

	messageWithArrayStruct := MessageWithArrayStruct{
		Payload: []Message{{Payload: "Hello"}, {Payload: "World"}},
	}

	data, err = structGo.Marshal(binary.BigEndian, messageWithArrayStruct)
	if err != nil {
		fmt.Printf("Failed to Marshal: %s", err.Error())
		return
	}

	fmt.Printf("byteArray: %d", data) // -> byteArray: [0 14 0 5 72 101 108 108 111 0 5 87 111 114 108 100]

	mixed := MixedStruct{
		MessageType: 2,
		UserName:    "Foo",
		Data: UserData{
			Age:     30,
			Country: "Sweden",
		},
	}

	data, err = structGo.Marshal(binary.BigEndian, mixed)
	if err != nil {
		fmt.Printf("Failed to Marshal: %s", err.Error())
		return
	}

	fmt.Printf("byteArray: %d", data) // -> byteArray: [2 3 70 111 111 30 6 83 119 101 100 101 110]
}
