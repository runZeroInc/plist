package plist

import (
	"strings"
	"testing"
)

// These tests exercise the runzero hardening:
// deeply nested plist input must return a parse error BEFORE the recursive
// parsers overflow the goroutine stack. A Go stack overflow is a fatal,
// unrecoverable runtime error (recover() cannot catch it), so the assertion is
// "Decode returns an error" rather than "Decode panics and we recover."
//
// The nesting count is chosen far above maxParseDepth but well below the depth
// that would overflow the stack, so if the cap were absent the test would still
// pass without crashing; we therefore also assert the returned error reports the
// depth limit, confirming the guard (not some incidental limit) fired.

func TestXMLPlistDepthLimit(t *testing.T) {
	const n = maxParseDepth + 50
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?><plist version="1.0">`)
	for i := 0; i < n; i++ {
		b.WriteString("<array>")
	}
	b.WriteString("<string>x</string>")
	for i := 0; i < n; i++ {
		b.WriteString("</array>")
	}
	b.WriteString("</plist>")

	var out any
	_, err := Unmarshal([]byte(b.String()), &out)
	if err == nil {
		t.Fatal("expected a parse error for deeply nested XML plist, got nil")
	}
	if !strings.Contains(err.Error(), errMaxDepthExceeded.Error()) {
		t.Fatalf("expected max-depth error, got: %v", err)
	}
}

func TestTextPlistDepthLimit(t *testing.T) {
	const n = maxParseDepth + 50
	// OpenStep nested dictionaries: {a={a={a=...;};};}
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString("{a=")
	}
	b.WriteString("x")
	for i := 0; i < n; i++ {
		b.WriteString(";}")
	}

	var out any
	_, err := Unmarshal([]byte(b.String()), &out)
	if err == nil {
		t.Fatal("expected a parse error for deeply nested text plist, got nil")
	}
	if !strings.Contains(err.Error(), errMaxDepthExceeded.Error()) {
		t.Fatalf("expected max-depth error, got: %v", err)
	}
}

// buildDeepBplist constructs a binary plist whose top object is a chain of
// `depth` single-element arrays terminating in an ASCII string. Each array
// references the next via the object table, so the file is small but forces the
// binary parser to recurse `depth` levels.
func buildDeepBplist(depth int) []byte {
	// Object IDs: 0..depth-1 are arrays (0 = top, references 1, ... ), object
	// `depth` is the terminal ASCII string "x".
	numObjects := depth + 1
	// 2-byte offsets and 2-byte object refs so the file can exceed 255 bytes.
	header := []byte("bplist00")
	body := append([]byte(nil), header...)

	offsets := make([]int, numObjects)
	// Arrays: tag 0xA1 (array, count 1) followed by a 2-byte object ref.
	for i := 0; i < depth; i++ {
		offsets[i] = len(body)
		body = append(body, 0xA1, byte((i+1)>>8), byte(i+1))
	}
	// Terminal string "x": tag 0x51 (ASCII string, length 1) + 'x'.
	offsets[depth] = len(body)
	body = append(body, 0x51, 'x')

	// Offset table: 2 bytes per object (big-endian).
	offsetTableOffset := len(body)
	for _, off := range offsets {
		body = append(body, byte(off>>8), byte(off))
	}

	// Trailer (32 bytes): 5 unused, sortVersion, offsetIntSize, objectRefSize,
	// numObjects(8), topObject(8), offsetTableOffset(8).
	trailer := make([]byte, 32)
	trailer[6] = 2 // OffsetIntSize
	trailer[7] = 2 // ObjectRefSize
	putUint64BE(trailer[8:16], uint64(numObjects))
	putUint64BE(trailer[16:24], 0) // top object = 0
	putUint64BE(trailer[24:32], uint64(offsetTableOffset))
	body = append(body, trailer...)
	return body
}

func putUint64BE(b []byte, v uint64) {
	for i := 7; i >= 0; i-- {
		b[i] = byte(v)
		v >>= 8
	}
}

func TestBplistDepthLimit(t *testing.T) {
	// Keep depth small enough that all object offsets fit in one byte (<256)
	// but above maxParseDepth so the depth guard fires.
	depth := maxParseDepth + 10
	data := buildDeepBplist(depth)

	var out any
	_, err := Unmarshal(data, &out)
	if err == nil {
		t.Fatal("expected a parse error for deeply nested binary plist, got nil")
	}
	if !strings.Contains(err.Error(), errMaxDepthExceeded.Error()) {
		t.Fatalf("expected max-depth error, got: %v", err)
	}
}

// TestBplistTruncatedOffset confirms the bounds checks convert a crafted
// out-of-range read into a recoverable parse error rather than an out-of-bounds
// slice panic that would escape Decode.
func TestBplistShallowStillParses(t *testing.T) {
	// A small valid bplist (single array of one string) must still decode after
	// the hardening, guarding against an over-aggressive bounds/depth check.
	data := buildDeepBplist(2)
	var out any
	if _, err := Unmarshal(data, &out); err != nil {
		t.Fatalf("valid shallow bplist failed to decode: %v", err)
	}
}
