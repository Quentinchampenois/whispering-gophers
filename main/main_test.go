package main

import (
	"sync"
	"testing"
	"time"
)

type Message struct {
	Addr string
	Body string
}

type Peers struct {
	m map[string]chan<- Message
	mu sync.RWMutex
}

func (p *Peers) Add(addr string)<-chan Message {
	if p.m[addr] != nil {
		return nil
	}
	msgCh := make(chan Message)
	p.m[addr] = msgCh
	return msgCh
}

func (p *Peers) Remove(key string) {
	delete(p.m, key)
}

func (p *Peers) List() []chan<- Message {
	// Instanciation d'une variable slice de type []chan<- Message, de taille 0 et de capacité maximale de la taille des adresses
	slice := make([]chan<- Message, 0, len(p.m))

	for _, ch := range p.m {
		slice = append(slice, ch)
	}

	return slice
}


func TestPeers(t *testing.T) {
	peers := &Peers{m: make(map[string]chan<- Message)}
	done := make(chan bool, 1)

	var chA, chB <-chan Message
	go func() {
		defer func() { done <- true }()
		if chA = peers.Add("a"); chA == nil {
			t.Fatal(`peers.Add("a") returned nil, want channel`)
		}
	}()
	go func() {
		defer func() { done <- true }()
		if chB = peers.Add("b"); chB == nil {
			t.Fatal(`peers.Add("b") returned nil, want channel`)
		}
	}()
	<-done
	<-done
	if chA == chB {
		t.Fatal(`peers.Add("b") returned same channel as "a"!`)
	}
	if ch := peers.Add("a"); ch != nil {
		t.Fatal(`second peers.Add("a") returned non-nil channel, want nil`)
	}
	if ch := peers.Add("b"); ch != nil {
		t.Fatal(`second peers.Add("b") returned non-nil channel, want nil`)
	}

	list := peers.List()
	if len(list) != 2 {
		t.Fatalf("peers.List() returned a list of length %d, want 2", len(list))
	}

	go func() {
		for _, ch := range list {
			select {
			case ch <- Message{Body: "foo"}:
			case <-time.After(10 * time.Millisecond):
			}
		}
		done <- true
	}()
	select {
	case m := <-chA:
		if m.Body != "foo" {
			t.Fatal("received message %q, want %q", m.Body, "foo")
		}
	case <-done:
		t.Fatal(`didn't receive message on "a" channel`)
	}
	<-done

	peers.Remove("a")

	list = peers.List()
	if len(list) != 1 {
		t.Fatalf("peers.List() returned a list of length %d, want 1", len(list))
	}

	go func() {
		select {
		case list[0] <- Message{Body: "bar"}:
		case <-time.After(10 * time.Millisecond):
		}
		done <- true
	}()
	select {
	case m := <-chB:
		if m.Body != "bar" {
			t.Fatalf("received message %q, want %q", m.Body, "bar")
		}
	case <-done:
		t.Fatal(`didn't receive message on "b" channel`)
	}
}
