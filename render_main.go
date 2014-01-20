package core

import ("fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/fire/go-ogre3d"
	"github.com/op/go-nanomsg")

type RenderThreadParams struct {
	root ogre.Root
	window *sdl.Window
	ogreWindow ogre.RenderWindow
}

type RenderThreadSockets struct {
	inputPush *nanomsg.PushSocket
	inputMouseSub *nanomsg.SubSocket
}

func renderThread(params RenderThreadParams) {
	fmt.Printf("Render Thread:\n")
	var rsockets RenderThreadSockets
	
	controlSocket, err := nanomsg.NewBusSocket()
	if err != nil {
		panic(err)
    }
	_, err = controlSocket.Connect("tcp://127.0.0.1:60207") // Control render
    if err != nil {
		panic(err)
    }
    
    gameSocket, err := nanomsg.NewBusSocket()
    if err != nil {
		panic(err)
	}
	_, err = gameSocket.Connect("tcp://127.0.0.1:60210") // Game render
	if err != nil {
		panic(err)
    }
	// NOTE: since both render thread and game thread get spun at the same time,
    // and the connect needs to happen after the bind,
    // it's possible this would fail on occasion? just loop a few times and retry?
	
	rsockets.inputMouseSub, err = nanomsg.NewSubSocket()
	if err != nil {
		panic(err)
	}
	rsockets.inputMouseSub.Subscribe("input.mouse:")
	_, err = rsockets.inputMouseSub.Connect("tcp://127.0.0.1:60208")
	if err != nil {
		panic(err)
    }
    
    // Internal render state, not part of the interpolation:
    var rs RenderState
    
    // Always interpolate between two states.
    var srs[2] SharedRenderState
    srsIndex := uint(0)
    srs[1].gameTime = 0
    srs[0].gameTime = srs[1].gameTime

	renderInit(params, rs, srs[0])
}
