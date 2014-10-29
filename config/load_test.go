package config

import (
	"reflect"
	"testing"

	"github.com/Wessie/plethora/config/testutil"
)

var storeAndLoadTests = []struct {
	key   string
	value interface{}
}{
	{"int", 150},
	{"string", "hello world"},
	{"float", 1.5},
	{"bytes", []byte("testing")},
}

func TestConfigStoreAndLoad(t *testing.T) {
	db, cleanup, err := testutil.NewTestDB()
	defer cleanup()
	if err != nil {
		t.Fatal(err)
	}

	c := Open(db)
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