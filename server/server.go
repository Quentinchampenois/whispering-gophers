package main

import(
    "net"
    "log"
    "encoding/json"
    "fmt"
    "util"
)

const (
    HOST="localhost"
    PORT="3000"
    TYPE="tcp"
)

type Message struct {
    Addr string
    Body string
}

// handleRequest handles incoming requests from client to server
func handleRequest(conn net.Conn) {

   // Creates json decoder on the connection
   d := json.NewDecoder(conn)
   var msg Message
   // Write message in Message.Body
   err := d.Decode(&msg)
     if err != nil {
         log.Println("Error reading:", err)
     }
   fmt.Println(msg.Body)
   // Write the message on user interface
     _, err = conn.Write([]byte(msg.Body + "\n"))
     if err != nil {
         log.Println("Error writing:", err)
     }

  defer conn.Close()
}

func main() {
    // New server
    server, err := net.Listen(TYPE, HOST+":"+PORT)
    if err != nil {
        log.Fatal(err)
    }
    defer server.Close()
    addr, _ := util.Listen()
    fmt.Println(addr.Addr())
    // Infinite loop waiting for user connection
    for {
        conn, err := server.Accept()
        if err != nil {
            log.Fatal(err)
        }
        // Define new goroutine foreach connection
        go handleRequest(conn)
    }
}
