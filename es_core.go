package main

import ("fmt"
	"os"
	"strconv"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/op/go-nanomsg"
	"github.com/fire/go-ogre3d")

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	window := sdl.CreateWindow("es_core::SDL",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_SHOWN)
	if window == nil {
		panic(fmt.Sprintf("sdl.CreateWindow failed: %s\n", sdl.GetError()))
	}
	var info sdl.SysWMInfo 
	if !window.GetWMInfo(&info) {
		panic(fmt.Sprintf("window.GetWMInfo failed.\n"))
	}
	// Parse and print info's version
	// Parse and print info's SYSWM_TYPE
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
	root.SetRenderSystem(renderers.RenderSystemListGet(0))
	root.Initialise(false, "es_core::ogre")
	params := ogre.CreateNameValuePairList()
	params.AddPair("macAPI", "cocoa")
	cocoaInfo := info.GetCocoaInfo()
	windowString := strconv.FormatUint(uint64(*(*uint32)(cocoaInfo.Window)), 10)
	params.AddPair("parentWindowHandle", windowString)
	
	renderWindow := root.CreateRenderWindow("es_core::ogre", 800, 600, false, params)
	renderWindow.SetVisible(true)
	
	game_socket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
        if err != nil {
                panic(err)
        }
        _, err = game_socket.Connect("tcp://127.0.0.1:60206")
        if err != nil {
                panic(err)
        }
	
	render_socket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
	if err != nil {
                panic(err)
        }
        _, err = render_socket.Connect("tcp://127.0.0.1:60207")
        if err != nil {
                panic(err)
        }

	input_pub, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.PUB)
        if err != nil {
                panic(err)
        }
        _, err = input_pub.Connect("tcp://127.0.0.1:60208")
        if err != nil {
                panic(err)
        }

	input_pull, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.PULL)
        if err != nil {
                panic(err)
        }
        _, err = input_pull.Connect("tcp://127.0.0.1:60209")
        if err != nil {
                panic(err)
        }

	go gameThread()
	go renderThread()
}
