package bytestruct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

type Simple struct {
	Length uint16 `bytePayload:"Name"`
	Name   string `byteSize:"Length"`
}

func Test_simple(t *testing.T) {

	name := "Hello World"

	data, err := Marshal(binary.BigEndian, Simple{
		Name: name,
	})

	fmt.Println(data)
	if err != nil {
		t.Errorf("faild with error %e", err)
	}

	var umData Simple
	reader := bytes.NewBuffer(data)
	err = Unmarshal(reader, binary.BigEndian, &umData)
	if err != nil {
		t.Errorf("faild with error %e", err)
	}

	nameLen, err := intCaster(len(name), reflect.TypeOf(umData.Length).Kind())
	if err != nil {
		t.Errorf("faild with error %e", err)
	}

	if umData.Length != nameLen {
		t.Errorf("failed with DataLength")
	}

	if umData.Name != name {
		t.Errorf("failed with name")
	}
}

type StructTwo struct {
	Value      int8
	TextLength uint8  `bytePayload:"Name"`
	Name       string `byteSize:"TextLength"`
}

type StructOne struct {
	Type    uint8
	Length  uint16    `bytePayload:"Payload"`
	Payload StructTwo `byteSize:"Length"`
}

func Test_marshal(t *testing.T) {
	testData := StructOne{
		Type: 25,
		Payload: StructTwo{
			Value: 1,
			Name:  "HelloBigWorld",
		},
	}
	data, err := Marshal(binary.BigEndian, testData)

	if err != nil {
		t.Errorf("faild with error %e", err)
	}

	var umData StructOne
	reader := bytes.NewBuffer(data)
	err = Unmarshal(reader, binary.BigEndian, &umData)
	if err != nil {
		t.Errorf("faild with error %e", err)
	}

	if umData.Type != testData.Type {
		t.Errorf("Marsheld failed with Type data ")
	}

	if umData.Payload.Value != testData.Payload.Value {
		t.Errorf("Marsheld failed with Type Value ")
	}

	if umData.Payload.Name != testData.Payload.Name {
		t.Errorf("Marsheld failed with Type Name ")
	}

}

type ArrayData struct {
	Type   uint8
	Length uint8  `bytePayload:"Name"`
	Name   string `byteSize:"Length"`
}

type StructArray struct {
	LengthUint    uint16      `bytePayload:"PayloadUint"`
	PayloadUint   []uint8     `byteSize:"LengthUint"`
	LengthStruct  uint16      `bytePayload:"PayloadStruct"`
	PayloadStruct []ArrayData `byteSize:"LengthStruct"`
}

func Test_marshal_array(t *testing.T) {
	testData := StructArray{
		PayloadStruct: []ArrayData{{Type: 1, Name: "Hello"}, {Type: 2, Name: "world"}},
		PayloadUint:   []uint8{8, 2, 3, 5, 1},
	}

	data, err := Marshal(binary.BigEndian, testData)

	if err != nil {
		t.Errorf("faild with error %e", err)
		return
	}

	var umData StructArray
	reader := bytes.NewBuffer(data)
	err = Unmarshal(reader, binary.BigEndian, &umData)

	if err != nil {
		t.Errorf("faild with error %e", err)
	}

}
