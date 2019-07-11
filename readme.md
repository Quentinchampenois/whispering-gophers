# Activity Whispering Gophers

This repo contains a network programming with Go where users can transmit messages on a peer-to-peer network.
Source: https://whispering-gophers.appspot.com/talk.slide#1


## Project architecture

* client/ => client part
    * message/ => message manager
    * network/ => client network manager     


* server/ => server part

## Getting started

Just fetch go packages from this repository, build them if needed, and run them

### Get the go packages

- Server

    `go get github.com/quentinchampenois/whispering-gophers/server`


- Client

    `go get github.com/quentinchampenois/whispering-gophers/client`

### Run go programs

Run on different terminal windows, __one server__ and every client you want.
