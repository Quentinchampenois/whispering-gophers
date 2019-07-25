package main

//  Chaque utilisateur à un serveur de réception et communique sur le serveur de l'autre

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "encoding/json"
    "flag"
    "net"
    "io"
    "github.com/campoy/whispering-gophers/util"
)

const (
    HOST = "localhost"
    PORT = "3000"
    METHOD = "tcp"
)

// Définir un flag "-addr" dont la valeur par défaut est "localhost:3000"
var (
    addr = flag.String("addr", HOST+":"+PORT, "host:port - Connect to server")
	dialAddr   = flag.String("dial", HOST+":3001", "host:port listener server")
    self string;
)
type Message struct {
    Addr string
    Body string
}

func main() {
    flag.Parse()


    // New server
    server, err := util.Listen()
    if err != nil {
        log.Println(err)
        return
    }
    self = server.Addr().String()
    log.Println("Listening on ", self)

    go dial(*dialAddr, self)
    defer server.Close()
    // Infinite loop waiting for user connection
    for {
        conn, err := server.Accept()
        if err != nil {
            log.Println(err)
            return
        }

        go request(conn, server.Addr().String())
    }
}

func request(conn net.Conn, lAddr string) {
    defer conn.Close()

    decoder := json.NewDecoder(conn)
    m := Message{}
    for {
        if err := decoder.Decode(&m); err != nil {
            log.Println(err)
            return
        }
        fmt.Println("Message reçu !")
        fmt.Println(m.Addr)
        fmt.Println(m.Body)

    }

    io.Copy(conn, conn)
    dial(m.Addr, lAddr)
    conn.Close()
}

func dial(addr, lAddr string) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
        return
	}

fmt.Println(lAddr)
	s := bufio.NewScanner(os.Stdin)
	e := json.NewEncoder(c)
    m := Message{Addr: lAddr}
	for s.Scan() {
		m.Body = s.Text()
		err := e.Encode(m)
		if err != nil {
			log.Println(err)
            return
		}
	}
	if err := s.Err(); err != nil {
		log.Println(err)
        return
	}
}
