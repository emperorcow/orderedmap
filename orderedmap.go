/*
Provides a map that stores its items in a specific order that can be used
in a protected and concurrent fashion. Map key must be a string, but the data
can be anything.

Let's start by first getting a new orderedmap and adding a few new things to it.
You must get a map by using the New() function, and can add things to it using
Set(Key, Value):

	om := orderedmap.New()
	om.Set("one", 1)
	om.Set("two", 2)

Now we can access our data through either it's position in order or the map key
using the GetKey and GetIndex functions:

	datakey := om.GetKey("one")
	dataidx := om.GetIndex(0)

*/
package orderedmap

import (
	"errors"
	"sort"
	"sync"
)

// A map structure that stores data within an ordered fashion.
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

// Add an object to a specific position in the map.  Position is zero indexed,
// so to add to the very beginning, you would use 0, to add to the end you would
// use Count() - 1.
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

// Get a specific object out of the map based on its map key.  In the event the
// key does not exist or the data is out of range, the function will have a
// second return of false.
//
// Get key:
// 	data, ok := om.GetKey("mykey")
//
// Test if key exists:
// 	if _, ok := om.GetKey("mykey"); ok {
// 		... DO SOMETHING HERE ...
// 	}
func (m OrderedMap) GetKey(key string) (interface{}, bool) {
	m.lock.RLock()
	data, ok := m.data[key]
	m.lock.RUnlock()
	return data, ok
}

// Get a specific object and it's key out of the map based on it's order index,
// with 0 being the first item in the order.  Will return a false in the event
// The key does not exist.
func (m OrderedMap) GetIndex(index int) (string, interface{}, bool) {
	m.lock.RLock()
	key := m.order[index]
	data, ok := m.data[key]
	m.lock.RUnlock()
	return key, data, ok
}

// Get a slice of strings containing the current order of the array
func (m OrderedMap) GetOrder() []string {
	m.lock.RLock()
	tmp := make([]string, len(m.order))
	copy(tmp, m.order)
	m.lock.RUnlock()
	return tmp
}

// Set a new order for this map.  SetOrder will return an error if either the
// number of items in the provided slice is different than those in the map, or
// if the keys are different that those currently in use.
func (m *OrderedMap) SetOrder(order []string) error {
	if !compareOrder(m.order, order) {
		return errors.New("Provided order does not contain the same data as existing.")
	}
	m.lock.Lock()
	copy(m.order, order)
	m.lock.Unlock()
	return nil
}

// Get the order index of a specific key
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

// A struct used to provide the ability to loop through all items in the
// orderedmap in order.
type OrderedMapIterator struct {
	returnchan chan Tuple
	breakchan  chan bool
	data       *OrderedMap
}

// A data structure to hold returned information on each iteration
type Tuple struct {
	Key string
	Val interface{}
}

// Returns an OrderedMapIterator type that can be used to loop through the
// entire map, in order.
//
// This function will return a struct with two functions that should be used
// to iterate through the map: Loop() and Break().  Loop() should be provided
// to range and will return a Tuple for each item in the map.
//
// IMPORTANT NOTE: You must use the Break() function before you use the break
// go command, otherwise you might have deadlock, race, or garbage issues.
func (m *OrderedMap) Iterator() OrderedMapIterator {
	return OrderedMapIterator{
		returnchan: make(chan Tuple),
		breakchan:  make(chan bool),
		data:       m,
	}
}

// Provides access to a channel that will allow looping through the entire
// map in order.  Returns a channel that can be passed to range and returns a
// Tuple struct with the key and value of each item.
//
// 	iter = mymap.Iterator()
// 	for data := range iter.Loop() {
// 		fmt.Printf("%s > %v\n", data.Key, data.Val)
// 	}
func (it *OrderedMapIterator) Loop() <-chan Tuple {
	go func() {
		max := it.data.Count()

		for i := 0; i < max; i++ {
			k, v, ok := it.data.GetIndex(i)
			if ok {
				select {
				case it.returnchan <- Tuple{k, v}:
				case <-it.breakchan:
					close(it.returnchan)
					return
				}
			}
		}

		close(it.returnchan)
		close(it.breakchan)
	}()

	return it.returnchan
}

// Signals the iterator that you no longer want to loop, allowing us to clean
// up, stop looping, and allows the garbage collector to clean up.  Finally,
// also makes sure all channels are closed and all mutex locks are clean, so
// that there are no issues with deadlocks.
func (it *OrderedMapIterator) Break() {
	select {
	case _, _ = <-it.breakchan:
	default:
		it.breakchan <- true
	}
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
