package main

import(
    "net"
    "log"
    "flag"
    "io"
    "fmt"
    "encoding/json"
)

const (
    HOST="localhost"
    PORT="3000"
    TYPE="tcp"
)

type Message struct {
    Body string
}

var addr = flag.String("addr", HOST+":"+PORT, "host:port - Connect to server")

func main() {
    flag.Parse()
    // New server
    server, err := net.Listen(TYPE, *addr)
    if err != nil {
        log.Fatal(err)
    }
    defer server.Close()
    // Infinite loop waiting for user connection
    for {
        conn, err := server.Accept()
        if err != nil {
            log.Fatal(err)
        }

        go request(conn)
    }
}

func request(conn net.Conn) {
    defer conn.Close()

    decoder := json.NewDecoder(conn)
    m := Message{}
    for {
        if err := decoder.Decode(&m); err != nil {
            log.Fatal(err)
        }
        fmt.Println(m.Body)
    }

    fmt.Fprintln(conn, "Welcome on the tchat !")

    io.Copy(conn, conn)
    conn.Close()
}
