package core

import ("fmt"
		"time"
		"github.com/fire/go-ogre3d")

type OrientationHistory struct {
	t uint64
	o ogre.Quaternion
}

type GameState struct {
	bounce float32 			// Limits of the bounce area:
	speed float32 			// Picked a speed to bounce around at startup.
	mousePressed bool 		// Go from mouse is pressed to click each time to change the control scheme.
	//direction ogre.Vector2	// Direction the head is moving on the plane:
	//rotation ogre.Vector3	// Rotation axis of the head:
	rotationSpeed float32  	// Rotation speed of the head in degrees:
}

func gameInit(){
	fmt.Printf("Game Init.\n")
}

func gameTick(gs *GameState, srs *SharedRenderState, now time.Duration){
	fmt.Printf("Game Tick.\n")
}
