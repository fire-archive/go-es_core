package core

import (
	"math"
)

type Degree struct {
	degree float32
}

func CreateDegree(deg float32) Degree {
	var result Degree
	result.degree = deg
	return result
}

func (deg *Degree) Mul(real float32) {
	deg.degree = deg.degree * real
}

func (deg *Degree) Add(real float32) {
	deg.degree = deg.degree + real
}

func (deg *Degree) ValueDegreesFloat() float32 {
	return deg.degree
}

func (deg *Degree) ValueDegrees() Degree {
	return *deg
}

func (deg *Degree) ValueRadianFloat() float32 {
	return deg.degree * math.Pi / 180
}

func (deg *Degree) ValueRadian() Radian {
	var result Radian
	result.radian = deg.ValueRadianFloat()
	return result
}

type Radian struct {
	radian float32
}

func CreateRadian(rad float32) Radian {
	var result Radian
	result.radian = rad
	return result
}

func (rad *Radian) ValueDegreesFloat() float32 {
	return rad.radian * 180 / math.Pi
}

func (rad *Radian) ValueDegrees() Degree {
	var result Degree
	result.degree = rad.ValueDegreesFloat()
	return result
}

func (rad *Radian) ValueRadian() Radian {
	return *rad
}

func (rad *Radian) ValueRadianFloat() float32 {
	return rad.radian
}
