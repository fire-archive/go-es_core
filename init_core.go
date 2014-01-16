
package core

import ("fmt"
	"os"
	"strconv"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/op/go-nanomsg"
	"github.com/fire/go-ogre3d"
	"github.com/jmckaskill/go-capnproto"
	"runtime"
	"math")

type InputState struct {
	yawSens float32
	pitchSens float32
	orientationFactor float32 // +1/-1 easy switch between look around and manipulate something
	yaw float32 // degrees, modulo [-180,180] range
	pitch float32 // degrees, clamped [-90,90] range
	roll float32
	orientation ogre.Quaternion // current orientation
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
	if runtime.GOOS == "windows" {
		root.LoadPlugin(wd  + "/RenderSystem_GL3Plus")
		}
	if runtime.GOOS == "darwin" {
		root.LoadPlugin(wd  + "/../frameworks/RenderSystem_GL3Plus")
	}
				
	renderers := root.GetAvailableRenderers()
	if renderers.RenderSystemListSize() != 1 {
		
		panic(fmt.Sprintf("Failed to initalize RendererRenderSystem_GL"))
	}
	root.SetRenderSystem(renderers.RenderSystemListGet(0))
	root.Initialise(false, "es_core::ogre")
	params := ogre.CreateNameValuePairList()
	if runtime.GOOS == "windows" {
		windowsInfo := info.GetWindowsInfo()
		windowString := strconv.FormatUint(uint64(*(*uint32)(windowsInfo.Window)), 10)
		params.AddPair("parentWindowHandle", windowString)
	}
	if runtime.GOOS == "darwin" {
		params.AddPair("macAPI", "cocoa")
		cocoaInfo := info.GetCocoaInfo()
		windowString := strconv.FormatUint(uint64(*(*uint32)(cocoaInfo.Window)), 10)
		params.AddPair("parentWindowHandle", windowString)
	}
	
	renderWindow := root.CreateRenderWindow("es_core::ogre", 800, 600, false, params)
//	renderWindow.SetVisible(true)
	
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
		var b []byte
		// We wait here.
		b, err = nnInputPull.Recv(0)	
		s, _, err := capn.ReadFromMemoryZeroCopy(b)
		if err != nil {
			fmt.Printf("Read error %v\n", err)
			return
		}	
		state := ReadRootState(s)
		
		fmt.Printf("Game push received:\n")
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		// poll for events before processing the request
		// NOTE: this is how SDL builds the internal mouse and keyboard state
		// TODO: done this way does not meet the objectives of smooth, frame independent mouse view control,
		// Plus it throws some latency into the calling thread

		var event sdl.Event
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent {
			switch t := event.(type) {
			case *sdl.KeyDownEvent:
				fmt.Printf("SDL keyboard event:\n")
			case *sdl.KeyUpEvent:
				fmt.Printf("SDL keyboard event:\n")
				if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
					// Todo
					sendShutdown(nnRenderSocket, nnGameSocket)
					shutdownRequested = true
				}
			case *sdl.MouseMotionEvent:
				// + when manipulating an object, - when doing a first person view .. needs to be configurable?
				is.yaw += is.orientationFactor * is.yawSens * float32(t.XRel)
				if is.yaw >= 0.0 {
					is.yaw = float32(math.Mod(float64(is.yaw) + 180.0, 360.0) - 180.0)
				} else {
					is.yaw = float32(math.Mod(float64(is.yaw) - 180.0, 360.0) + 180.0)
				}
				// + when manipulating an object, - when doing a first person view .. needs to be configurable?
				is.pitch += is.orientationFactor * is.pitchSens * float32(t.YRel)
				if is.pitch > 90.0 {
					is.pitch = 90.0
				} else if ( is.pitch < -90.0 ) {
					is.pitch = -90.0
				}
				// build a quaternion of the current orientation
				var r ogre.Matrix3
				r.FromEulerAnglesYXZ( deg2Rad(is.yaw), deg2Rad(is.pitch), deg2Rad(is.roll)) 
				is.orientation.FromRotationMatrix(r)
			case *sdl.MouseButtonEvent:
				fmt.Printf("SDL mouse button event:\n")
			case *sdl.QuitEvent:
			    // push a shutdown on the control socket, game and render will pick it up later
				// NOTE: if the message patterns change we may still have to deal with hangs here
				sendShutdown(nnRenderSocket, nnGameSocket)
				
				shutdownRequested = true
			default:
				fmt.Printf("SDL_Event %T\n", event);
			}
		}
		switch {
		case state.Mouse():
			
		}
	}
}

func deg2Rad(deg float32) float32 {
	return deg * math.Pi / 180
}

func sendShutdown(nnRenderSocket *nanomsg.Socket, nnGameSocket *nanomsg.Socket) {
	fmt.Printf("Render socket shutdown.\n")
	fmt.Printf("Game socket shutdown.\n")
}
