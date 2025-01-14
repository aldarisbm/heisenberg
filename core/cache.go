package core

import (
	"github.com/quantanotes/heisenberg/utils"
)

type Cache struct {
	name    string
	kv      map[string]Value
	mapping map[uint]string
	index   *Index
}

func NewCache(name string, dim uint, space utils.SpaceType) *Cache {
	return &Cache{
		name:    name,
		kv:      make(map[string]Value),
		mapping: make(map[uint]string),
		index:   NewIndex(indexConfig{Name: name, Count: 0, Dim: dim, Space: uint(space)}),
	}
}

func (c *Cache) Get(key string) (Entry, error) {
	if val, ok := c.kv[key]; ok {
		return Entry{Collection: c.name, Key: key, Value: val}, nil
	}
	return Entry{}, utils.InvalidKey(c.name, key)
}

func (c *Cache) Put(key string, vector []float32, meta map[string]any) error {
	mapping := c.index.Next()
	if err := c.index.Insert(mapping, vector); err != nil {
		return err
	}
	c.kv[key] = Value{mapping, vector, meta}
	c.mapping[mapping] = key
	return nil
}

func (c *Cache) Delete(key string) error {
	val, ok := c.kv[key]
	if !ok {
		return utils.InvalidKey(key, c.name)
	}
	mapping := val.Index
	c.index.Delete(mapping)
	delete(c.mapping, mapping)
	delete(c.kv, key)
	return nil
}

func (c *Cache) Search(query []float32, k uint) ([]Entry, error) {
	ids, err := c.index.Search(query, int(k))
	if err != nil {
		return nil, err
	}
	results := make([]Entry, len(ids))
	for i, ids := range ids {
		key := c.mapping[uint(ids)]
		value := c.kv[key]
		entry := Entry{Collection: c.name, Key: key, Value: value}
		results[i] = entry
	}
	return results, nil
}
