package core

import ("fmt"
	"time"
	"crypto/rand"
	"math"
	"math/big"
	"bytes"
	"github.com/jackyb/go-sdl2/sdl"
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
	direction ogre.Vector2	// Direction the head is moving on the plane:
	rotation ogre.Vector3	// Rotation axis of the head:
	rotationSpeed float32  	// Rotation speed of the head in degrees:
	
	// use the last few frames of mouse input to build a smoothed angular velocity
	orientationIndex int
	orientationHistory[ORIENTATIONLOG] OrientationHistory
	smoothedAngular ogre.Vector3
	smoothedAngularVelocity float32 // Degree
}

func gameInit(gsockets *GameThreadSockets, gs *GameState, rs *SharedRenderState){
	fmt.Printf("Game Init.\n")
	gs.speed = randFloat32(59) + 40
	fmt.Printf("Random speed: %f\n", gs.speed)
	gs.bounce = 25.0
	angle := deg2Rad(randFloat32(359))
	gs.direction = ogre.CreateVector2FromValues(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle))))
	unitZ := ogre.CreateVector3()
	unitZ.UnitZ()
	rs.orientation = ogre.CreateQuaternion()
	rs.orientation.FromAngleAxis(0.0, unitZ)
	rs.position = ogre.CreateVector3()
	rs.position.Zero()
	gs.mousePressed = false
	gs.rotation = ogre.CreateVector3()
	gs.rotation.UnitX()
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

func gameTick(gsockets *GameThreadSockets, gs *GameState, srs *SharedRenderState, now time.Duration){
	fmt.Printf("Game Tick.\n")

	// Get the latest mouse buttons state and orientation.
	s := capn.NewBuffer(nil)
	state := NewRootState(s)
	state.SetMouse(true)
	buf := bytes.Buffer{}
	s.WriteTo(&buf)
	gsockets.inputPush.Send(buf.Bytes(), 0)

	b, err := gsockets.inputMouseSub.Recv(0)
	if err != nil {
		fmt.Printf("%s\n", err)
	}	
	s, _, err = capn.ReadFromMemoryZeroCopy(b)
	if err != nil {
		fmt.Printf("Read error %v\n", err)
		return
	}	
	input := ReadRootInputMouse(s)
	orientation := ogre.CreateQuaternionFromValues(input.W(), input.X(), input.Y(), input.Z())
	buttons := input.Buttons()

	// At 16 ms tick and the last 10 orientations buffered, that's 150ms worth of orientation history.
	gs.orientationHistory[gs.orientationIndex].t = uint64(now)
	gs.orientationHistory[gs.orientationIndex].o = orientation
	gs.orientationIndex = (gs.orientationIndex + 1) % ORIENTATIONLOG
	
	// Oldest Orientation
	q1Index := gs.orientationIndex
	// NOTE: the problem with using the successive orientations to infer an angular speed,
	// is that if the orientation is changing fast enough, this code will 'flip' the speed around
	// e.g. this doesn't work, need to use the XY mouse data to track angular speed
	// NOTE: uncomment the following line to use the full history, notice the 'flip' happens at much lower speed
	q1Index = ( q1Index + ORIENTATIONLOG - 2) % ORIENTATIONLOG
	q1 := gs.orientationHistory[q1Index].o
	q1T := gs.orientationHistory[q1Index].t
	omega := orientation.SubtractQuaternion(q1)
	omega = omega.MultiplyScalar(2.0)
	omega = omega.UnitInverse()
	omega = omega.MultiplyScalar(float32(float64(time.Second)/float64(now - time.Duration(q1T))))
	omega.Normalise()
	omega.ToAngleAxisDegree(&gs.smoothedAngularVelocity, gs.smoothedAngular)
	//  fmt.Printf("%f %f %f - %f\n", gs.smoothed_angular.X(), gs.smoothed_angular.Y(), gs.smoothed_angular.Z(), gs.smoothed_angular_velocity.valueDegrees())
	srs.smoothedAngular = gs.smoothedAngular

	if (buttons & sdl.Button(sdl.BUTTON_LEFT)) != 0 {
		if !gs.mousePressed {
			gs.mousePressed = true
			// changing the control scheme: the player is now driving the orientation of the head directly with the mouse
			// tell the input logic to reset the orientation to match the current orientation of the head
			s := capn.NewBuffer(nil)
			state := NewRootState(s)
			state.SetMouseReset(true)
			state.Quaternion().SetW(srs.orientation.W())
			state.Quaternion().SetX(srs.orientation.X())
			state.Quaternion().SetY(srs.orientation.Y())
			state.Quaternion().SetZ(srs.orientation.Z())			
			buf := bytes.Buffer{}
			s.WriteTo(&buf)
			gsockets.inputPush.Send(buf.Bytes(), 0)	
		}
	}
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
