package plethora

import (
	"io"
	"reflect"
	"strings"
)

// seperator is used as seperator in various string utilities
// Note: changing this will break existing databases
const seperator = ":"

var nameToProvider map[string]DataProvider
var providerToName map[string]string

// DataType is a kind of data
//
// Each DataType is backed by a DataProvider that stores data in a
// place the provider sees fit. It is possible for multiple providers
// to exist for a single DataType.
type DataType string

// DataProvider returns the Data associated with the identifier passed.
type DataProvider func(identifier string) (Data, error)

// Data is the interface used to support arbitrary kinds
// of data in plethora. All data types need to support the
// interface to be able to register with plethora.
type Data interface {
	// Type returns the type of data this is
	Type() DataType
	// Provider returns the DataProvider of this data
	Provider() DataProvider
	// Identifier is called to get an unique identifier to this
	// data. The identifier only has to be unique to the
	// DataProvider returned by Provider.
	Identifier() string
	// Render should write a html representation of the data to
	// the writer given.
	Render(w io.Writer) error
}

type identifier string

// Identifier returns an identifier unique to this Data, this is
// different from Data.Identifier in that the former is only unique
// to the provider associated with the Data. The identifier returned
// by Identifier is unique in the whole system.
func Identifier(d Data) identifier {
	p := d.Provider()
	if p == nil {
		panic("illegal: data returned nil provider")
	}

	return identifier(providerName(p) + seperator + d.Identifier())
}

// RegisterProvider registers a DataProvider with the given name.
//
// A DataProvider can be registered under multiple names
//
// Note: changing the name of a data provider after any data has
// entered the system will result in data becoming unreachable.
// Therefore if you want to change the name of your data provider
// you should register a backwards-compatible version under the old
// name to be used with the existing data.
func RegisterProvider(name string, provider DataProvider) {
	if nameToProvider[name] != nil {
		panic("illegal: double register for single name")
	}

	if strings.Index(name, seperator) > 0 {
		panic("illegal: seperator contained in provider name")
	}

	nameToProvider[name] = provider
	providerToName[providerName(provider)] = name
}

// providerName returns a name unique to this DataProvider
func providerName(p DataProvider) string {
	return reflect.TypeOf(p).String()
}
