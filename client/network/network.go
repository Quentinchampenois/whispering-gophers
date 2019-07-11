// network package allows client communicate with server
package network

import (
    "log"
    "net"
    "io"
    "os"
)

// TCPConnection make a tcp connexion on server, passing data
func TCPConnection( method string, url string, data []byte ) {
    c, err := net.Dial(method, url)
    if err != nil {
        log.Fatal(err)
    }
    c.Write(data)
    io.Copy(os.Stdout, c)
    defer c.Close()
}
