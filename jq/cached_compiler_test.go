package jq

import (
	"github.com/itchyny/gojq"
	"testing"
)

type mockCompiler struct {
	Calls uint
}

func (c *mockCompiler) Compiler(rawQuery string) (*gojq.Code, error) {
	c.Calls++
	return QueryCompiler(rawQuery)
}

func TestCachedCompilerValidQuery(t *testing.T) {
	rawQuery := "."
	compiler := mockCompiler{}
	cachedCompiler, _ := NewCachedCompiler(compiler.Compiler, 1)
	code, err := cachedCompiler.Compiler(rawQuery)
	if err != nil {
		t.Fatal("Compilation failed")
	}
	if code == nil {
		t.Fatal("Compilation failed")
	}
	if compiler.Calls != 1 {
		t.Fatal("Compiler not called")
	}
}

func TestCachedCompilerValidQueryCacheHit(t *testing.T) {
	rawQuery := "."
	compiler := mockCompiler{}
	cachedCompiler, _ := NewCachedCompiler(compiler.Compiler, 1)
	code, err := cachedCompiler.Compiler(rawQuery)
	if err != nil {
		t.Fatal("Compilation failed")
	}
	if code == nil {
		t.Fatal("Compilation failed")
	}
	code, err = cachedCompiler.Compiler(rawQuery)
	if err != nil {
		t.Fatal("Compilation failed")
	}
	if code == nil {
		t.Fatal("Compilation failed")
	}
	if compiler.Calls != 1 {
		t.Fatal("Compiler called too often")
	}
}

func TestCachedCompilerValidQueryCacheMiss(t *testing.T) {
	rawQuery := "."
	compiler := mockCompiler{}
	cachedCompiler, _ := NewCachedCompiler(compiler.Compiler, 1)
	code, err := cachedCompiler.Compiler(rawQuery)
	if err != nil {
		t.Fatal("Compilation failed")
	}
	if code == nil {
		t.Fatal("Compilation failed")
	}
	code, err = cachedCompiler.Compiler(".[]")
	if err != nil {
		t.Fatal("Compilation failed")
	}
	if code == nil {
		t.Fatal("Compilation failed")
	}
	if compiler.Calls != 2 {
		t.Fatal("Compiler not called often enough")
	}
}

func TestCachedCompilerInvalidQuery(t *testing.T) {
	rawQuery := "!"
	compiler := mockCompiler{}
	cachedCompiler, _ := NewCachedCompiler(compiler.Compiler, 1)
	code, err := cachedCompiler.Compiler(rawQuery)
	if err == nil {
		t.Fatal("Compilation succeeded")
	}
	if code != nil {
		t.Fatal("Compilation succeeded")
	}
	code, err = cachedCompiler.Compiler(rawQuery)
	if err == nil {
		t.Fatal("Compilation succeeded")
	}
	if code != nil {
		t.Fatal("Compilation succeeded")
	}
	if compiler.Calls != 2 {
		t.Fatal("Cached invalid query")
	}
}
