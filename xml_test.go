package plist

import (
	"bytes"
	"fmt"
	"testing"
)

var InvalidXMLPlists = []struct {
	Name string
	Data string
}{
	{"hex integer with no digits", `<plist version="1.0"><integer>0x</integer></plist>`},
	{"unknown element doct", "<plist><doct><key>helo</key><string></string></doct></plist>"},
	{"dict with string instead of key", "<plist><dict><string>helo</string></dict></plist>"},
	{"dict with key but no value", "<plist><dict><key>helo</key></dict></plist>"},
	{"integer with non-numeric value", "<plist><integer>helo</integer></plist>"},
	{"empty integer", "<plist><integer></integer></plist>"},
	{"real with non-numeric value", "<plist><real>helo</real></plist>"},
	{"data with invalid base64", "<plist><data>*@&amp;%#helo</data></plist>"},
	{"date with invalid format", "<plist><date>*@&amp;%#helo</date></plist>"},
	{"unclosed integer tag", "<plist><integer>10</plist>"},
	{"unclosed real tag", "<plist><real>10</plist>"},
	{"unclosed string tag", "<plist><string>10</plist>"},
	{"unclosed dict tag", "<plist><dict>10</plist>"},
	{"unclosed key tag", "<plist><dict><key>10</plist>"},
	{"truncated plist open tag", "<plist>"},
	{"truncated data tag", "<plist><data>"},
	{"truncated date tag", "<plist><date>"},
	{"truncated array tag", "<plist><array>"},
	{"self-closing empty plist", "<plist/>"},
	{"truncated XML", "<pl"},
	{"binary plist magic as XML", "bplist00"},
}

func TestVariousIllegalXMLPlists(t *testing.T) {
	for i, tc := range InvalidXMLPlists {
		t.Run(fmt.Sprintf("%d_%s", i, tc.Name), func(t *testing.T) {
			buf := bytes.NewReader([]byte(tc.Data))
			d := newXMLPlistParser(buf)
			_, err := d.parseDocument()
			if err == nil {
				t.Errorf("expected error for invalid plist %q, got nil", tc.Data)
			} else {
				t.Logf("Error: %v", err)
			}
		})
	}
}
