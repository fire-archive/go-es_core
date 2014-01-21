
package core

import ("fmt"
	"os"
	"strconv"
	"runtime"
	"math"
	"bytes"
	"time"
	"reflect"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/op/go-nanomsg"
	"github.com/fire/go-ogre3d"
	"github.com/jmckaskill/go-capnproto")

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
	var gameThreadParams GameThreadParams
	gameThreadParams.start = time.Now() // There's an small time before this variable is initalized,
	// it probably doesn't matter... Someone timed Go initalization at 1.94us on Linux.
	sdl.Init(sdl.INIT_EVERYTHING)
	window := sdl.CreateWindow("es_core::SDL",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN)
	if window == nil {
		panic(fmt.Sprintf("sdl.CreateWindow failed: %s\n", sdl.GetError()))
	}
	defer sdl.Quit()
	glContext := sdl.GL_CreateContext(window)
	
	var info sdl.SysWMInfo 
	if !window.GetWMInfo(&info) {
		panic(fmt.Sprintf("window.GetWMInfo failed.\n"))
	}
	var version sdl.Version
	sdl.GetVersion(&version)
	
	fmt.Printf("Sdl Major Version: %d\n", version.Major)
	fmt.Printf("Sdl Minor Version: %d\n", version.Minor)
	fmt.Printf("Sdl Patch level: %d\n", version.Patch)
	fmt.Printf("Sdl Subsystem: %s\n", getSubsystemString(info)) 
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
		params.AddPair("externalGLControl", "1")
		params.AddPair("externalGLContext", strconv.FormatUint(uint64(uintptr(glContext)), 10))
		
		windowsInfo := info.GetWindowsInfo()
		windowString := strconv.FormatUint(uint64(uintptr(windowsInfo.Window)), 10)
		params.AddPair("externalWindowHandle", windowString)
	}
	if runtime.GOOS == "darwin" {
		params.AddPair("externalGLControl", "1")
		params.AddPair("externalGLContext", strconv.FormatUint(uint64(uintptr(glContext)), 10))
		params.AddPair("macAPI", "cocoa")
		cocoaInfo := info.GetCocoaInfo()
		params.AddPair("externalWindowHandle", strconv.FormatUint(uint64(*(*uint32)(cocoaInfo.Window)), 10))
	}
	
	renderWindow := root.CreateRenderWindow("es_core::ogre", 800, 600, false, params)
	if runtime.GOOS == "darwin" {
		// I suspect triple buffering is on by default, which makes vsync pointless?
		// except maybe for poorly implemented render loops which will then be forced to wait
		// window->SetVSyncEnabled(false)
	} else {
		// NOTE: SDL_GL_SWAP_CONTROL was SDL 1.2 and has been retired
		sdl.GL_SetSwapInterval(1);
	}
	
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
	go gameThread(gameThreadParams)
	var renderThreadParams RenderThreadParams
	renderThreadParams.start = gameThreadParams.start
	renderThreadParams.root = root
	renderThreadParams.window = window
	renderThreadParams.glContext = glContext
	renderThreadParams.ogreWindow = renderWindow
	
	sdl.GL_MakeCurrent(window, nil)
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

	for !shutdownRequested /* && SDL_GetTicks() < MAX_RUN_TIME */ {
		// We wait here.
		b, err := nnInputPull.Recv(0)
		if err != nil {
			fmt.Printf("%s\n", err)
		}	
		s, _, err := capn.ReadFromMemoryZeroCopy(b)
		if err != nil {
			fmt.Printf("Read error %v\n", err)
			return
		}	
		state := ReadRootState(s)
		fmt.Printf("Game push received:\n")
		// poll for events before processing the request
		// NOTE: this is how SDL builds the internal mouse and keyboard state
		// TODO: done this way does not meet the objectives of smooth, frame independent mouse view control,
		// Plus it throws some latency into the calling thread

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			default:
				fmt.Printf("SDL_Event: %s\n", reflect.TypeOf(t).String()[5:]);
			case *sdl.KeyDownEvent:
				fmt.Printf("SDL keyboard event:\n")
			case *sdl.KeyUpEvent:
				fmt.Printf("SDL keyboard event:\n")
				if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
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
			}
		}
		switch {
		// we are ready to process the request now
		case state.Mouse():
			buttons := sdl.GetMouseState(nil, nil)
			fmt.Printf("buttons: %d\n", buttons)
			s := capn.NewBuffer(nil)
			ms := NewRootInputMouse(s)
			ms.SetW(is.orientation.W())
			ms.SetX(is.orientation.X())
			ms.SetY(is.orientation.Y())
			ms.SetZ(is.orientation.Z())
			ms.SetButtons(buttons)
			buf := bytes.Buffer{}
			s.WriteTo(&buf)
			nnInputPub.Send(append([]byte("input.mouse:"), buf.Bytes()...), 0)
			fmt.Printf("Mouse input sent.\n")
			
		case state.Kb():
		// looking at a few hardcoded keys for now
		// NOTE: I suspect it would be perfectly safe to grab that pointer once, and read it from a different thread?
			state := sdl.GetKeyboardState()
			t := capn.NewBuffer(nil)
			kbs := NewRootInputKb(t)			
			kbs.SetW(state[sdl.SCANCODE_W] != 0)
			kbs.SetA(state[sdl.SCANCODE_A] != 0)
			kbs.SetS(state[sdl.SCANCODE_S] != 0)
			kbs.SetD(state[sdl.SCANCODE_D] != 0)
			kbs.SetSpace(state[sdl.SCANCODE_SPACE] != 0)
			kbs.SetLalt(state[sdl.SCANCODE_LALT] != 0)
			b := bytes.Buffer{}
			t.WriteTo(&b)
			nnInputPub.Send(append([]byte("input.kb:"), b.Bytes()...), 0)
			fmt.Printf("Keyboard input sent.\n")
				
		case state.MouseReset():
			var q ogre.Quaternion;
			is.orientation = q.FromValues(state.Quaternion().W(), state.Quaternion().X(),
				state.Quaternion().Y(), state.Quaternion().Z())
			var r ogre.Matrix3
			is.orientation.ToRotationMatrix(&r)
			var rfYAngle, rfPAngle, rfRAngle float32
			r.ToEulerAnglesYXZ(&rfYAngle, &rfPAngle, &rfRAngle)
			is.yaw = rad2Deg(rfYAngle)
			is.pitch = rad2Deg(rfPAngle)
			is.roll = rad2Deg(rfRAngle)
		case state.ConfigLookAround():
			if state.LookAround().ManipulateObject() {
				fmt.Printf("Input configuration: manipulate object\n");
				is.orientationFactor = 1.0;
			} else {
				fmt.Printf("Input configuration: look around\n");
				is.orientationFactor = -1.0
			}
		}
	}
	if !shutdownRequested {
      sendShutdown(nnRenderSocket, nnGameSocket)
      shutdownRequested = true
    }
    // make the GL context again before proceeding with the teardown
    if runtime.GOOS != "darwin" {
		sdl.GL_MakeCurrent(window, glContext)
	}
    waitShutdown(nnInputPull, &gameThreadParams)
}

