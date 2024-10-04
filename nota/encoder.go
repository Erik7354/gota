package nota

import (
	"github.com/erik7354/gota/kim"
	"io"
	"math"
	"reflect"
)

type notaType = byte

const (
	notaBlob          notaType = 0b0000_0000
	notaText          notaType = 0b0001_0000
	notaArray         notaType = 0b0010_0000
	notaRecord        notaType = 0b0011_0000
	notaFloatingPoint notaType = 0b0100_0000
	notaInteger       notaType = 0b0110_0000
	notaSymbol        notaType = 0b0111_0000
)

// cm Continuation Mask
const cm = 0b0111_1111

// cb Continuation Bit
const cb = 0b1000_0000

// An Encoder writes NOTA values to an output stream.
type Encoder struct {
	w   io.Writer
	err error
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the NOTA encoding of v to the stream.
//
// See the documentation for [Marshal] for details about the
// conversion of Go values to NOTA.
func (enc Encoder) Encode(v any) error {
	return enc.encodeVal(v)
}

func (enc Encoder) encodeVal(v any) error {
	var err error
	switch reflect.TypeOf(v).Kind() {
	case reflect.String: // nota: text
		err = enc.encodeString(reflect.ValueOf(v))
	case reflect.Array, reflect.Slice: // nota: array
		err = enc.encodeArray(reflect.ValueOf(v))
	case reflect.Map: // nota: record
		if reflect.TypeOf(v).Key().Kind() != reflect.String {
			panic("unsupported key type of map, must be string")
		}
		err = enc.encodeMap(reflect.ValueOf(v))
	case reflect.Float32, reflect.Float64: // nota: floating point
		panic("unimplemented: go uses IEEE754 instead of DEC64 floats")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // nota: integer
		err = enc.encodeInt(reflect.ValueOf(v))
	case reflect.Bool: // nota: symbol
		err = enc.encodeBool(reflect.ValueOf(v))
	default:
		panic("unsupported type")
	}

	return err
}

func (enc Encoder) encodeInt(v reflect.Value) (err error) {
	var ui uint64
	var preamble byte = 0b0_110_0_000

	if v.CanInt() {
		i := v.Int()
		if i < 0 {
			preamble += 0b0_000_1_000
			i *= -1
		}
		ui = uint64(i)
	} else {
		ui = v.Uint()
	}

	switch {
	case ui < 1<<3:
		_, err = enc.w.Write([]byte{
			preamble + uint8(ui),
		})
	case ui < 1<<10:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>7),
			uint8(ui) & cm,
		})
	case ui < 1<<17:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>14),
			cb + uint8(ui>>7)&cm,
			uint8(ui) & cm,
		})
	case ui < 1<<24:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>21),
			cb + uint8(ui>>14)&cm,
			cb + uint8(ui>>7)&cm,
			uint8(ui) & cm,
		})
	case ui < 1<<31:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>28),
			cb + uint8(ui>>21)&cm,
			cb + uint8(ui>>14)&cm,
			cb + uint8(ui>>7)&cm,
			uint8(ui) & cm,
		})
	case ui < 1<<38:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>35),
			cb + uint8(ui>>28)&cm,
			cb + uint8(ui>>21)&cm,
			cb + uint8(ui>>14)&cm,
			cb + uint8(ui>>7)&cm,
			uint8(ui) & cm,
		})
	case ui < 1<<45:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>42),
			cb + uint8(ui>>35)&cm,
			cb + uint8(ui>>28)&cm,
			cb + uint8(ui>>21)&cm,
			cb + uint8(ui>>14)&cm,
			cb + uint8(ui>>7)&cm,
			uint8(ui) & cm,
		})
	case ui < 1<<52:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>49),
			cb + uint8(ui>>42)&cm,
			cb + uint8(ui>>35)&cm,
			cb + uint8(ui>>28)&cm,
			cb + uint8(ui>>21)&cm,
			cb + uint8(ui>>14)&cm,
			cb + uint8(ui>>7)&cm,
			uint8(ui) & cm,
		})
	case ui < 1<<59:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>56),
			cb + uint8(ui>>49)&cm,
			cb + uint8(ui>>42)&cm,
			cb + uint8(ui>>35)&cm,
			cb + uint8(ui>>28)&cm,
			cb + uint8(ui>>21)&cm,
			cb + uint8(ui>>14)&cm,
			cb + uint8(ui>>7)&cm,
			uint8(ui) & cm,
		})
	case ui <= math.MaxUint64:
		_, err = enc.w.Write([]byte{
			preamble + cb + uint8(ui>>63),
			cb + uint8(ui>>56)&cm,
			cb + uint8(ui>>49)&cm,
			cb + uint8(ui>>42)&cm,
			cb + uint8(ui>>35)&cm,
			cb + uint8(ui>>28)&cm,
			cb + uint8(ui>>21)&cm,
			cb + uint8(ui>>14)&cm,
			cb + uint8(ui>>7)&cm,
			uint8(ui) & cm,
		})
	}
	if err != nil {
		return err
	}

	return nil
}

