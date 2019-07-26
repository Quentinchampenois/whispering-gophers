package main

//  Chaque utilisateur à un serveur de réception et communique sur le serveur de l'autre

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/campoy/whispering-gophers/util"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var (
	dialAddr = flag.String("dial", "", "host:port listener server")
	self     string
	// Si le channel est unbuffered c'est synchrone, à l'inverse c'est asynchrone
	ch = make(chan Message)
)

type Message struct {
	Addr string
	Body string
}

type Peers struct {
	m map[string]chan<- Message
	mu sync.RWMutex
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

	go dial(*dialAddr)
	go readUserMsg()

	defer server.Close()
	// Infinite loop waiting for user connection
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
			return
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
			log.Println(err)
			return
		}
		fmt.Println("Message reçu !")
		fmt.Println(m.Addr)
		fmt.Println(m.Body)

	}

	io.Copy(conn, conn)
	conn.Close()
}

func dial(addr string) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}

	// Définir où le json doit être écrit
	e := json.NewEncoder(c)

	for m := range ch {
		err := e.Encode(m)
		if err != nil {
			log.Fatal(err)
		}
	}


}

func readUserMsg() {

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := Message {
			Addr: self,
			Body: s.Text(),
		}
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}

		// Envoyer le struct Message sur le chan ch de type Message
		ch <- m
	}

}

func (p *Peers) Add(addr string)chan<- Message {
	defer p.mu.Unlock()
	p.mu.Lock()
	if p.m[addr] != nil {
		return nil
	}
	msgCh := make(chan Message)
	p.m[addr] = msgCh
	return msgCh
}

func (p *Peers) Remove(key string) {
	defer p.mu.Unlock()
	p.mu.Lock()
	delete(p.m, key)
}

func (p *Peers) List() []chan<- Message {
	defer p.mu.Unlock()
	p.mu.Lock()
	// Instanciation d'une variable slice de type []chan<- Message, de taille 0 et de capacité maximale de la taille des adresses
	slice := make([]chan<- Message, 0, len(p.m))

	for _, ch := range p.m {
		slice = append(slice, ch)
	}

	return slice
}