package plist

import (
	"reflect"
	"testing"
	"time"
)

func TestInvalidMarshal(t *testing.T) {
	tests := []struct {
		Name  string
		Thing any
	}{
		{"Function", func() {}},
		{"Nil", nil},
		{"Map with integer keys", map[int]string{1: "hi"}},
		{"Channel", make(chan int)},
	}

	for _, v := range tests {
		t.Run(v.Name, func(t *testing.T) {
			data, err := Marshal(v.Thing, OpenStepFormat)
			if err == nil {
				t.Fatalf("expected error; got plist data: %x", data)
			} else {
				t.Log(err)
			}
		})
	}
}

type Cat struct{}

func (c *Cat) MarshalPlist() (any, error) {
	return "cat", nil
}

func TestInterfaceMarshal(t *testing.T) {
	var c Cat
	b, err := Marshal(&c, XMLFormat)
	if err != nil {
		t.Log(err)
	} else if len(b) == 0 {
		t.Log("expect non-zero data")
	}
}

func TestInterfaceFieldMarshal(t *testing.T) {
	type X struct {
		C any // C's type does not implement Marshaler
	}
	x := &X{
		C: &Cat{}, // C's value implements Marshaler
	}

	b, err := Marshal(x, XMLFormat)
	if err != nil {
		t.Log(err)
	} else if len(b) == 0 {
		t.Log("expect non-zero data")
	}
}

func TestMarshalInterfaceFieldPtrTime(t *testing.T) {
	type X struct {
		C any // C's type is unknown
	}

	sentinelTime := time.Date(2013, 11, 27, 0, 34, 0, 0, time.UTC)
	x := &X{
		C: &sentinelTime,
	}

	e := &Encoder{}
	rval := reflect.ValueOf(x)
	cf := e.marshal(rval)

	if dict, ok := cf.(*cfDictionary); ok {
		if _, ok := dict.values[0].(cfDate); !ok {
			t.Error("inner value is not a cfDate")
		}
	} else {
		t.Error("failed to marshal toplevel dictionary (?)")
	}
}

type Dog struct {
	Name string
}

type Animal any

func TestInterfaceSliceMarshal(t *testing.T) {
	x := make([]Animal, 0)
	x = append(x, &Dog{Name: "dog"})

	b, err := Marshal(x, XMLFormat)
	if err != nil {
		t.Error(err)
	} else if len(b) == 0 {
		t.Error("expect non-zero data")
	}
}

func TestInterfaceGeneralSliceMarshal(t *testing.T) {
	x := make([]any, 0) // accept any type
	x = append(x, &Dog{Name: "dog"}, "a string", 1, true)

	b, err := Marshal(x, XMLFormat)
	if err != nil {
		t.Error(err)
	} else if len(b) == 0 {
		t.Error("expect non-zero data")
	}
}
