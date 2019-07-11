package main

import (
    "os"
    "bufio"
    "encoding/json"
    "log"
    "github.com/quentinchampenois/whispering-gophers/client/message"
    "github.com/quentinchampenois/whispering-gophers/client/network"
)

type Message struct {
    Body string
}

func main() {
    // Define new reader on standard input
    reader := bufio.NewReader(os.Stdin)
    // Get user input from previous reader
    userMsg, err := message.UserReaderMessage(reader)
    // Message saved in Message.Body
    msg := Message{Body: userMsg}
    // JSON encoding
    jsondata, err := json.Marshal(msg)
    if err != nil {
        log.Fatal(err)
    }

    // Make tcp connection on localhost:3000 passing jsondata
    network.TCPConnection("tcp", "localhost:3000", jsondata)
}
