package main

//  Chaque utilisateur à un serveur de réception et communique sur le serveur de l'autre

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/campoy/whispering-gophers/util"
	"log"
	"net"
	"os"
	"sync"
)

var (
	dialAddr = flag.String("dial", "", "host:port listener server")
	self     string
	peers = &Peers{m: make(map[string]chan<- Message)}
)

type Message struct {
	ID string
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

	if *dialAddr != "" {
		go dial(*dialAddr)
	}
	go readUserMsg()
	
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
			log.Println(err)
			return
		}
		if Seen(m.ID) {
			continue
		} 
		fmt.Println("Message reçu !")
		fmt.Println(m.Addr)
		fmt.Println(m.Body)
		go dial(m.Addr)
	}

}

func dial(addr string) {

	if addr == self {
		return
	}

	ch := peers.Add(addr)
	if ch == nil {
		return
	}
	defer peers.Remove(addr)

	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(addr, err)
		return
	}

	defer c.Close()

	// Définir où le json doit être écrit
	e := json.NewEncoder(c)

	for m := range ch {
		err := e.Encode(m)
		if err != nil {
			log.Fatal(addr, err)
			return
		}
	}
}

func broadcast(m Message) {
	for _, ch := range peers.List() {
		select {
		case ch <- m:
			fmt.Println("Sent message ", m.Body)
		default:
			fmt.Println("Message not sent")
		}
	}
}

func readUserMsg() {

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := Message {
			ID: util.RandomID(),
			Addr: self,
			Body: s.Text(),
		}
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
		Seen(m.ID)
		// Envoyer le struct Message sur le chan ch de type Message
		broadcast(m)
	}

}

func (p *Peers) Add(addr string) <-chan Message {
	defer p.mu.Unlock()
	p.mu.Lock()
	if _, ok := p.m[addr]; ok {
		return nil
	}
	msgCh := make(chan Message)
	p.m[addr] = msgCh

	fmt.Println(len(p.m))
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
	var slice []chan<- Message

	for _, ch := range p.m {
		slice = append(slice, ch)
	}

	return slice
}

var seenIDs = struct{
	m map[string]bool
	sync.Mutex
}{m: make(map[string]bool)}

func Seen(id string) bool {
	seenIDs.Lock()
	ok := seenIDs.m[id]
	seenIDs.m[id] = true
	seenIDs.Unlock()
	return ok
}
