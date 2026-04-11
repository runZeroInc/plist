package plist

import "unsafe"

func zeroCopy8BitString(buf []byte, off, n int) string {
	if n == 0 {
		return ""
	}
	return unsafe.String(&buf[off], n)
}
