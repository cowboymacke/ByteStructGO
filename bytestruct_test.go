package bytestruct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

type StructTwo struct {
	Value    int8
	ValueTwo uint8  `byteLength:"Name"`
	Name     string `byte:"ValueTwo"`
}

type StructOne struct {
	Length  uint16
	Payload StructTwo
}

func Test_Marshal(t *testing.T) {
	data, err := Marshal(binary.BigEndian, StructOne{
		Length: 2,
		Payload: StructTwo{
			Value:    1,
			ValueTwo: 2,
			Name:     "Hello world",
		},
	})

	if err != nil {
		t.Errorf("faild with error %w", err)
	}

	fmt.Printf("sucess with data : %d", data)

}

func Test_unmarshalt(t *testing.T) {

	var data StructOne

	reader := bytes.NewBuffer([]byte{0, 2, 1, 11, 72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100})
	err := Unmarshal(reader, binary.BigEndian, &data)

	if err != nil {
		t.Errorf("faild with error %w", err)
	}

	fmt.Println("sucess")
}
