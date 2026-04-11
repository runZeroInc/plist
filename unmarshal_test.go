package plist

import (
	"reflect"
	"testing"
)

type CustomDate struct{}

func (cd *CustomDate) UnmarshalPlist(unmarshal func(any) error) error { return nil }

func TestCustomDateUnmarshal(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <date>2003-02-03T09:00:00.00Z</date>
</plist>`

	var custom CustomDate
	if _, err := Unmarshal([]byte(input), &custom); err != nil {
		t.Error(err)
	}
}

func TestInvalidMapKeyTypeUnmarshal(t *testing.T) {
	m := make(map[int]string)
	dict := &cfDictionary{
		keys: []string{"1", "2"},
		values: []cfValue{
			cfString("first"),
			cfString("second"),
		},
	}

	var caught bool
	{
		defer func() {
			if e := recover(); e != nil {
				t.Log("error:", e)
				caught = true
			}
		}()

		d := &Decoder{}
		d.unmarshalDictionary(dict, reflect.ValueOf(m))
	}

	if !caught {
		t.Fail()
	}
}

func TestValidButAliasedMapKeyTypeUnmarshal(t *testing.T) {
	type sortaString string
	m := make(map[sortaString]string)
	dict := &cfDictionary{
		keys: []string{"1", "2"},
		values: []cfValue{
			cfString("first"),
			cfString("second"),
		},
	}

	var caught bool
	{
		defer func() {
			if e := recover(); e != nil {
				t.Error("error:", e)
				caught = true
			}
		}()

		d := &Decoder{}
		d.unmarshalDictionary(dict, reflect.ValueOf(m))
	}

	if caught {
		t.Fail()
	}
}
