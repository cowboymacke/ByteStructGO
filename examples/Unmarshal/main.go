package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	structGo "github.com/MarcusLing/ByteStructGO"
)

type Message struct {
	Length  uint16
	Payload string `byteSize:"Length"`
}

type MessageWithStruct struct {
	Length  uint16
	Payload Message `byteSize:"Length"`
}

type MessageWithArrayStruct struct {
	Length  uint16
	Payload []Message `byteSize:"Length"`
}

type MixedStruct struct {
	MessageType    uint8
	UserNameLength uint8
	Username       string `byteSize:"UserNameLength"`
	Data           UserData
}

type UserData struct {
	Age           uint8
	CountryLength uint8
	Country       string `byteSize:"CountryLength"`
}

func main() {

	reader := bytes.NewBuffer([]byte{0, 11, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100})

	var message Message
	err := structGo.Unmarshal(reader, binary.BigEndian, &message)
	if err != nil {
		fmt.Printf("Failed to Unmarshal: %s", err.Error())
		return
	}

	fmt.Printf("Payload: %s", message.Payload) // -> Payload: Hello World

	reader = bytes.NewBuffer([]byte{0, 13, 0, 11, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100})
	var messageWithStruct MessageWithStruct
	err = structGo.Unmarshal(reader, binary.BigEndian, &messageWithStruct)
	if err != nil {
		fmt.Printf("Failed to Unmarshal: %s", err.Error())
		return
	}

	fmt.Printf("Payload: %s", messageWithStruct.Payload.Payload) // -> Payload: Hello World

	reader = bytes.NewBuffer([]byte{0, 14, 0, 5, 72, 101, 108, 108, 111, 0, 5, 87, 111, 114, 108, 100})
	var messageWithArrayStruct MessageWithArrayStruct
	err = structGo.Unmarshal(reader, binary.BigEndian, &messageWithArrayStruct)
	if err != nil {
		fmt.Printf("Failed to Unmarshal: %s", err.Error())
		return
	}

	fmt.Printf("Payload[0].Payload: %s", messageWithArrayStruct.Payload[0].Payload) // -> Payload[0]: Hello

	reader = bytes.NewBuffer([]byte{2, 3, 70, 111, 111, 30, 6, 83, 119, 101, 100, 101, 110})
	var mixedStruct MixedStruct
	err = structGo.Unmarshal(reader, binary.BigEndian, &mixedStruct)
	if err != nil {
		fmt.Printf("Failed to Unmarshal: %s", err.Error())
		return
	}

	fmt.Printf("Username: %s", mixedStruct.Username)    // -> Username: Foo
	fmt.Printf("Age: %d", mixedStruct.Data.Age)         // -> Age: 30
	fmt.Printf("Country: %s", mixedStruct.Data.Country) // -> Country: Sweden
}
