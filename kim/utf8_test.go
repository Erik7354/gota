package kim

import (
	"testing"
	"unicode/utf8"
)

func BenchmarkUTF8Encode(b *testing.B) {
	in := 'A'
	out := make([]byte, 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = utf8.EncodeRune(out, in)
	}
}

func BenchmarkUTF8Encode_100(b *testing.B) {
	in := []rune("Lórêm ípsüm dôlôr sît ämet, cönsectetür ädîpïsîcïng ëlît, sêd dô ëïüsmod têmpör.")
	out := make([]byte, 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, r := range in {
			_ = utf8.EncodeRune(out, r)
		}
	}
}

func BenchmarkUTF8Encode_250(b *testing.B) {
	in := []rune("Lórêm ípsüm dôlôr sît ämet, cönsectetür ädîpïsîcïng ëlît, sêd dô ëïüsmod têmpör. Üt ënäim âd mïním vëniâm, qüïs nöstrüd ëxërcïtätïön üllämçö lårëa 🐱, nïsi üt äliqüïp êx êä cömmo🚀 cönsëqüät. Düïs äü🌟dô löre💡 dölör ëü fügïåt nullä pârïätür.")
	out := make([]byte, 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, r := range in {
			_ = utf8.EncodeRune(out, r)
		}
	}
}

func BenchmarkUTF8Encode_250_ASCII(b *testing.B) {
	in := []rune("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit.")
	out := make([]byte, 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, r := range in {
			_ = utf8.EncodeRune(out, r)
		}
	}
}

func BenchmarkUTF8Decode(b *testing.B) {
	in := make([]byte, 4)
	utf8.EncodeRune(in, '😂')
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = utf8.DecodeRune(in)
	}
}
