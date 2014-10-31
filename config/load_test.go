package config

import (
	"reflect"
	"testing"
)

var storeAndLoadTests = []struct {
	key   string
	value interface{}
}{
	{"int", 150},
	{"int", 200},
	{"int", 6000},
	{"string", "hello world"},
	{"string", "hello"},
	{"string", "world"},
	{"float", 1.5},
	{"bytes", []byte("testing")},
}

func TestConfigStoreAndLoad(t *testing.T) {
	defer TestConfiguration()()

	c := OpenConfig("test")
	for _, tCase := range storeAndLoadTests {
		err := c.Store(tCase.key, tCase.value)
		if err != nil {
			t.Error(err)
			continue
		}

		v := reflect.New(reflect.TypeOf(tCase.value))

		err = c.Load(tCase.key, v.Interface())
		if err != nil {
			t.Error(err)
			continue
		}

		nv := v.Elem().Interface()
		if !reflect.DeepEqual(nv, tCase.value) {
			t.Error("mismatched values")
		}
	}
}
