# plist - A pure Go property list transcoder

A modern Go library for encoding and decoding Apple property list files.

This is a modernized fork of [howett.net/plist](https://pkg.go.dev/howett.net/plist), updated for idiomatic Go 1.20+ development. See [Changes from upstream](#changes-from-upstream) for details.

## Install

```
go get github.com/moond4rk/plist
```

## Features

- Encode and decode property lists from/to arbitrary Go types
- Supports all four formats: **XML**, **Binary**, **OpenStep**, and **GNUStep**
- Automatic format detection on decode
- Struct tags (`plist:"name,omitempty"`) for field customization
- Custom marshaler/unmarshaler interfaces
- Zero third-party dependencies

## Usage

### Encode a struct to XML plist

```go
type SparseBundleHeader struct {
    InfoDictionaryVersion string `plist:"CFBundleInfoDictionaryVersion"`
    BandSize              uint64 `plist:"band-size"`
    BackingStoreVersion   int    `plist:"bundle-backingstore-version"`
    DiskImageBundleType   string `plist:"diskimage-bundle-type"`
    Size                  uint64 `plist:"size"`
}

data := &SparseBundleHeader{
    InfoDictionaryVersion: "6.0",
    BandSize:              8388608,
    Size:                  4 * 1048576 * 1024 * 1024,
    DiskImageBundleType:   "com.apple.diskimage.sparsebundle",
    BackingStoreVersion:   1,
}

// Encode to XML with indentation
out, err := plist.MarshalIndent(data, plist.XMLFormat, "\t")
```

### Decode a plist file

```go
f, _ := os.Open("Info.plist")
defer f.Close()

var config map[string]any
decoder := plist.NewDecoder(f)
err := decoder.Decode(&config)
// decoder.Format contains the detected format (XMLFormat, BinaryFormat, etc.)
```

### Marshal to different formats

```go
value := map[string]any{
    "name":    "example",
    "count":   42,
    "enabled": true,
    "tags":    []string{"go", "plist"},
}

xmlBytes, _    := plist.Marshal(value, plist.XMLFormat)
binaryBytes, _ := plist.Marshal(value, plist.BinaryFormat)
```

### Unmarshal from bytes

```go
xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>name</key>
    <string>example</string>
    <key>count</key>
    <integer>42</integer>
</dict>
</plist>`)

var result map[string]any
format, err := plist.Unmarshal(xmlData, &result)
fmt.Println(format == plist.XMLFormat) // true
```

### Stream encoding

```go
encoder := plist.NewEncoderForFormat(os.Stdout, plist.XMLFormat)
encoder.Indent("\t")
err := encoder.Encode(value)
```

## Supported formats

| Constant | Format | Description |
|----------|--------|-------------|
| `XMLFormat` | Apple XML | Standard `.plist` files used by macOS/iOS |
| `BinaryFormat` | Apple Binary | Compact binary representation (`bplist00`) |
| `OpenStepFormat` | OpenStep | Legacy NeXTSTEP text format |
| `GNUStepFormat` | GNUStep | Extended ASCII format with typed values |

## Struct tags

```go
type Example struct {
    Name    string `plist:"display_name"`         // custom key name
    Count   int    `plist:"count,omitempty"`       // omit if zero value
    Ignored string `plist:"-"`                     // skip this field
}
```

## Changes from upstream

This fork modernizes [howett.net/plist](https://pkg.go.dev/howett.net/plist) for current Go development practices:

- **Go 1.20 minimum** - updated from Go 1.12
- **Zero dependencies** - removed `go-flags` and `yaml.v3` (CLI tools removed)
- **Modern Go idioms** - `any` instead of `interface{}`, `//go:build` tags, `unsafe.String` instead of deprecated `reflect.StringHeader`
- **Deprecated APIs removed** - replaced `io/ioutil` with `io`/`os` equivalents
- **Performance improvements** - `strings.Builder` for text format encoding
- **Bug fixes** - empty `<plist/>` now correctly returns an error; stronger test assertions
- **CI/Lint** - golangci-lint v2 with comprehensive linter configuration
- **Dead code removed** - Go 1.6/1.7 compatibility shims, AppEngine build variants

## Requirements

- Go 1.20+
