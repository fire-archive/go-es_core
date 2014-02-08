package core

import (
	"fmt"
	"bytes"
	"github.com/fire/go-ogre3d"
	"github.com/jmckaskill/go-capnproto"
)

type RenderState struct {
	mouseControl bool
	headNode     ogre.SceneNode
	rotationVectorObj ogre.ManualObject
	rotationVectorNode ogre.SceneNode
}

func parseRenderState(rs *RenderState, srs *SharedRenderState, b *[]byte) {	
	var buf bytes.Buffer 
	r := bytes.NewReader(*b)
	s, err := capn.ReadFromPackedStream(r, &buf)
	if err != nil {
		fmt.Printf("Read error %v\n", err)
		return
	}
        state := ReadRootRenderStateMsg(s)
	if state.HeadTrigger() {
		if state.FreeSpin() == false {
			rs.mouseControl = false
			rs.rotationVectorNode.SetVisible(false)		
		}
		if state.FreeSpin() {
			rs.mouseControl = true
			// Note: Uncomment to visualize the rotation vector
			rs.rotationVectorNode.SetVisible(true)
		}
		// Resume updating render state on next loop
		return
	}	  
	
	renderState := ReadRootEmittedRenderState(s)
	srs.gameTime = renderState.Time()
	srs.position.SetX(renderState.Position().X())
	srs.position.SetY(renderState.Position().Y())
	srs.orientation = ogre.CreateQuaternionFromValues(renderState.Orientation().W(), renderState.Orientation().X(), renderState.Orientation().Y(), renderState.Orientation().Z())
	srs.position.SetZ(0.0)
}

func renderInit(params *RenderThreadParams, rs *RenderState, srs *SharedRenderState) {
	fmt.Printf("Render Init:\n")
	rs.mouseControl = false
	
	mgr := ogre.GetResourceGroupManager()
	
	mgr.AddResourceLocation("media/models", "FileSystem", "General");
	mgr.AddResourceLocation("media/materials/scripts", "FileSystem", "General")
	mgr.AddResourceLocation("media/materials/textures", "FileSystem", "General")
	mgr.AddResourceLocation("media/materials/programs", "FileSystem", "General")

	mgr.InitialiseAllResourceGroups()

	scene := params.root.CreateSceneManager("DefaultSceneManager", "SimpleStaticCheck")
	scene.SetAmbientLight(0.5, 0.5, 0.5)
	head := scene.CreateEntity("head", "ogre.mesh", "head_group")
	rootNode := scene.GetRootSceneNode()
	zero := ogre.CreateVector3()
	zero.Zero()
	rs.headNode = rootNode.CreateChildSceneNode("head_node", zero , ogre.CreateQuaternion())
	rs.headNode.AttachObject(ogre.GetEntityBase(head))
	light := scene.CreateLight("light")
	light.SetPosition(20.0, 80.0, 50.0)
	cam := scene.CreateCamera("cam")
	cam.SetPosition(0,0,90)
	cam.LookAt(0,0,-300)
	cam.SetNearClipDistance(5)
	
	viewport := params.ogreWindow.AddViewport(cam)
	viewport.SetBackgroundColour(0, 0, 0, 0)
	
	cam.SetAspectRatio(viewport.GetActualWidth(), viewport.GetActualHeight())
	
	rs.rotationVectorObj = scene.CreateManualObject("rotation_vector")
	rs.rotationVectorObj.SetDynamic(true)
	rs.rotationVectorObj.Begin("BaseWhiteNoLighting", ogre.OT_LINE_LIST, "head")
	rs.rotationVectorObj.Position(0.0, 0.0, 0.0)
	rs.rotationVectorObj.Position(0.0, 0.0, 0.0)
	rs.rotationVectorObj.Position(0.0, 0.0, 0.0)
	rs.rotationVectorObj.End()
	rot := scene.GetRootSceneNode()
	rs.rotationVectorNode = rot.CreateChildSceneNode("rotation_vector_node", zero, ogre.CreateQuaternion())
	rs.rotationVectorNode.AttachObject(ogre.GetManualObjectBase(rs.rotationVectorObj))
	rs.rotationVectorNode.SetVisible(false)
}

func interpolateAndRender(rsockets *RenderThreadSockets, rs *RenderState,
	ratio float32, previousRender *SharedRenderState, nextRender *SharedRenderState) {
	temp := previousRender.position.MultiplyScalar(1.0 - ratio)
	interpPosition := temp.AddVector3(nextRender.position.MultiplyScalar(ratio))
	rs.headNode.SetPosition(interpPosition.X(), interpPosition.Y(), interpPosition.Z())
	if rs.mouseControl {
		t := capn.NewBuffer(nil)
		inputMouse := NewRootState(t)
		inputMouse.SetMouse(true)
		buf := bytes.Buffer{}
		t.WriteToPacked(&buf)
		rsockets.inputPush.Send(buf.Bytes(), 0)
		fmt.Printf("Render mouse_state requested\n")

		b, err := rsockets.inputMouseSub.Recv(0)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		b = bytes.TrimPrefix(b, []byte("input.mouse:"))
		var bBuf bytes.Buffer
		bBuf.Read(b)
		r := bytes.NewReader(b)
		fmt.Printf("Bytestring START%sEND\n", bBuf.String())
		var rBuf bytes.Buffer 
		s, err := capn.ReadFromPackedStream(r, &rBuf)
		if err != nil {
			fmt.Printf("Read error %v\n", err)
			return
		}		
		fmt.Printf("Render mouse_state received\n")
		input := ReadRootInputMouse(s)
		orientation := ogre.CreateQuaternionFromValues(input.W(), input.X(), input.Y(), input.Z())
		// Use latest mouse data to orient the head.
		rs.headNode.SetOrientation(orientation)
		// Update the rotation axis of the head (smoothed over a few frames in the game thread)
		rs.rotationVectorObj.BeginUpdate(0)
		rs.rotationVectorObj.Position(interpPosition.X(), interpPosition.Y(), interpPosition.Z())
		temp := previousRender.smoothedAngular.MultiplyScalar(1.0 - ratio)
		interpSmoothedAngular := temp.AddVector3(nextRender.smoothedAngular.MultiplyScalar(ratio))
		rotationVectorEnd := interpPosition.AddVector3(interpSmoothedAngular.MultiplyScalar(40.0))
		rs.rotationVectorObj.Position(rotationVectorEnd.X(), rotationVectorEnd.Y(), rotationVectorEnd.Z())
		rs.rotationVectorObj.End()
	} else {
		rs.headNode.SetOrientation(ogre.QuaternionSlerp(ratio, previousRender.orientation, nextRender.orientation, false))
	}	}

