package kim

const MaxRune = 0x10FFFF // unicode.MaxRune

const RuneError = 0xFFFD

// cb continuation bit set
const cb = 0b1000_0000

const (
	rune1Max = 0x7F
	rune2Max = 0x3FFF
)

// EncodeRune writes into p (which must be large enough) the KIM encoding of the rune.
// It returns the number of bytes written.
func EncodeRune(p []byte, r rune) int {
	switch i := uint32(r); {
	case i <= rune1Max:
		p[0] = byte(i)
		return 1
	case i <= rune2Max:
		p[0] = byte(cb | i>>7)
		p[1] = byte(i)
		return 2
	case i > MaxRune:
		i = RuneError
		fallthrough
	default: // <= MaxRune
		p[0] = byte(cb | i>>14)
		p[1] = byte(cb | i>>7)
		p[2] = byte(i)
		return 3
	}
}

// DecodeRune unpacks the first KIM encoding in p and returns the rune and
// its width in bytes. If p is empty it returns (RuneError, 0). Otherwise, if
// the encoding is invalid, it returns (RuneError, 1). Both are impossible
// results for correct, non-empty KIM.
func DecodeRune(p []byte) (r rune, size int) {
	if len(p) < 1 {
		return RuneError, 0
	}

	r += int32(0b0111_1111 & p[0])
	if p[0] < cb {
		return r, 1
	} else if p[0] == cb {
		return RuneError, 1
	}

	r = r << 7
	r += int32(0b0111_1111 & p[1])
	if p[1] < cb {
		return r, 2
	}

	r = r << 7
	r += int32(0b0111_1111 & p[2])

	if r > MaxRune {
		return RuneError, 1
	}

	return r, 3
}
