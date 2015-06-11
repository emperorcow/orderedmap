package orderedmap

import (
	"errors"
	"sort"
	"sync"
)

type OrderedMap struct {
	data  map[string]interface{}
	order []string
	sync.RWMutex
}

// Create a new ordered map object
func New() OrderedMap {
	return OrderedMap{
		data: make(map[string]interface{}),
	}
}

// Add an object onto the end of the map
func (m OrderedMap) Add(key string, value interface{}) {
	m.Lock()
	m.data[key] = value
	m.order = append(m.order, key)
	m.Unlock()
}

// Add an object to a specific position in the map
func (m OrderedMap) Insert(position int, key string, value interface{}) {
	m.Lock()
	m.data[key] = value
	m.order = append(m.order, "")
	copy(m.order[position+1:], m.order[position:])
	m.order[position] = key
	m.Unlock()
}

// Get a specific object out of the map based on its key
func (m OrderedMap) Get(key string) (interface{}, bool) {
	m.RLock()
	data, ok := m.data[key]
	m.RUnlock()
	return data, ok
}

// Get an unordered map of all objects currently stored within this one
func (m OrderedMap) GetAll() map[string]interface{} {
	tmp := make(map[string]interface{})
	m.RLock()
	for k, v := range m.data {
		tmp[k] = v
	}
	m.RUnlock()
	return tmp
}

// Get a slice of strings containing the current order
func (m OrderedMap) GetOrder() []string {
	m.RLock()
	tmp := make([]string, len(m.order))
	copy(tmp, m.order) 
	m.RUnlock()
	return tmp
}

// Set a new order for this map, must contain all keys.
func (m OrderedMap) SetOrder(order []string) error {
	if !compareOrder(m.order, order) {
		return errors.New("Provided order does not contain the same data as existing.")
	}
	m.Lock()
	copy(m.order, order)
	m.Unlock()
	return nil
}

// Get the index in the order of a specific key
func (m OrderedMap) IndexOf(key string) int {
	m.RLock()
	defer m.RUnlock()
	for i := 0; i < len(m.order); i++ {
		if m.order[i] == key {
			return i
		}
	}
	return -1
}

// Delete a specific key and all associated data from the map
func (m OrderedMap) Delete(key string) {
	m.Lock()
	delete(m.data, key)
	idx := m.IndexOf(key)
	m.order = m.order[:idx+copy(m.order[idx:], m.order[idx+1:])]
	m.Unlock()
}

// Get the total size of the map
func (m OrderedMap) Count() int {
	m.RLock()
	cnt := len(m.data)
	m.RUnlock()
	return cnt
}

// Compare two orders and determine if they have the same data even if not in the same order
func compareOrder(f []string, s []string) bool {
	// Check to see if the two slices have the same length, if not they obviously aren't the same
	if len(f) != len (s) {
		return false
	}

	// Let's copy our own slices, so we don't edit the originals
	var tmp1 []string
	var tmp2 []string
	copy(tmp1, f)
	copy(tmp2, s)

	// Now we'll sort them and go through each index to make sure it's the same
	sort.Strings(tmp1)
	sort.Strings(tmp2)
	for i, v := range f {
		if s[i] != v {
			return false
		}
	}

	return true
}
