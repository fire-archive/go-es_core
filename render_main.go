package core

import ("fmt"
		"time"
		"github.com/jackyb/go-sdl2/sdl"
		"github.com/fire/go-ogre3d"
		"github.com/op/go-nanomsg"
		"github.com/jmckaskill/go-capnproto")

type RenderThreadParams struct {
	root ogre.Root
	window *sdl.Window
	ogreWindow ogre.RenderWindow
	start time.Time
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

	renderInit(&params, &rs, &srs[0])
	
	for true {
		b, err := controlSocket.Recv(nanomsg.DontWait)
		if err != nil {
			fmt.Printf("%s\n", err)
		}	
		if b != nil {
			s, _, err := capn.ReadFromMemoryZeroCopy(b)
			if err != nil {
				fmt.Printf("Read error %v\n", err)
			}
			stop := ReadRootStop(s)	
			if stop.Stop() {
				break // exit the thread
			}
		}
		for true {
			// Any message from the game thread?
			b, err := gameSocket.Recv(nanomsg.DontWait)
			if err != nil {
				fmt.Printf("%s\n", err)
			}	
			if b == nil {
				break
			}
			
			srsIndex ^= 1
			parseRenderState(&rs, &srs[srsIndex], &b)
		}
		// Skip rendering until enough data has come in to support interpolation
		if srs[0].gameTime == srs[1].gameTime { // 0 == 0
			continue
		}
		preRenderTime := time.Since(params.start)
		ratio := float32(uint64(preRenderTime) - srs[srsIndex ^ 1].gameTime) /
			float32(srs[srsIndex].gameTime - srs[srsIndex ^ 1].gameTime) 
		fmt.Printf("Render ratio $f\n", ratio)
		
		interpolateAndRender(&rsockets, &rs, ratio, &srs[srsIndex^1], &srs[srsIndex])
		
		//params.root._fireFrameStarted()
		params.root.RenderOneFrame()
		//params.root._fireFrameRenderingQueued()
		// 'render latency': How late is the image we displayed?
		// If vsync is off, it's the time it took to render the frame.
		// If vsync is on, it's render time + time waiting for the buffer swap.
		// NOTE: could display it terms of ratio also?
		postRenderTime := time.Since(params.start)
		fmt.Printf("Render latency %f ms.\n", float32(postRenderTime - preRenderTime) / float32(time.Millisecond))
	}
}
