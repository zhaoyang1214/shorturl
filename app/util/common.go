package util

import (
	"errors"
	"math"
	"strings"
)

const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func FormatInt(u uint64, base int) string {
	if base < 2 || base > len(digits) {
		panic("FormatInt base error")
	}
	var a [64]byte
	i := len(a)
	b := uint64(base)
	for u >= b {
		i--
		a[i] = digits[uint(u%b)]
		u = u / b
	}
	i--
	a[i] = digits[uint(u)]

	return string(a[i:])
}

func ParseUint(s string, base int) (uint64, error) {
	if base < 2 || base > len(digits) {
		return 0, errors.New("base error")
	}
	j := len(s) - 1
	if j < 0 || j > 63 {
		return 0, errors.New("string invalid")
	}
	var a uint64
	for i := j; i >= 0; i-- {
		d := strings.IndexByte(digits, s[i])
		if d == -1 {
			return 0, errors.New("string invalid")
		}
		ud := uint64(d)

		b := uint64(math.Pow(float64(base), float64(j-i)))
		if math.MaxUint64/b < ud || math.MaxUint64-a < ud*b {
			return 0, errors.New("overflows uint64")
		}

		a += ud * b
	}
	return a, nil
}
