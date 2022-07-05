package bytestruct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

type StructTwo struct {
	Value    int8
	ValueTwo int8
	name     string `byte:"ValueTwo"`
}

type StructOne struct {
	Length  int16
	Payload StructTwo `byte:"Length"`
}

func Test_Marshal(t *testing.T) {
	data, err := Marshal(binary.BigEndian, StructOne{
		Length: 2,
		Payload: StructTwo{
			Value:    1,
			ValueTwo: 2,
			name:     "he",
		},
	})

	if err != nil {
		t.Errorf("faild with error %w", err)
	}

	fmt.Printf("sucess with data : %d", data)

}

func Test_unmarshalt(t *testing.T) {

	var data StructOne

	reader := bytes.NewBuffer([]byte{0, 2, 1, 2, 104, 101})
	err := Unmarshal(reader, binary.BigEndian, &data)

	if err != nil {
		t.Errorf("faild with error %w", err)
	}

	fmt.Println("sucess")
}
