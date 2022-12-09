package flexbuffers

type BitWidth int

const (
	Width8 BitWidth = iota
	Width16
	Width32
	Width64
	WidthInvalid = -1
)

func toByteWidth(bitWidth BitWidth) int {
	return 1 << bitWidth
}

func paddingSize(bufSize, bitWidth int) int {
	return (^bufSize + 1) & (bitWidth - 1)
}

type ValueType int

const (
	Null ValueType = iota
	Int8
	Int16
	Int32
	Int64
	UInt
	Float
	Key
	String
	IndirectInt
	IndirectUInt
	IndirectFloat
	Map
	Vector
	VectorInt
	VectorUInt
	VectorFloat
	VectorKey
	VectorInt2
	VectorUInt2
	VectorFloat2
	VectorInt3
	VectorUInt3
	VectorFloat3
	VectorInt4
	VectorUInt4
	VectorFloat4
	Blob
	Bool
	VectorBool
)

func isInline(valuetype ValueType) bool {
	// TODO: Change Float once more float types are introduced.
	return valuetype == Bool || valuetype >= Int8 && valuetype < Float
}
