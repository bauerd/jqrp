package jq

import "github.com/itchyny/gojq"

// Compiler compiles raw queries.
type Compiler = func(string) (*gojq.Code, error)

// QueryCompiler converts raw query strings into compiled queries.
func QueryCompiler(rawQuery string) (*gojq.Code, error) {
	query, err := gojq.Parse(rawQuery)
	if err != nil {
		return nil, err
	}
	code, err := gojq.Compile(query)
	if err != nil {
		return nil, err
	}
	return code, nil
}
