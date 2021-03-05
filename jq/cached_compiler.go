package jq

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/itchyny/gojq"
)

// CachedCompiler wraps a compiler with an LRU cache.
type CachedCompiler struct {
	compiler Compiler
	cache    *lru.Cache
}

// NewCachedCompiler returns a new LRU-caching compiler of given cache size.
func NewCachedCompiler(compiler Compiler, size int) (*CachedCompiler, error) {
	lruCache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &CachedCompiler{
		compiler: compiler,
		cache:    lruCache,
	}, nil
}

// Compiler looks up the rawQuery in the cache. If found, it returns the
// precompiled query. Otherwise, it compiles rawQuery, and caches the result.
func (c *CachedCompiler) Compiler(rawQuery string) (*gojq.Code, error) {
	code, found := c.cache.Get(rawQuery)
	if found {
		return code.(*gojq.Code), nil
	}
	code, err := c.compiler(rawQuery)
	if err != nil {
		return nil, err
	}
	c.cache.Add(rawQuery, code)
	return code.(*gojq.Code), nil
}
