package core

import ("fmt"
	"time"
	"crypto/rand"
	"math"
	"math/big"
	"bytes"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/fire/go-ogre3d"
	"github.com/jmckaskill/go-capnproto"
	"github.com/op/go-nanomsg")

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
	rotationSpeed Degree  	// Rotation speed of the head in degrees:
	
	// use the last few frames of mouse input to build a smoothed angular velocity
	orientationIndex int
	orientationHistory[ORIENTATIONLOG] OrientationHistory
	smoothedAngular ogre.Vector3
	smoothedAngularVelocity Degree // Degree
}

func gameInit(gsockets *GameThreadSockets, gs *GameState, rs *SharedRenderState){
	fmt.Printf("Game Init.\n")
	gs.speed = randFloat32(59) + 40
	fmt.Printf("Random speed: %f\n", gs.speed)
	gs.bounce = 25.0
	dAngle := CreateDegree(randFloat32(359))
	gs.direction = ogre.CreateVector2FromValues(float32(math.Cos(float64(dAngle.ValueRadianFloat()))), float32(math.Sin(float64(dAngle.ValueRadianFloat()))))
	unitZ := ogre.CreateVector3()
	unitZ.UnitZ()
	rs.orientation = ogre.CreateQuaternion()
	rs.orientation.FromAngleAxis(0.0, unitZ)
	rs.position = ogre.CreateVector3()
	rs.position.Zero()
	gs.mousePressed = false
	gs.rotation = ogre.CreateVector3()
	gs.rotation.UnitX()
	gs.rotationSpeed = CreateDegree(0.0)
	gs.orientationIndex = 0
	fmt.Printf("Random angle: %f\n", dAngle.ValueDegreesFloat())
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
	b = bytes.TrimPrefix(b, []byte("input.mouse:"))
	// var by bytes.Buffer
	// by.Write(b)
	// fmt.Printf("Bytestring START%sEND\n", by.String())
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
	var tempDeg float32
	omega.ToAngleAxisDegree(&tempDeg, gs.smoothedAngular)
	gs.smoothedAngularVelocity = CreateDegree(tempDeg)
	fmt.Printf("%f %f %f - %f\n", gs.smoothedAngular.X(), gs.smoothedAngular.Y(), gs.smoothedAngular.Z(), gs.smoothedAngularVelocity.ValueDegreesFloat())
	srs.smoothedAngular = gs.smoothedAngular
	

	if (buttons & sdl.Button(sdl.BUTTON_LEFT)) != 0 {
		if !gs.mousePressed {
			gs.mousePressed = true
			// changing the control scheme: the player is now driving the orientation of the head directly with the mouse
			// tell the input logic to reset the orientation to match the current orientation of the head
			s := capn.NewBuffer(nil)
			state := NewRootState(s)
			state.SetMouseReset(true)
			state.Orientation().SetW(srs.orientation.W())
			state.Orientation().SetX(srs.orientation.X())
			state.Orientation().SetY(srs.orientation.Y())
			state.Orientation().SetZ(srs.orientation.Z())			
			buf := bytes.Buffer{}
			s.WriteTo(&buf)
			gsockets.inputPush.Send(buf.Bytes(), 0)	
			
			// IF RENDER TICK HAPPENS HERE: render will not know that it should grab the orientation directly from the mouse,
			// but the orientation coming from game should still be ok?

			s = capn.NewBuffer(nil)
			renderState := NewRootControlScheme(s)
			renderState.SetFreeSpin(true)			
			buf = bytes.Buffer{}
			s.WriteTo(&buf)
			gsockets.renderSocket.Send(buf.Bytes(), 0)                    
			// IF RENDER TICK HAPPENS HERE (before a new gamestate):
			// the now reset input orientation will combine with the old game state, that's bad
		}
	} else {
		if gs.mousePressed {
			gs.mousePressed = false
			// Changing the control scheme: the head will free spin and slow down for a bit, then it will resume bouncing around
			// the player looses mouse control, the game grabs latest orientation and angular velocity
			// the input thread was authoritative on orientation until now, so accept that as our starting orientation
			srs.orientation = orientation
			gs.rotationSpeed = gs.smoothedAngularVelocity
			gs.rotation = gs.smoothedAngular

			s := capn.NewBuffer(nil)
			renderState := NewRootControlScheme(s)
			renderState.SetFreeSpin(false)			
			buf := bytes.Buffer{}
			s.WriteTo(&buf)
			gsockets.renderSocket.Send(buf.Bytes(), 0)
			// IF RENDER TICK HAPPENS HERE (before a new gamestate): render will pull the head orientation from the game state rather than input, but game state won't have the fixed orientation yet
		}
	}

	if srs.position.X() > gs.bounce || srs.position.X() < -gs.bounce {
		gs.direction.SetX(gs.direction.X() * -1.0)
	}
	if srs.position.Y() > gs.bounce || srs.position.Y() < -gs.bounce {
		gs.direction.SetY(gs.direction.Y() * -1.0)
	}
	delta := gs.direction.MultiplyScalar(gs.speed * float32(float64(GAMEDELAY)/ float64(time.Second)))
	if !gs.mousePressed {
		if gs.rotationSpeed.ValueDegreesFloat() == 0.9 {
			srs.position.SetX(srs.position.X() + delta.X())
			srs.position.SetY(srs.position.Y() + delta.Y())
		}
		//    fmt.Printf("game tick position: %f %f\n", rs.position.X(), rs.position.Y());

		// update the orientation of the head on a free roll
		// gs.rotation is unit length
		// gs.rotationSpeed is in degrees/seconds
		// NOTE: sinf/cosf really needed there?
		gs.rotationSpeed.Mul(0.97)
		if gs.rotationSpeed.ValueDegreesFloat() < 20.0 {
			gs.rotationSpeed = CreateDegree(0.0)
		}
		tempDegree := CreateDegree(gs.rotationSpeed.ValueDegreesFloat() * float32(GAMETICKFLOAT))
		var factor float32 = float32(math.Sin(float64(0.5 * tempDegree.ValueRadianFloat())))
		rotationTick := ogre.CreateQuaternionFromValues(
			float32(math.Cos(float64(0.5 * tempDegree.ValueRadianFloat()))),
			factor * gs.rotation.X(),
			factor * gs.rotation.Y(),
			factor * gs.rotation.Z())
		rotationTick.Normalise()
		srs.orientation = rotationTick.MultiplyQuaternion(srs.orientation)
	} else {
		 // Keep updating the orientation in the render state, even while the render thread is ignoring it:
		// when the game thread resumes control of the head orientation, it will interpolate from one of these states,
		// so we keep updating the orientation to avoid a short glitch at the discontinuity
		srs.orientation = orientation
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

func emitRenderState(socket *nanomsg.BusSocket, time uint64, srs *SharedRenderState) {
        s := capn.NewBuffer(nil)
	emitted := NewRootEmittedRenderState(s)
	emitted.SetTime(time)
	emitted.Position().SetX(srs.position.X())
	emitted.Position().SetY(srs.position.Y())
	emitted.Orientation().SetW(srs.orientation.W())
	emitted.Orientation().SetX(srs.orientation.X())
	emitted.Orientation().SetY(srs.orientation.Y())
	emitted.Orientation().SetZ(srs.orientation.Z())
	emitted.SmoothedAngular().SetX(srs.smoothedAngular.X())
	emitted.SmoothedAngular().SetY(srs.smoothedAngular.Y())
	emitted.SmoothedAngular().SetZ(srs.smoothedAngular.Z())
	buf := bytes.Buffer{}
	s.WriteTo(&buf)
	socket.Send(buf.Bytes(), 0)
}
