package core

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/jmckaskill/go-capnproto"
	"math"
	"unsafe"
)

type State C.Struct
type StateOrientation State
type StateLookAround State
type State_Which uint16

const (
	STATE_MOUSE            State_Which = 0
	STATE_KB                           = 1
	STATE_MOUSERESET                   = 2
	STATE_CONFIGLOOKAROUND             = 3
)

func NewState(s *C.Segment) State                    { return State(s.NewStruct(24, 0)) }
func NewRootState(s *C.Segment) State                { return State(s.NewRootStruct(24, 0)) }
func ReadRootState(s *C.Segment) State               { return State(s.Root(0).ToStruct()) }
func (s State) Which() State_Which                   { return State_Which(C.Struct(s).Get16(2)) }
func (s State) Mouse() bool                          { return C.Struct(s).Get1(0) }
func (s State) SetMouse(v bool)                      { C.Struct(s).Set16(2, 0); C.Struct(s).Set1(0, v) }
func (s State) Kb() bool                             { return C.Struct(s).Get1(0) }
func (s State) SetKb(v bool)                         { C.Struct(s).Set16(2, 1); C.Struct(s).Set1(0, v) }
func (s State) MouseReset() bool                     { return C.Struct(s).Get1(0) }
func (s State) SetMouseReset(v bool)                 { C.Struct(s).Set16(2, 2); C.Struct(s).Set1(0, v) }
func (s State) ConfigLookAround() bool               { return C.Struct(s).Get1(0) }
func (s State) SetConfigLookAround(v bool)           { C.Struct(s).Set16(2, 3); C.Struct(s).Set1(0, v) }
func (s State) Orientation() StateOrientation        { return StateOrientation(s) }
func (s StateOrientation) W() float32                { return math.Float32frombits(C.Struct(s).Get32(4)) }
func (s StateOrientation) SetW(v float32)            { C.Struct(s).Set32(4, math.Float32bits(v)) }
func (s StateOrientation) X() float32                { return math.Float32frombits(C.Struct(s).Get32(8)) }
func (s StateOrientation) SetX(v float32)            { C.Struct(s).Set32(8, math.Float32bits(v)) }
func (s StateOrientation) Y() float32                { return math.Float32frombits(C.Struct(s).Get32(12)) }
func (s StateOrientation) SetY(v float32)            { C.Struct(s).Set32(12, math.Float32bits(v)) }
func (s StateOrientation) Z() float32                { return math.Float32frombits(C.Struct(s).Get32(16)) }
func (s StateOrientation) SetZ(v float32)            { C.Struct(s).Set32(16, math.Float32bits(v)) }
func (s State) LookAround() StateLookAround          { return StateLookAround(s) }
func (s StateLookAround) ManipulateObject() bool     { return C.Struct(s).Get1(1) }
func (s StateLookAround) SetManipulateObject(v bool) { C.Struct(s).Set1(1, v) }

type State_List C.PointerList