func deg2Rad(deg float32) float32 {
	return deg * math.Pi / 180
}

func rad2Deg (rad float32) float32 {
	return rad * 180 / math.Pi
}

func sendShutdown(nnRenderSocket *nanomsg.Socket, nnGameSocket *nanomsg.Socket) {
	s := capn.NewBuffer(nil)
	stop := NewRootStop(s)
	stop.SetStop(true)
	buf := bytes.Buffer{}
	s.WriteTo(&buf)
	fmt.Printf("Render socket shutdown.\n")
	nnRenderSocket.Send(buf.Bytes(), 0)
	fmt.Printf("Game socket shutdown.\n")
	nnGameSocket.Send(buf.Bytes(), 0)
}

func waitShutdown(nnInputPull *nanomsg.Socket, params *GameThreadParams) {
	// For now, loop the input thread for a bit to flush out any events
	continueTime := time.Since(params.start) + 500 * time.Millisecond  // An eternity.
	for time.Since(params.start) < continueTime {	
		msg, _ := nnInputPull.Recv(nanomsg.DontWait)
		if msg == nil {
			sdl.Delay(10)
		}
	}
}

func getSubsystemString(info sdl.SysWMInfo) string {
	switch info.Subsystem {
	case 0:	
	    return "Unknown"
	case 1:
		return "Windows"
	case 2:
		return "X11"
	case 3:
		return "DirectFB"
	case 4: 
		return "Cocoa"
	case 5:
		return "UiKit"
	}
	return "Unknown"
}
