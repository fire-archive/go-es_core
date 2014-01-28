package core

import "github.com/fire/go-ogre3d"

type SharedRenderState struct {
  gameTime uint64 // needs to be the first parameter
  // control the head position and orientation
  position ogre.Vector3
  orientation ogre.Quaternion
  // an extra vector to visualize the rotation
  //smoothedAngular ogre.Vector3 
}

