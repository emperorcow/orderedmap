package orderedmap

import (
	"errors"
	"sort"
	"sync"
)

type OrderedMap struct {
	data  map[string]interface{}
	order []string
	lock  sync.RWMutex
}

// Create a new ordered map object
func New() OrderedMap {
	return OrderedMap{
		data:  make(map[string]interface{}),
		order: make([]string, 0),
	}
}

// Add an object onto the end of the map
func (m *OrderedMap) Add(key string, value interface{}) {
	m.lock.Lock()
	m.data[key] = value
	m.order = append(m.order, key)
	m.lock.Unlock()
}

// Add an object to a specific position in the map
func (m *OrderedMap) Insert(position int, key string, value interface{}) error {
	if position >= len(m.data) {
		return errors.New("Position is larger than the current map size.")
	}

	if position < 0 {
		return errors.New("Position is less than 0.")
	}

	m.lock.Lock()
	m.data[key] = value
	pre := m.order[:position]
	post := m.order[position:]
	m.order = make([]string, len(pre))
	copy(m.order, pre)
	m.order = append(m.order, key)
	m.order = append(m.order, post...)
	m.lock.Unlock()

	return nil
}

// Get a specific object out of the map based on its key
func (m OrderedMap) GetKey(key string) (interface{}, bool) {
	m.lock.RLock()
	data, ok := m.data[key]
	m.lock.RUnlock()
	return data, ok
}

// Get a specific object and it's key out of the map based on it's order index
func (m OrderedMap) GetIndex(index int) (string, interface{}, bool) {
	m.lock.RLock()
	key := m.order[index]
	data, ok := m.data[key]
	m.lock.RUnlock()
	return key, data, ok
}

// Get a slice of strings containing the current order
func (m OrderedMap) GetOrder() []string {
	m.lock.RLock()
	tmp := make([]string, len(m.order))
	copy(tmp, m.order)
	m.lock.RUnlock()
	return tmp
}

// Set a new order for this map, must contain all keys.
func (m *OrderedMap) SetOrder(order []string) error {
	if !compareOrder(m.order, order) {
		return errors.New("Provided order does not contain the same data as existing.")
	}
	m.lock.Lock()
	copy(m.order, order)
	m.lock.Unlock()
	return nil
}

// Get the index in the order of a specific key
func (m OrderedMap) IndexOf(key string) int {
	m.lock.RLock()
	index := -1
	for i := 0; i < len(m.order); i++ {
		if m.order[i] == key {
			index = i
		}
	}
	m.lock.RUnlock()
	return index
}

// Delete a specific key and all associated data from the map
func (m *OrderedMap) Delete(key string) {
	idx := m.IndexOf(key)

	m.lock.Lock()
	delete(m.data, key)
	tmp := make([]string, len(m.order))
	copy(tmp, m.order)
	m.order = make([]string, len(tmp))

	m.order = append(tmp[:idx], tmp[idx+1:]...)
	m.lock.Unlock()
}

// Get the total size of the map
func (m OrderedMap) Count() int {
	m.lock.RLock()
	cnt := len(m.data)
	m.lock.RUnlock()
	return cnt
}

type OrderedMapIterator struct {
	returnchan chan Tuple
	breakchan  chan bool
	data       *OrderedMap
}

type Tuple struct {
	Key string
	Val interface{}
}

func (m *OrderedMap) Iterator() OrderedMapIterator {
	return OrderedMapIterator{
		returnchan: make(chan Tuple),
		breakchan:  make(chan bool),
		data:       m,
	}
}

func (it *OrderedMapIterator) Loop() <-chan Tuple {
	go func() {
		max := it.data.Count()

		for i := 0; i < max; i++ {
			k, v, ok := it.data.GetIndex(i)
			if ok {
				it.returnchan <- Tuple{k, v}
			}
		}

		close(it.returnchan)
	}()

	return it.returnchan
}

func (it *OrderedMapIterator) Break() {

}

// Compare two orders and determine if they have the same data even if not in the same order
func compareOrder(f []string, s []string) bool {
	// Check to see if the two slices have the same length, if not they obviously aren't the same
	if len(f) != len(s) {
		return false
	}

	// Let's copy our own slices, so we don't edit the originals
	tmpf := make([]string, len(f))
	tmps := make([]string, len(s))
	copy(tmpf, f)
	copy(tmps, s)

	// Now we'll sort them and go through each index to make sure it's the same
	sort.Strings(tmpf)
	sort.Strings(tmps)

	for i, v := range tmps {
		if tmpf[i] != v {
			return false
		}
	}

	return true
}
