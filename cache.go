//
// Primitive Cach feature
//
// Last used last deleted [string]interface{}.
//
// Copyright 2020 Aterier UEDA
// Author: Takeyuki UEDA

package cache2

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type CachedValueType struct {
	Mu    sync.Mutex
	Value interface{}
	cache *Cache
}

type fifoElmType struct {
	lastUpdated int64
	id          interface{}
}

type Cache struct {
	maxSize   int
	valueMap  map[interface{}]*CachedValueType
	fifoArray []fifoElmType
	debug     bool
}

var mu sync.Mutex

func NewCache(maxSize int, debug bool) (*Cache, error) {
	cache := Cache{} // initialize
	cache.maxSize = maxSize
	cache.valueMap = map[interface{}]*CachedValueType{}
	cache.debug = debug
	return &cache, nil
}

// AddOrReplace
func (cache *Cache) AddOrReplace(id interface{}, value interface{}) (result *CachedValueType) { // Add & Replace

	mu.Lock()
	defer mu.Unlock()

	_, isExist := cache.valueMap[id]
	if isExist {
		// remove from fifo (then add to bottom later)
		cache.removefromFifo(id)
	} else if len(cache.valueMap) >= cache.maxSize {
		// delete oldest
		delete(cache.valueMap, cache.fifoArray[0].id)
		cache.fifoArray = cache.fifoArray[1:]
	}
	// add (or replace) new one
	result = makeValue(value, cache)
	cache.valueMap[id] = result
	cache.fifoArray = append(cache.fifoArray, makeFifoElm(id))

	if cache.debug {
		cache.DumpValueMap()
		cache.DumpFifoArray()
	}

	return
}

// Get
func (cache *Cache) Get(id interface{}) (c *CachedValueType, isExist bool) {

	mu.Lock()
	defer mu.Unlock()

	c, isExist = cache.valueMap[id]
	if isExist {
		cache.toBottom(id)
	}

	if cache.debug {
		cache.DumpValueMap()
		cache.DumpFifoArray()
	}

	return
}

// Delete
func (cache *Cache) Delete(id interface{}) {

	mu.Lock()
	defer mu.Unlock()

	// remove from CacheTable
	delete(cache.valueMap, id)
	// remove from CacheOrder
	cache.removefromFifo(id)

	if cache.debug {
		cache.DumpValueMap()
		cache.DumpFifoArray()
	}

	return
}

// move to Bottom
func (cache *Cache) MoveToBottom(id interface{}) (err error) {

	mu.Lock()
	defer mu.Unlock()

	_, isExist := cache.valueMap[id]
	if isExist {
		cache.toBottom(id)
	} else {
		err = fmt.Errorf("err: ID=%v is not exist on this cache.", id)
	}
	return
}

//
// internal functions
//
// Prerequisite:
//   The key is confirmed to exist in the cache
//

func (cache *Cache) removefromFifo(id interface{}) {
	for i, fifo := range cache.fifoArray {
		if fifo.id == id {
			cache.fifoArray = append(cache.fifoArray[:i], cache.fifoArray[i+1:]...)
			break
		}
	}
}

func (cache *Cache) toBottom(id interface{}) {
	// first of all, remove from fifo
	cache.removefromFifo(id)
	// then, add to bottom again
	cache.fifoArray = append(cache.fifoArray, makeFifoElm(id))
}

//
// Helper functions
//

// make body
func makeValue(value interface{}, cache *Cache) (v *CachedValueType) {
	v = new(CachedValueType)
	v.Value = value
	v.cache = cache
	return
}

// make fifoElm
func makeFifoElm(key interface{}) fifoElmType {
	return fifoElmType{id: key, lastUpdated: time.Now().Unix()}
}

// Dump Keys
func (cache *Cache) DumpKeys() {
	log.Println("*** Dump Cache Keys ***")
	for key, _ := range cache.valueMap {
		log.Println(key)
	}
	log.Println("***********************")
}

// Dump valueMap
func (cache *Cache) DumpValueMap() {
	log.Println("len(cache.valueMap)", len(cache.valueMap))
	log.Println("cache.valueMap", cache.MakeKeyValueMap())
}

// Dump fifoArray
func (cache *Cache) DumpFifoArray() {
	log.Println("len(cache.fifoArray)", len(cache.fifoArray))
	log.Println("cache.fifoArray", cache.fifoArray)
}

func (cache *Cache) MakeKeyValueMap() (result map[interface{}]interface{}) {
	result = map[interface{}]interface{}{}
	for key, value := range cache.valueMap {
		result[key] = value.Value
	}
	return
}
