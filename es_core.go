package main

import "fmt"
import "github.com/op/go-nanomsg"

func main() {
	fmt.Printf("Hello, game!\n")
	
	game_socket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
        if err != nil {
                panic(err)
        }
        _, err = game_socket.Connect("tcp://127.0.0.1:60206")
        if err != nil {
                panic(err)
        }
}
