package main

import ("fmt"
	"github.com/fire/go-sdl2/sdl"
	"github.com/op/go-nanomsg")

func main() {
	fmt.Printf("Hello, game!\n")
	
	sdl.Init(sdl.INIT_EVERYTHING)
	window := sdl.CreateWindow("es_core::SDL",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_OPENGL | sdl.WINDOW_SHOWN)
	if window == nil {
		panic(fmt.Sprintf("SDL_CreateWindow failed: %s\n", sdl.GetError()))
	}
	
	game_socket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
        if err != nil {
                panic(err)
        }
        _, err = game_socket.Connect("tcp://127.0.0.1:60206")
        if err != nil {
                panic(err)
        }
}
