package core

import ("fmt"
		"github.com/fire/go-ogre3d")

type RenderState struct {
	mouseControl bool
	headNode *ogre.SceneNode
	//rotationVectorObj *ogre.ManualObject
	rotationVectorNode *ogre.SceneNode
}

func parseRenderState(rs *RenderState, srs *SharedRenderState, b *[]byte) {
	// Stuff
}

func renderInit(params *RenderThreadParams, rs *RenderState, srs *SharedRenderState ) {
	fmt.Printf("Render Init:\n")
}

func interpolateAndRender( rsockets *RenderThreadSockets, rs *RenderState,
	ratio float32, previousRender *SharedRenderState, nextRender *SharedRenderState) {
	// Stuff	
}
