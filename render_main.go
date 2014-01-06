package core

import ("fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/fire/go-ogre3d")

type RenderThreadParams struct {
	root ogre.Root
	window *sdl.Window
	ogreWindow ogre.RenderWindow
}

func renderThread(_parms RenderThreadParams) {
	fmt.Printf("Render Thread:\n")
}