func (enc Encoder) encodeArray(v reflect.Value) (err error) {
	if v.IsNil() {
		return nil
	}

	c := v.Len()
	if err = enc.writePreamble(notaArray, c); err != nil {
		return err
	}

	for i := range v.Len() {
		if err = enc.encodeVal(v.Index(i).Interface()); err != nil {
			return err
		}
	}

	return nil
}

func (enc Encoder) encodeMap(v reflect.Value) (err error) {
	if v.IsNil() {
		return nil
	}

	c := v.Len()
	if err = enc.writePreamble(notaRecord, c); err != nil {
		return err
	}

	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value().Interface()

		if err = enc.encodeString(key); err != nil {
			return err
		}
		if err = enc.encodeVal(val); err != nil {
			return err
		}
	}

	return nil
}

func (enc Encoder) encodeString(v reflect.Value) (err error) {
	s := v.String()

	c := len([]rune(s))
	if err = enc.writePreamble(notaText, c); err != nil {
		return err
	}

	b := make([]byte, 3)
	for _, r := range []rune(s) {
		n := kim.EncodeRune(b, r)

		_, err = enc.w.Write(b[0:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (enc Encoder) encodeBool(v reflect.Value) error {
	var b byte = 0b0_111_0000
	if v.Bool() {
		b = 0b0_111_0001
	}

	_, err := enc.w.Write([]byte{b})
	if err != nil {
		return err
	}

	return nil
}

// writePreamble writes preamble for text,
func (enc Encoder) writePreamble(t notaType, c int) error {
	var preamble []byte
	switch {
	case c < 1<<4:
		preamble = []byte{
			t + uint8(c),
		}
	case c < 1<<11:
		preamble = []byte{
			t + cb + uint8(c>>7),
			uint8(c) & cm,
		}
	case c < 1<<18:
		preamble = []byte{
			t + cb + uint8(c>>14),
			cb + uint8(c>>7)&cm,
			uint8(c) & cm,
		}
	case c < 1<<25:
		preamble = []byte{
			t + cb + uint8(c>>21),
			cb + uint8(c>>14)&cm,
			cb + uint8(c>>7)&cm,
			uint8(c) & cm,
		}
	case c < 1<<32:
		preamble = []byte{
			t + cb + uint8(c>>28),
			cb + uint8(c>>21)&cm,
			cb + uint8(c>>14)&cm,
			cb + uint8(c>>7)&cm,
			uint8(c) & cm,
		}
	case c < 1<<39:
		preamble = []byte{
			t + cb + uint8(c>>35),
			cb + uint8(c>>28)&cm,
			cb + uint8(c>>21)&cm,
			cb + uint8(c>>14)&cm,
			cb + uint8(c>>7)&cm,
			uint8(c) & cm,
		}
	case c < 1<<46:
		preamble = []byte{
			t + cb + uint8(c>>42),
			cb + uint8(c>>35)&cm,
			cb + uint8(c>>28)&cm,
			cb + uint8(c>>21)&cm,
			cb + uint8(c>>14)&cm,
			cb + uint8(c>>7)&cm,
			uint8(c) & cm,
		}
	case c < 1<<53:
		preamble = []byte{
			t + cb + uint8(c>>49),
			cb + uint8(c>>42)&cm,
			cb + uint8(c>>35)&cm,
			cb + uint8(c>>28)&cm,
			cb + uint8(c>>21)&cm,
			cb + uint8(c>>14)&cm,
			cb + uint8(c>>7)&cm,
			uint8(c) & cm,
		}
	case c < 1<<60:
		preamble = []byte{
			t + cb + uint8(c>>56),
			cb + uint8(c>>49)&cm,
			cb + uint8(c>>42)&cm,
			cb + uint8(c>>35)&cm,
			cb + uint8(c>>28)&cm,
			cb + uint8(c>>21)&cm,
			cb + uint8(c>>14)&cm,
			cb + uint8(c>>7)&cm,
			uint8(c) & cm,
		}
	case c <= math.MaxInt64:
		preamble = []byte{
			t + cb + uint8(c>>63),
			cb + uint8(c>>56)&cm,
			cb + uint8(c>>49)&cm,
			cb + uint8(c>>42)&cm,
			cb + uint8(c>>35)&cm,
			cb + uint8(c>>28)&cm,
			cb + uint8(c>>21)&cm,
			cb + uint8(c>>14)&cm,
			cb + uint8(c>>7)&cm,
			uint8(c) & cm,
		}
	}

	_, err := enc.w.Write(preamble)
	if err != nil {
		return err
	}

	return nil
}
