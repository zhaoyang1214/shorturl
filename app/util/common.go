package util

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
