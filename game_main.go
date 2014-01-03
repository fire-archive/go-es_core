package main

import ("fmt"
	"github.com/op/go-nanomsg"
	"github.com/fire/go-sdl2/sdl"
	)

func gameThread() {
	controlSocket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
	if err != nil {
                panic(err)
        }
        _, err = controlSocket.Connect("tcp://127.0.0.1:60206")
        if err != nil {
                panic(err)
        }

	renderSocket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
	if err != nil {
                panic(err)
        }
        _, err = renderSocket.Bind("tcp://127.0.0.1:60210")
        if err != nil {
                panic(err)
        }

	inputMouseSub, err := nanomsg.NewSubSocket()
	inputMouseSub.Subscribe("input.mouse:")
	if err != nil {
		panic(err)
	}
	_, err = inputMouseSub.Connect("tcp://127.0.0.1:60208")
	if err != nil {
		panic(err)
	}
	
	inputKbSub, err := nanomsg.NewSubSocket()
	inputKbSub.Subscribe("input.kb:")
	if err != nil {
		panic(err)
	}
	_, err = inputKbSub.Connect("tcp://127.0.0.1:60208")

	gameInit()
	baseLine := sdl.GetTicks()
	fmt.Printf("baseline: %d\n", baseLine)
}
