package main

import ("fmt"
	"os"
	"github.com/fire/go-sdl2/sdl"
	"github.com/op/go-nanomsg"
	"github.com/fire/go-ogre3d")

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
		panic(fmt.Sprintf("sdl.CreateWindow failed: %s\n", sdl.GetError()))
	}
	glcontext := sdl.GL_CreateContext(window)
	if glcontext == nil {
		panic(fmt.Sprintf("sdl.CreateContext failed: %s\n", sdl.GetError()))
	}
	// sdl.GetWindowWMInfo

	root := ogre.NewRoot("", "", "ogre.log")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root.LoadPlugin(wd  + "/../frameworks/RenderSystem_GL.framework")
	renderers := root.GetAvailableRenderers()
	if renderers.RenderSystemListSize() != 1 {
		panic(fmt.Sprintf("Failed to initalize RendererRenderSystem_GL"))
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
