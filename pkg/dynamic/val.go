package dynamic

type String func() string

func (s String) Val() string {
	return s()
}

func NewString(cb func() string) String {
	return cb
}

type Bool func() bool

func (s Bool) Val() bool {
	return s()
}

func NewBool(cb func() bool) Bool {
	return cb
}

type Int func() int

func (s Int) Val() int {
	return s()
}

func NewInt(cb func() int) Int {
	return cb
}

type Int8 func() int8

func (s Int8) Val() int8 {
	return s()
}

func NewInt8(cb func() int8) Int8 {
	return cb
}

type Int16 func() int16

func (s Int16) Val() int16 {
	return s()
}

func NewInt16(cb func() int16) Int16 {
	return cb
}

type Int32 func() int32

func (s Int32) Val() int32 {
	return s()
}

func NewInt32(cb func() int32) Int32 {
	return cb
}

type Int64 func() int64

func (s Int64) Val() int64 {
	return s()
}

func NewInt64(cb func() int64) Int64 {
	return cb
}

type UnInt func() uint

func (s UnInt) Val() uint {
	return s()
}

func NewUnInt(cb func() uint) UnInt {
	return cb
}

type UnInt8 func() uint8

func (s UnInt8) Val() uint8 {
	return s()
}

func NewUnInt8(cb func() uint8) UnInt8 {
	return cb
}

type UnInt16 func() uint16

func (s UnInt16) Val() uint16 {
	return s()
}

func NewUnInt16(cb func() uint16) UnInt16 {
	return cb
}

type UnInt32 func() uint32

func (s UnInt32) Val() uint32 {
	return s()
}

func NewUnInt32(cb func() uint32) UnInt32 {
	return cb
}

type UnInt64 func() uint64

func (s UnInt64) Val() uint64 {
	return s()
}

func NewUnInt64(cb func() uint64) UnInt64 {
	return cb
}

type Float32 func() float32

func (s Float32) Val() float32 {
	return s()
}

func NewFloat32(cb func() float32) Float32 {
	return cb
}

type Float64 func() float64

func (s Float64) Val() float64 {
	return s()
}

func NewFloat64(cb func() float64) Float64 {
	return cb
}
