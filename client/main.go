package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "encoding/json"
    "flag"
    "net"
    "io"
)

const (
    HOST = "localhost"
    PORT = "3000"
    METHOD = "tcp"
)
// Définir un flag "-addr" dont la valeur par défaut est "localhost:3000"
var addr = flag.String("addr", HOST+":"+PORT, "host:port - Connect to server")

type Message struct {
    Body string
}

func main() {

    flag.Parse()
    fmt.Println(*addr)

    // Définir un nouveau listener sur l'entrée standard
    s := bufio.NewScanner(os.Stdin)
    fmt.Println("Veuillez écrire votre message :")

    c, err := net.Dial(METHOD, *addr)

    m := Message{}
    for s.Scan() {
        // Définir la sortie d'ecriture sur os.Stdout
        e := json.NewEncoder(c)
        // On enregistre le message dans le Message.Body
        m.Body = s.Text()
        // Encoder m
        e.Encode(m)
        break;
    }
    if err := s.Err(); err != nil {
        log.Fatal(err)
    }

    if err != nil {
        log.Fatal(err)
    }
    //c.Write(json.March)


    io.Copy(os.Stdout, c)
    c.Close()
    fmt.Println("okokok")
}
