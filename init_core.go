
package core

import ("fmt"
	"os"
	"strconv"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/op/go-nanomsg"
	"github.com/fire/go-ogre3d")

type InputState struct {
	yawSens float32
	pitchSens float32
	orientationFactor float32 // +1/-1 easy switch between look around and manipulate something
	yaw float32 // degrees, modulo [-180,180] range
	pitch float32 // degrees, clamped [-90,90] range
	roll float32
	// orientation ogre.Quaternion // current orientation
}

func InitCore() {
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
	defer sdl.Quit()
	var info sdl.SysWMInfo 
	if !window.GetWMInfo(&info) {
		panic(fmt.Sprintf("window.GetWMInfo failed.\n"))
	}
	// Parse and print info's version
	// Parse and print info's SYSWM_TYPE
	root := ogre.NewRoot("", "", "ogre.log")
	defer root.Destroy()
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
	
	nnGameSocket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
        if err != nil {
                panic(err)
        }
        _, err = nnGameSocket.Bind("tcp://127.0.0.1:60206")
        if err != nil {
                panic(err)
        }
	
	nnRenderSocket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
	if err != nil {
                panic(err)
        }
        _, err = nnRenderSocket.Bind("tcp://127.0.0.1:60207")
        if err != nil {
                panic(err)
        }

	nnInputPub, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.PUB)
        if err != nil {
                panic(err)
        }
        _, err = nnInputPub.Bind("tcp://127.0.0.1:60208")
        if err != nil {
                panic(err)
        }

	nnInputPull, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.PULL)
        if err != nil {
                panic(err)
        }
        _, err = nnInputPull.Bind("tcp://127.0.0.1:60209")
        if err != nil {
                panic(err)
        }
	go gameThread()
	var renderThreadParams RenderThreadParams
	renderThreadParams.root = root
	renderThreadParams.window = window
	renderThreadParams.ogreWindow = renderWindow
	
	go renderThread(renderThreadParams)

	window.SetGrab(true)
	sdl.SetRelativeMouseMode(true)

	shutdownRequested := false
	var is InputState
	is.yawSens = 0.1
	is.yaw = 0.0
	is.pitchSens = 0.1
	is.pitch = 0.0
	is.roll = 0.0
	is.orientationFactor = -1.0 // Look around config

	for !shutdownRequested {
		var inputPull string
		string, err := nnInputPull.RecvString()
	}
}
