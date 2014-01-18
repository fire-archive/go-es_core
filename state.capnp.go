package core

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/jmckaskill/go-capnproto"
	"math"
	"unsafe"
)

type State C.Struct
type StateQuaternion State
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
func (s State) Quaternion() StateQuaternion          { return StateQuaternion(s) }
func (s StateQuaternion) W() float32                 { return math.Float32frombits(C.Struct(s).Get32(4)) }
func (s StateQuaternion) SetW(v float32)             { C.Struct(s).Set32(4, math.Float32bits(v)) }
func (s StateQuaternion) X() float32                 { return math.Float32frombits(C.Struct(s).Get32(8)) }
func (s StateQuaternion) SetX(v float32)             { C.Struct(s).Set32(8, math.Float32bits(v)) }
func (s StateQuaternion) Y() float32                 { return math.Float32frombits(C.Struct(s).Get32(12)) }
func (s StateQuaternion) SetY(v float32)             { C.Struct(s).Set32(12, math.Float32bits(v)) }
func (s StateQuaternion) Z() float32                 { return math.Float32frombits(C.Struct(s).Get32(16)) }
func (s StateQuaternion) SetZ(v float32)             { C.Struct(s).Set32(16, math.Float32bits(v)) }
func (s State) LookAround() StateLookAround          { return StateLookAround(s) }
func (s StateLookAround) ManipulateObject() bool     { return C.Struct(s).Get1(1) }
func (s StateLookAround) SetManipulateObject(v bool) { C.Struct(s).Set1(1, v) }
func (s StateLookAround) LookAround() bool           { return C.Struct(s).Get1(2) }
func (s StateLookAround) SetLookAround(v bool)       { C.Struct(s).Set1(2, v) }

type State_List C.PointerList

func NewStateList(s *C.Segment, sz int) State_List { return State_List(s.NewCompositeList(24, 0, sz)) }
func (s State_List) Len() int                      { return C.PointerList(s).Len() }
func (s State_List) At(i int) State                { return State(C.PointerList(s).At(i).ToStruct()) }
func (s State_List) ToArray() []State              { return *(*[]State)(unsafe.Pointer(C.PointerList(s).ToArray())) }

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
