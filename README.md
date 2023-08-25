### ByteStructGO: A Go library for reading and writing binary data using go structures
:construction: This is **not** production ready :construction:

ByteStructGO is a Go library that simplifies the reading and writing of binary data. It allows you to define the structure of your data using tags, and then automatically handles the encoding and decoding of bytes. ByteStructGO supports both little-endian and big-endian byte orders.

## Features
-  Easy to use: just define your data structure using tags, and let ByteStructGO do the rest.
-  Nested Structure: ByteStructGO can decode nessted and array structs

## Missing
- cant handle array of strings


## Install

```sh
go get github.com/cowboymacke/ByteStructGO
```

## Usage

To use ByteStructGO, you need to import the package and define your data structure using tags.


```go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	structGo "github.com/cowboymacke/ByteStructGO"
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


## Feedback
If you have any questions, suggestions, or issues, please feel free to open an issue or submit a pull request.




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->


[examples]:https://github.com/MarcusLing/ByteStructGO/tree/main/examples
