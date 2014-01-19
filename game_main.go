package core

import ("fmt"
		"time"
		"github.com/op/go-nanomsg")
const MAXFRAMERATE = 60 
const GAMEDELAY = time.Duration(time.Second / MAXFRAMERATE) 
const GAMETICKFLOAT = float64(GAMEDELAY) / float64(time.Millisecond)

type GameThreadSockets struct {
	controlSocket *nanomsg.Socket
	inputMouseSub *nanomsg.SubSocket
	inputKbSub *nanomsg.SubSocket
	inputPush *nanomsg.Socket
	renderSocket *nanomsg.Socket
}

func gameThread(params GameThreadParams) {
	var gsockets GameThreadSockets
	var gs GameState
	var srs SharedRenderState
	var err error
	gsockets.controlSocket, err = nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
	if err != nil {
		panic(err)
    }
    _, err = gsockets.controlSocket.Connect("tcp://127.0.0.1:60206")
    if err != nil {
		panic(err)
    }
	gsockets.renderSocket, err = nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.BUS)
	if err != nil {
		panic(err)
    }
    _, err = gsockets.renderSocket.Bind("tcp://127.0.0.1:60210")
    if err != nil {
		panic(err)
    }

	gsockets.inputMouseSub, err = nanomsg.NewSubSocket()
	if err != nil {
		panic(err)
	}
	gsockets.inputMouseSub.Subscribe("input.mouse:")
	_, err = gsockets.inputMouseSub.Connect("tcp://127.0.0.1:60208")
	if err != nil {
		panic(err)
	}
	
	gsockets.inputKbSub, err = nanomsg.NewSubSocket()
	gsockets.inputKbSub.Subscribe("input.kb:")
	if err != nil {
		panic(err)
	}
	_, err = gsockets.inputKbSub.Connect("tcp://127.0.0.1:60208")

	gameInit()
	baseLine := time.Since(params.start)
	var framenum uint64
	framenum = 0
	for true {
		now := time.Since(params.start)
		targetFrame := uint64 ((now - baseLine) / GAMEDELAY)
		if framenum <= targetFrame {
			framenum++
			// NOTE: build the state of the world at t = framenum * GAME_DELAY,
			// under normal conditions that's a time in the future
			// (the exception to that is if we are catching up on ticking game frames)
			gameTick(&gs, &srs, now);
			// Notify the render thread that a new game state is ready.
			// On the next render frame, it will start interpolating between the previous state and this new one
		} else {
			ahead := time.Duration(framenum) * GAMEDELAY - (now - baseLine)
			if ahead < 0 {
				panic(fmt.Sprintf("Ahead is less than 0: %d\n", ahead))
			}
			fmt.Printf("Game sleep %d ms\n", ahead)
			time.Sleep(ahead)
		}
		//cmd := 
	}
}
