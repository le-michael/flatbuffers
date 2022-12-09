package flexbuffers

import "fmt"

type BuilderOpts struct {
	size int
}

var DefaultBuilderOpts = BuilderOpts{
	size: 2048,
}

type Builder struct {
	Bytes         []byte
	finished      bool
	stack         []*stackValue
	stackPointers []*stackPointer
	offset        int
}

// NewBuilder creates a new flexbuffer builder with default options.
func NewBuilder() *Builder {
	return NewBuilderWithOpts(DefaultBuilderOpts)
}

// NewBuilderWithOpts creates a new flexbuffer with custom options.
func NewBuilderWithOpts(opts BuilderOpts) *Builder {
	return &Builder{
		Bytes:    make([]byte, opts.size),
		finished: false,
	}
}

func (b *Builder) AddNull() error {
	if err := b.integrityCheckOnValueAddition(); err != nil {
		return err
	}
	b.stack = append(b.stack, newStackValueWithNull())
	return nil
}

func (b *Builder) AddInt8(value int8) error {
	if err := b.integrityCheckOnValueAddition(); err != nil {
		return err
	}
	b.stack = append(b.stack, newStackValueWithInt8(value))
	return nil
}

func (b *Builder) AddInt16(value int16) error {
	if err := b.integrityCheckOnValueAddition(); err != nil {
		return err
	}
	b.stack = append(b.stack, newStackValueWithInt16(value))
	return nil
}

func (b *Builder) AddInt32(value int32) error {
	if err := b.integrityCheckOnValueAddition(); err != nil {
		return err
	}
	b.stack = append(b.stack, newStackValueWithInt32(value))
	return nil
}

func (b *Builder) AddInt64(value int64) error {
	if err := b.integrityCheckOnValueAddition(); err != nil {
		return err
	}
	b.stack = append(b.stack, newStackValueWithInt64(value))
	return nil
}

// Finish building the flexbuffer and return a slice of bytes.
func (b *Builder) Finish() ([]byte, error) {
	if b.finished {
		return b.Bytes, nil
	}

	if err := b.finish(); err != nil {
		return nil, fmt.Errorf("Unable to finish builder: %v", err)
	}

	return b.Bytes, nil
}

func (b *Builder) finish() error {
	if len(b.stack) != 1 {
		return fmt.Errorf("Stack length has to be exactly 1, but is %d.", len(b.stack))
	}

	value := b.stack[0]
	elementWidth, err := value.elementWidth(b.offset, 0)
	if err != nil {
		return fmt.Errorf("Unable to calculate element width: %v", err)
	}
	byteWidth := b.align(elementWidth)
	b.writeStackValue(value, byteWidth)

	b.finished = true
	return nil
}

func (b *Builder) writeStackValue(value *stackValue, byteWidth int) {

}

func (b *Builder) newOffset(newValueSize int) int {
	newOffset := b.offset + newValueSize
	size := len(b.Bytes)
	prevSize := size

	for size < newOffset {
		size <<= 1;
	}

	return newOffset
}

func (b *Builder) align(width BitWidth) int {
	byteWidth := toByteWidth(width)
	b.offset += paddingSize(b.offset, byteWidth)
	return byteWidth
}

func (b *Builder) integrityCheckOnValueAddition() error {
	if b.finished {
		return fmt.Errorf("Cannot add value to finished builder.")
	}

	if len(b.stackPointers) > 0 && !b.stackPointers[len(b.stackPointers)-1].isVector {
		if b.stack[len(b.stack)-1].valueType != Key {
			return fmt.Errorf("Cannot add value to a map before adding a key")
		}
	}

	return nil
}

func (b *Builder) integrityCheckOnKeyAddition() error {
	if b.finished {
		return fmt.Errorf("Cannot add value to finished builder.")
	}

	if len(b.stackPointers) == 0 || b.stackPointers[len(b.stackPointers)-1].isVector {
		return fmt.Errorf("Cannot add key before starting a map")
	}

	return nil
}

type stackValue struct {
	value     any
	offset    int
	valueType ValueType
	bitWidth  BitWidth
}

func newStackValueWithNull() *stackValue {
	return &stackValue{
		valueType: Null,
		bitWidth:  Width8,
	}
}

func newStackValueWithInt8(value int8) *stackValue {
	return &stackValue{
		value:     value,
		valueType: Int8,
		bitWidth:  Width8,
	}
}
func newStackValueWithInt16(value int16) *stackValue {
	return &stackValue{
		value:     value,
		valueType: Int16,
		bitWidth:  Width16,
	}
}
func newStackValueWithInt32(value int32) *stackValue {
	return &stackValue{
		value:     value,
		valueType: Int32,
		bitWidth:  Width32,
	}
}

func newStackValueWithInt64(value int64) *stackValue {
	return &stackValue{
		value:     value,
		valueType: Int64,
		bitWidth:  Width64,
	}
}

func (s *stackValue) elementWidth(size, index int) (BitWidth, error) {
	if isInline(s.valueType) {
		return s.bitWidth, nil
	}
	// TODO: Support non-inline value types
	return WidthInvalid, fmt.Errorf("Element is of unknown size.")
}

type stackPointer struct {
	stackPosition int
	isVector      bool
}

func newStatckPointer(stackPosition int, isVector bool) *stackPointer {
	return &stackPointer{stackPosition: stackPosition, isVector: isVector}
}
