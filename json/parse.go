// Package json exists because dynamic unmarshalling of JSON to either a Slice
// or Map is not easily afforded by encoding/json.
package json

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
)

var (
	// ErrPrimitiveRootType signals that a decoded data structure is primitive,
	// i.e. not an object or array, and therefore cannot be root.
	ErrPrimitiveRootType = errors.New("root type is primitive")

	// ErrMultipleRoots signals that multiple data structures were decoded.
	ErrMultipleRoots = errors.New("multiple roots")
)

// Parse returns parsed JSON. It decodes the first JSON data structure from
// reader. If there is more than one structure in reader, it errors. If the data
// structure is a primitive type, i.e. not an object or array, it errors.
func Parse(reader io.Reader) (interface{}, error) {
	var result interface{}
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	switch reflect.TypeOf(result).Kind() {
	case reflect.Slice:
	case reflect.Map:
		break
	default:
		return nil, ErrPrimitiveRootType
	}
	if decoder.More() {
		return nil, ErrMultipleRoots
	}
	return result, nil
}
