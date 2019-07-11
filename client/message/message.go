// message package allows to manage messages
package message

import (
    "fmt"
    "log"
    "bufio"
)

// UserReaderMessage ask to the user to write a message
// It returns the message ( without "\n" ) and error
func UserReaderMessage(r *bufio.Reader) (userMsg string, err error) {
    fmt.Println("Saisissez votre message : ")
    userMsg, err = r.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }
    userMsg = userMsg[:len(userMsg) - 1]
    return
}
