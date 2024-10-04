package nota

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

// string

func TestEncoder_encodeString(t *testing.T) {
	var res bytes.Buffer

	d := "abc😂"
	err := NewEncoder(&res).Encode(d)
	if err != nil {
		t.Fatal(err.Error())
	}

	ex := []byte{
		0b0_001_0100,
		97,
		98,
		99,
		0b1000_0111,
		0b1110_1100,
		0b0000_0010,
	}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}
}

// array, slice

func TestEncoder_encodeArray(t *testing.T) {
	var res bytes.Buffer

	d := []uint{1, 2, 3}
	err := NewEncoder(&res).Encode(d)
	if err != nil {
		t.Fatal(err.Error())
	}

	ex := []byte{
		0b0_010_0011,
		0b0_110_0_001,
		0b0_110_0_010,
		0b0_110_0_011,
	}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}
}

// map

func TestEncoder_encodeMap(t *testing.T) {
	var res bytes.Buffer

	err := NewEncoder(&res).Encode(map[string]bool{
		"a": true,
		"b": false,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	ex := []byte{
		0x32,       // record preamble with two entries
		0x11, 0x61, // "a"
		0x71,       // true
		0x11, 0x62, // "b"
		0x70, // false
	}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}
}

// int

func TestEncoder_encodeInt(t *testing.T) {
	var res bytes.Buffer

	// 1
	err := NewEncoder(&res).Encode(1)
	if err != nil {
		t.Fatal(err.Error())
	}

	ex := []byte{0b0_110_0_001}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}

	// 10
	res.Reset()

	err = NewEncoder(&res).Encode(10)
	if err != nil {
		t.Fatal(err.Error())
	}

	ex = []byte{0b1_110_0_000, 0b0000_1010}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}

	// 12345
	res.Reset()

	err = NewEncoder(&res).Encode(12345)
	if err != nil {
		t.Fatal(err.Error())
	}

	ex = []byte{0b1_110_0_000, 0b1_1100000, 0b0_0111001}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}

	// -12345
	res.Reset()

	err = NewEncoder(&res).Encode(-12345)
	if err != nil {
		t.Fatal(err.Error())
	}

	ex = []byte{0b1_110_1_000, 0b1_1100000, 0b0_0111001}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}
}

func BenchmarkEncoder_encodeInt(b *testing.B) {
	var res bytes.Buffer

	for i := 0; i < b.N; i++ {
		_ = NewEncoder(&res).Encode(math.MaxInt64)
	}
}

// uint

func TestEncoder_encodeUInt(t *testing.T) {
	var res bytes.Buffer

	// 0
	err := NewEncoder(&res).Encode(uint(0))
	if err != nil {
		t.Fatal(err.Error())
	}

	ex := []byte{0b0_110_0_000}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}

	// 7
	res.Reset()
	err = NewEncoder(&res).Encode(uint(7))
	if err != nil {
		t.Fatal(err.Error())
	}

	ex = []byte{0b0_110_0_111}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}

	// 8
	res.Reset()
	err = NewEncoder(&res).Encode(uint(8))
	if err != nil {
		t.Fatal(err.Error())
	}

	ex = []byte{0b1_110_0_000, 0b0_00001000}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}

	// 10
	res.Reset()
	err = NewEncoder(&res).Encode(uint(10))
	if err != nil {
		t.Fatal(err.Error())
	}

	ex = []byte{0b1_110_0_000, 0b0_00001010}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}

	// 18446744073709551615
	// 1 1111111 1111111 1111111 1111111 1111111 1111111 1111111 1111111 1111111
	res.Reset()
	err = NewEncoder(&res).Encode(uint(math.MaxUint64))
	if err != nil {
		t.Fatal(err.Error())
	}

	ex = []byte{0b1_110_0_001, 0b1_1111111, 0b1_1111111, 0b1_1111111, 0b1_1111111, 0b1_1111111, 0b1_1111111, 0b1_1111111, 0b1_1111111, 0b0_1111111}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}
}

func BenchmarkEncoder_encodeUInt(b *testing.B) {
	var res bytes.Buffer

	for i := 0; i < b.N; i++ {
		_ = NewEncoder(&res).Encode(uint(math.MaxUint64))
	}
}

// bool

func TestEncoder_encodeBool(t *testing.T) {
	var res bytes.Buffer

	// false
	err := NewEncoder(&res).Encode(false)
	if err != nil {
		t.Fatal(err.Error())
	}

	ex := []byte{0b0_111_0000}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}

	// true
	res.Reset()
	err = NewEncoder(&res).Encode(true)
	if err != nil {
		t.Fatal(err.Error())
	}

	ex = []byte{0b0_111_0001}
	if !reflect.DeepEqual(res.Bytes(), ex) {
		t.Fatalf("expected: %v \t got: %v", ex, res.Bytes())
	}
}

func BenchmarkEncoder_encodeBool(b *testing.B) {
	var res bytes.Buffer

	for i := 0; i < b.N; i++ {
		_ = NewEncoder(&res).Encode(false)
	}
}
