### ByteStructGO
:construction: This is **not** production ready :construction:

ByteStructGO is a Go library that simplifies the reading and writing of binary data. It allows you to define the structure of your data using tags, and then automatically handles the encoding and decoding of bytes. ByteStructGO supports both little-endian and big-endian byte orders.

## Features
-  Easy to use: just define your data structure using tags, and let ByteStructGO do the rest.
-  Nested Structure: ByteStructGO can decode nessted and array structs

## Missing
- cant handle array of strings


## Install

```sh
go get github.com/MarcusLing/ByteStructGO
```

## Usage

To use ByteStructGO, you need to import the package and define your data structure using tags.


Offset |  Size	| Description
0 | 2|  payloadSize (uint16) [0x0 0xb]
2 | x |  payload (string) [0x48 0x65 0x6c 0x6c 0x6f 0x20 0x57 0x6f 0x72 0x6c 0x64]


```go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	structGo "github.com/MarcusLing/ByteStructGO"
)

// Define the message structure
type Message struct {
	Length  uint16 
	Payload string `byteSize:"Length"`
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

}

```




## Examples
You can find more examples of using ByteStructGO in the [examples] folder.

