package kim

import (
	"testing"
)

func TestEncodeRune(t *testing.T) {
	buf := make([]byte, 3)

	n := EncodeRune(buf, 'c')
	if buf[0] != 99 || n != 1 {
		t.Fail()
	}

	n = EncodeRune(buf, 'a')
	if buf[0] != 97 || n != 1 {
		t.Fail()
	}

	n = EncodeRune(buf, 't')
	if buf[0] != 116 || n != 1 {
		t.Fail()
	}

	n = EncodeRune(buf, '😂')
	if buf[0] != 0b_1000_0111 || buf[1] != 0b_1110_1100 || buf[2] != 0b_0000_0010 || n != 3 {
		t.Fail()
	}
}

func BenchmarkEncode(b *testing.B) {
	in := 'A'
	out := make([]byte, 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = EncodeRune(out, in)
	}
}

func BenchmarkEncode_100(b *testing.B) {
	in := []rune("Lórêm ípsüm dôlôr sît ämet, cönsectetür ädîpïsîcïng ëlît, sêd dô ëïüsmod têmpör.")
	out := make([]byte, 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, r := range in {
			_ = EncodeRune(out, r)
		}
	}
}

func BenchmarkEncode_250(b *testing.B) {
	in := []rune("Lórêm ípsüm dôlôr sît ämet, cönsectetür ädîpïsîcïng ëlît, sêd dô ëïüsmod têmpör. Üt ënäim âd mïním vëniâm, qüïs nöstrüd ëxërcïtätïön üllämçö lårëa 🐱, nïsi üt äliqüïp êx êä cömmo🚀 cönsëqüät. Düïs äü🌟dô löre💡 dölör ëü fügïåt nullä pârïätür.")
	out := make([]byte, 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, r := range in {
			_ = EncodeRune(out, r)
		}
	}
}

func BenchmarkEncode_250_ASCII(b *testing.B) {
	in := []rune("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit.")
	out := make([]byte, 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, r := range in {
			_ = EncodeRune(out, r)
		}
	}
}

func TestDecodeRune(t *testing.T) {
	r, size := DecodeRune([]byte{})
	if r != RuneError || size != 0 {
		t.Fail()
	}

	r, size = DecodeRune([]byte{99, 100, 101})
	if r != 'c' || size != 1 {
		t.Fail()
	}

	r, size = DecodeRune([]byte{97})
	if r != 'a' || size != 1 {
		t.Fail()
	}

	r, size = DecodeRune([]byte{116})
	if r != 't' || size != 1 {
		t.Fail()
	}

	r, size = DecodeRune([]byte{0b_1000_0010, 0b_0001_0000})
	if r != 'Đ' || size != 2 {
		t.Fail()
	}

	r, size = DecodeRune([]byte{0b_1000_0111, 0b_1110_1100, 0b_0000_0010})
	if r != '😂' || size != 3 {
		t.Fail()
	}

	r, size = DecodeRune([]byte{0b_1100_0011, 0b_1111_1111, 0b_0111_1111})
	if r != MaxRune || size != 3 {
		t.Fail()
	}

	// >MaxRune
	r, size = DecodeRune([]byte{0b_1111_1111, 0b_1111_1111, 0b_0111_1111})
	if r != RuneError || size != 1 {
		t.Fail()
	}
}

func BenchmarkDecodeRune(b *testing.B) {
	in := []byte{0b_1000_0111, 0b_1110_1100, 0b_0000_0010} // 😂

	for i := 0; i < b.N; i++ {
		_, _ = DecodeRune(in)
	}
}