func NewStateList(s *C.Segment, sz int) State_List { return State_List(s.NewCompositeList(24, 0, sz)) }
func (s State_List) Len() int                      { return C.PointerList(s).Len() }
func (s State_List) At(i int) State                { return State(C.PointerList(s).At(i).ToStruct()) }
func (s State_List) ToArray() []State              { return *(*[]State)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type RenderStateMsg C.Struct

func NewRenderStateMsg(s *C.Segment) RenderStateMsg      { return RenderStateMsg(s.NewStruct(8, 0)) }
func NewRootRenderStateMsg(s *C.Segment) RenderStateMsg  { return RenderStateMsg(s.NewRootStruct(8, 0)) }
func ReadRootRenderStateMsg(s *C.Segment) RenderStateMsg { return RenderStateMsg(s.Root(0).ToStruct()) }
func (s RenderStateMsg) HeadTrigger() bool               { return C.Struct(s).Get1(1) }
func (s RenderStateMsg) SetHeadTrigger(v bool)           { C.Struct(s).Set1(1, v) }
func (s RenderStateMsg) FreeSpin() bool                  { return C.Struct(s).Get1(0) }
func (s RenderStateMsg) SetFreeSpin(v bool)              { C.Struct(s).Set1(0, v) }

type RenderStateMsg_List C.PointerList

func NewRenderStateMsgList(s *C.Segment, sz int) RenderStateMsg_List {
	return RenderStateMsg_List(s.NewUInt8List(sz))
}
func (s RenderStateMsg_List) Len() int { return C.PointerList(s).Len() }
func (s RenderStateMsg_List) At(i int) RenderStateMsg {
	return RenderStateMsg(C.PointerList(s).At(i).ToStruct())
}
func (s RenderStateMsg_List) ToArray() []RenderStateMsg {
	return *(*[]RenderStateMsg)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type EmittedRenderState C.Struct
type EmittedRenderStatePosition EmittedRenderState
type EmittedRenderStateOrientation EmittedRenderState
type EmittedRenderStateSmoothedAngular EmittedRenderState

func NewEmittedRenderState(s *C.Segment) EmittedRenderState {
	return EmittedRenderState(s.NewStruct(48, 0))
}
func NewRootEmittedRenderState(s *C.Segment) EmittedRenderState {
	return EmittedRenderState(s.NewRootStruct(48, 0))
}
func ReadRootEmittedRenderState(s *C.Segment) EmittedRenderState {
	return EmittedRenderState(s.Root(0).ToStruct())
}
func (s EmittedRenderState) Time() uint64     { return C.Struct(s).Get64(0) }
func (s EmittedRenderState) SetTime(v uint64) { C.Struct(s).Set64(0, v) }
func (s EmittedRenderState) Position() EmittedRenderStatePosition {
	return EmittedRenderStatePosition(s)
}
func (s EmittedRenderStatePosition) X() float32     { return math.Float32frombits(C.Struct(s).Get32(8)) }
func (s EmittedRenderStatePosition) SetX(v float32) { C.Struct(s).Set32(8, math.Float32bits(v)) }
func (s EmittedRenderStatePosition) Y() float32     { return math.Float32frombits(C.Struct(s).Get32(12)) }
func (s EmittedRenderStatePosition) SetY(v float32) { C.Struct(s).Set32(12, math.Float32bits(v)) }
func (s EmittedRenderState) Orientation() EmittedRenderStateOrientation {
	return EmittedRenderStateOrientation(s)
}
func (s EmittedRenderStateOrientation) W() float32     { return math.Float32frombits(C.Struct(s).Get32(16)) }
func (s EmittedRenderStateOrientation) SetW(v float32) { C.Struct(s).Set32(16, math.Float32bits(v)) }
func (s EmittedRenderStateOrientation) X() float32     { return math.Float32frombits(C.Struct(s).Get32(20)) }
func (s EmittedRenderStateOrientation) SetX(v float32) { C.Struct(s).Set32(20, math.Float32bits(v)) }
func (s EmittedRenderStateOrientation) Y() float32     { return math.Float32frombits(C.Struct(s).Get32(24)) }
func (s EmittedRenderStateOrientation) SetY(v float32) { C.Struct(s).Set32(24, math.Float32bits(v)) }
func (s EmittedRenderStateOrientation) Z() float32     { return math.Float32frombits(C.Struct(s).Get32(28)) }
func (s EmittedRenderStateOrientation) SetZ(v float32) { C.Struct(s).Set32(28, math.Float32bits(v)) }
func (s EmittedRenderState) SmoothedAngular() EmittedRenderStateSmoothedAngular {
	return EmittedRenderStateSmoothedAngular(s)
}
func (s EmittedRenderStateSmoothedAngular) X() float32 {
	return math.Float32frombits(C.Struct(s).Get32(32))
}
func (s EmittedRenderStateSmoothedAngular) SetX(v float32) { C.Struct(s).Set32(32, math.Float32bits(v)) }
func (s EmittedRenderStateSmoothedAngular) Y() float32 {
	return math.Float32frombits(C.Struct(s).Get32(36))
}
func (s EmittedRenderStateSmoothedAngular) SetY(v float32) { C.Struct(s).Set32(36, math.Float32bits(v)) }
func (s EmittedRenderStateSmoothedAngular) Z() float32 {
	return math.Float32frombits(C.Struct(s).Get32(40))
}
func (s EmittedRenderStateSmoothedAngular) SetZ(v float32) { C.Struct(s).Set32(40, math.Float32bits(v)) }

type EmittedRenderState_List C.PointerList

func NewEmittedRenderStateList(s *C.Segment, sz int) EmittedRenderState_List {
	return EmittedRenderState_List(s.NewCompositeList(48, 0, sz))
}
func (s EmittedRenderState_List) Len() int { return C.PointerList(s).Len() }
func (s EmittedRenderState_List) At(i int) EmittedRenderState {
	return EmittedRenderState(C.PointerList(s).At(i).ToStruct())
}
func (s EmittedRenderState_List) ToArray() []EmittedRenderState {
	return *(*[]EmittedRenderState)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type InputMouse C.Struct

func NewInputMouse(s *C.Segment) InputMouse      { return InputMouse(s.NewStruct(24, 0)) }
func NewRootInputMouse(s *C.Segment) InputMouse  { return InputMouse(s.NewRootStruct(24, 0)) }
func ReadRootInputMouse(s *C.Segment) InputMouse { return InputMouse(s.Root(0).ToStruct()) }
func (s InputMouse) W() float32                  { return math.Float32frombits(C.Struct(s).Get32(0)) }
func (s InputMouse) SetW(v float32)              { C.Struct(s).Set32(0, math.Float32bits(v)) }
func (s InputMouse) X() float32                  { return math.Float32frombits(C.Struct(s).Get32(4)) }
func (s InputMouse) SetX(v float32)              { C.Struct(s).Set32(4, math.Float32bits(v)) }
func (s InputMouse) Y() float32                  { return math.Float32frombits(C.Struct(s).Get32(8)) }
func (s InputMouse) SetY(v float32)              { C.Struct(s).Set32(8, math.Float32bits(v)) }
func (s InputMouse) Z() float32                  { return math.Float32frombits(C.Struct(s).Get32(12)) }
func (s InputMouse) SetZ(v float32)              { C.Struct(s).Set32(12, math.Float32bits(v)) }
func (s InputMouse) Buttons() uint32             { return C.Struct(s).Get32(16) }
func (s InputMouse) SetButtons(v uint32)         { C.Struct(s).Set32(16, v) }

type InputMouse_List C.PointerList

func NewInputMouseList(s *C.Segment, sz int) InputMouse_List {
	return InputMouse_List(s.NewCompositeList(24, 0, sz))
}
func (s InputMouse_List) Len() int            { return C.PointerList(s).Len() }
func (s InputMouse_List) At(i int) InputMouse { return InputMouse(C.PointerList(s).At(i).ToStruct()) }
func (s InputMouse_List) ToArray() []InputMouse {
	return *(*[]InputMouse)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type InputKb C.Struct

func NewInputKb(s *C.Segment) InputKb      { return InputKb(s.NewStruct(8, 0)) }
func NewRootInputKb(s *C.Segment) InputKb  { return InputKb(s.NewRootStruct(8, 0)) }
func ReadRootInputKb(s *C.Segment) InputKb { return InputKb(s.Root(0).ToStruct()) }
func (s InputKb) W() bool                  { return C.Struct(s).Get1(0) }
func (s InputKb) SetW(v bool)              { C.Struct(s).Set1(0, v) }
func (s InputKb) A() bool                  { return C.Struct(s).Get1(1) }
func (s InputKb) SetA(v bool)              { C.Struct(s).Set1(1, v) }
func (s InputKb) S() bool                  { return C.Struct(s).Get1(2) }
func (s InputKb) SetS(v bool)              { C.Struct(s).Set1(2, v) }
func (s InputKb) D() bool                  { return C.Struct(s).Get1(3) }
func (s InputKb) SetD(v bool)              { C.Struct(s).Set1(3, v) }
func (s InputKb) Space() bool              { return C.Struct(s).Get1(4) }
func (s InputKb) SetSpace(v bool)          { C.Struct(s).Set1(4, v) }
func (s InputKb) Lalt() bool               { return C.Struct(s).Get1(5) }
func (s InputKb) SetLalt(v bool)           { C.Struct(s).Set1(5, v) }

type InputKb_List C.PointerList

func NewInputKbList(s *C.Segment, sz int) InputKb_List { return InputKb_List(s.NewUInt8List(sz)) }
func (s InputKb_List) Len() int                        { return C.PointerList(s).Len() }
func (s InputKb_List) At(i int) InputKb                { return InputKb(C.PointerList(s).At(i).ToStruct()) }
func (s InputKb_List) ToArray() []InputKb {
	return *(*[]InputKb)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Stop C.Struct

func NewStop(s *C.Segment) Stop      { return Stop(s.NewStruct(8, 0)) }
func NewRootStop(s *C.Segment) Stop  { return Stop(s.NewRootStruct(8, 0)) }
func ReadRootStop(s *C.Segment) Stop { return Stop(s.Root(0).ToStruct()) }
func (s Stop) Stop() bool            { return C.Struct(s).Get1(0) }
func (s Stop) SetStop(v bool)        { C.Struct(s).Set1(0, v) }

type Stop_List C.PointerList

func NewStopList(s *C.Segment, sz int) Stop_List { return Stop_List(s.NewBitList(sz)) }
func (s Stop_List) Len() int                     { return C.PointerList(s).Len() }
func (s Stop_List) At(i int) Stop                { return Stop(C.PointerList(s).At(i).ToStruct()) }
func (s Stop_List) ToArray() []Stop              { return *(*[]Stop)(unsafe.Pointer(C.PointerList(s).ToArray())) }
