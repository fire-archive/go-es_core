package core

import ("fmt"
		"time"
		"crypto/rand"
		"math/big"
		"bytes"
		"github.com/fire/go-ogre3d"
		"github.com/jmckaskill/go-capnproto")

const ORIENTATIONLOG int = 10;

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
	
	// use the last few frames of mouse input to build a smoothed angular velocity
	orientationIndex int
	orientationHistory[ORIENTATIONLOG] OrientationHistory
	//smoothedAngular ogre.Vector3
	smoothedAngularVelocity float32 // Degree
}

func gameInit(gsockets *GameThreadSockets, gs *GameState, srs *SharedRenderState){
	fmt.Printf("Game Init.\n")
	gs.speed = randFloat32(59) + 40
	fmt.Printf("Random speed: %f\n", gs.speed)
	gs.bounce = 25.0
	angle := deg2Rad(randFloat32(359))
	//gs.direction = ogre.Vector2(math.Cos(angle), math.Sin(angle))
	//rs.orientation.fromAngleAxis(0.0, ogre.Vector3.unitZ)
	//rs.position = ogre.Vector3.zero
	gs.mousePressed = false
	//gs.rotation = ogre.Vector3.unitX
	gs.rotationSpeed = 0.0
	gs.orientationIndex = 0
	fmt.Printf("Random angle: %f\n", angle)
	// Set the input code to manipulate an object rather than look around.
	s := capn.NewBuffer(nil)
	lookAround := NewRootState(s)
	lookAround.SetConfigLookAround(true)
	lookAround.LookAround().SetManipulateObject(true)
	buf := bytes.Buffer{}
	s.WriteTo(&buf)
	gsockets.inputPush.Send(buf.Bytes(), 0)
}

func gameTick(gs *GameState, srs *SharedRenderState, now time.Duration){
	fmt.Printf("Game Tick.\n")
}

// Create a random 32bit float from [1,max+1).
func randFloat32(max uint64) float32 {
	i := big.NewInt(0)
	r, err := rand.Int(rand.Reader, i.SetUint64(uint64(1) << 63))
	if err != nil {
			fmt.Printf("%s\n", err)
	}	
	return float32(float64(r.Uint64()) / float64(1 << 63) * float64(max))
}
