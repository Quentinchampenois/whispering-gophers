package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "encoding/json"
)

type Message struct {
    Body string
}

func main() {
    // Définir un nouveau listener sur l'entrée standard
    s := bufio.NewScanner(os.Stdin)

    for s.Scan() {
        fmt.Println("Veuillez écrire votre message : \n")
        // Définir la sortie d'ecriture sur os.Stdout
        e := json.NewEncoder(os.Stdout)
        // On enregistre le message dans le Message.Body
        m := Message{Body:s.Text()}
        // Encoder m 
        e.Encode(m)
    }
    if err := s.Err(); err != nil {
        log.Fatal(err)
    }
}
