package jq

import (
	"testing"
)

func TestCompilerValidQuery(t *testing.T) {
	rawQuery := "."
	code, err := QueryCompiler(rawQuery)
	if err != nil {
		t.Fatal("Compilation failed")
	}
	if code == nil {
		t.Fatal("Compilation failed")
	}
}

func TestCompilerInvalidQuery(t *testing.T) {
	rawQuery := "!"
	code, err := QueryCompiler(rawQuery)
	if err == nil {
		t.Fatal("Compilation succeeded")
	}
	if code != nil {
		t.Fatal("Compilation succeeded")
	}
}
