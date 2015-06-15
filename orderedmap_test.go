// Package orderedmap provides a map where the order of items is maintained.
// Furthermore, access to contained data is done in a way that is protected and
// concurrent.
package orderedmap

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type TestData struct {
	ID   int
	Name string
}

func TestNewOrderedMap(t *testing.T) {
	om := New()
	if reflect.TypeOf(om).Name() != "OrderedMap" {

		t.Error("Map is not the correct type")
	}

	if om.Count() != 0 {
		t.Error("New map is not empty")
	}
}

func TestAdd(t *testing.T) {
	om := New()
	one := TestData{ID: 1, Name: "one"}
	two := TestData{ID: 2, Name: "two"}

	om.Add("one", one)
	om.Add("two", two)

	if om.Count() != 2 {
		t.Error("Map does not contain two items")
	}
}

func TestInsert(t *testing.T) {
	om := New()

	om.Add("one", TestData{ID: 1, Name: "one"})
	om.Add("three", TestData{ID: 3, Name: "three"})
	om.Add("four", TestData{ID: 4, Name: "four"})
	om.Add("five", TestData{ID: 5, Name: "five"})

	err := om.Insert(1, "two", TestData{ID: 2, Name: "two"})

	if err != nil {
		t.Error("Error trying to insert into ordered map: " + err.Error())
	}

	if om.Count() != 5 {
		t.Error("Map does not contain correct number of items")
	}

	tmp := om.GetOrder()
	if tmp[1] != "two" {
		t.Logf("Order: %v\n", tmp)
		t.Error("Index two is not the correct object name")
	}

	err = om.Insert(30, "six", TestData{ID: 6, Name: "six"})
	if err == nil {
		t.Error("No error was received when trying to insert above the range.")
	}
	err = om.Insert(-1, "six", TestData{ID: 6, Name: "six"})
	if err == nil {
		t.Error("No error was received when trying to insert negative value.")
	}
}

func TestGetKey(t *testing.T) {
	om := New()

	om.Add("one", TestData{ID: 1, Name: "one"})
	om.Add("two", TestData{ID: 2, Name: "two"})
	om.Add("three", TestData{ID: 3, Name: "three"})

	test, ok := om.GetKey("two")
	gotten := test.(TestData)

	if !ok {
		t.Error("Unable to get item from map by key")
	}

	if gotten.ID != 2 || gotten.Name != "two" {
		t.Error("Wrong item was returned from map")
	}
}

func TestGetIndex(t *testing.T) {
	om := New()
	om.Add("one", TestData{ID: 1, Name: "one"})
	om.Add("two", TestData{ID: 2, Name: "two"})
	om.Add("three", TestData{ID: 3, Name: "three"})

	key, val, ok := om.GetIndex(1)
	gotten := val.(TestData)

	if !ok {
		t.Error("Unable to get item from map by index")
	}

	if key != "two" || gotten.ID != 2 || gotten.Name != "two" {
		t.Error("Wrong item was returned from map")
	}
}

func TestGetOrder(t *testing.T) {
	om := New()
	om.Add("one", TestData{ID: 1, Name: "one"})
	om.Add("two", TestData{ID: 2, Name: "two"})
	om.Add("three", TestData{ID: 3, Name: "three"})

	ord := om.GetOrder()
	if ord[0] != "one" {
		t.Error("First item was wrong")
	} else if ord[1] != "two" {
		t.Error("Second item was wrong")
	} else if ord[2] != "three" {
		t.Error("Third item was wrong")
	}
}

func TestSetOrder(t *testing.T) {
	om := New()
	om.Add("one", TestData{ID: 1, Name: "one"})
	om.Add("two", TestData{ID: 2, Name: "two"})
	om.Add("three", TestData{ID: 3, Name: "three"})

	err := om.SetOrder([]string{"three", "one", "two"})
	if err != nil {
		t.Error("An error occured setting order: " + err.Error())
	}

	ord := om.GetOrder()
	if ord[0] != "three" {
		t.Error("First item was wrong")
	} else if ord[1] != "one" {
		t.Error("Second item was wrong")
	} else if ord[2] != "two" {
		t.Error("Third item was wrong")
	}

	err = om.SetOrder([]string{"three", "one", "two", "five", "eleventy"})
	if err == nil {
		t.Error("No error occured when trying to use an order that was too large")
	}

	err = om.SetOrder([]string{"three", "one", "five"})
	if err == nil {
		t.Error("No error occured when trying to use an order the right size, but with the wrong items")
	}
}

func TestIndexOf(t *testing.T) {
	om := New()
	om.Add("one", TestData{ID: 1, Name: "one"})
	om.Add("two", TestData{ID: 2, Name: "two"})
	om.Add("three", TestData{ID: 3, Name: "three"})

	idx := om.IndexOf("one")
	if idx != 0 {
		t.Error("Index of one was not 0")
	}
	idx = om.IndexOf("three")
	if idx != 2 {
		t.Error("Index of three was not 2")
	}
}

func TestDelete(t *testing.T) {
	om := New()
	om.Add("one", TestData{ID: 1, Name: "one"})
	om.Add("two", TestData{ID: 2, Name: "two"})
	om.Add("three", TestData{ID: 3, Name: "three"})

	om.Delete("two")
	_, ok := om.GetKey("two")
	if ok {
		t.Error("Deleted key still exists")
	}
	if om.Count() != 2 {
		t.Error("Size of ordered map was wrong")
	}
}

func TestCount(t *testing.T) {
	om := New()
	om.Add("one", TestData{ID: 1, Name: "one"})
	if om.Count() != 1 {
		t.Error("First count was wrong")
	}
	om.Add("two", TestData{ID: 2, Name: "two"})
	if om.Count() != 2 {
		t.Error("Second count was wrong")
	}
	om.Add("three", TestData{ID: 3, Name: "three"})
	if om.Count() != 3 {
		t.Error("Third count was wrong")
	}
}

func TestIterator(t *testing.T) {
	om := New()
	for i := 0; i < 100; i++ {
		str := strconv.Itoa(i)
		om.Add(str, TestData{ID: i, Name: str})
	}

	itr := om.Iterator()
	j := 0
	for item := range itr.Loop() {
		if item.Key != strconv.Itoa(j) {
			t.Errorf("Index %v did not match", j)
		}
		j++
	}
}

func TestIteratorBreak(t *testing.T) {
	om := New()
	for i := 0; i < 1000; i++ {
		str := strconv.Itoa(i)
		om.Add(str, TestData{ID: i, Name: str})
	}

	itr := om.Iterator()
	j := 0
	for _ = range itr.Loop() {
		if j == 60 {
			itr.Break()
			break
		}
		j++
	}
}

func ExampleIterator_full() {
	om := New()

	om.Add("1", "one")
	om.Add("2", "two")
	om.Add("3", "three")
	om.Add("4", "four")
	om.Add("5", "five")

	iter := om.Iterator()
	for data := range iter.Loop() {
		fmt.Printf("%s > %v\n", data.Key, data.Val)
	}
}

func ExampleIterator_break() {
	om := New()

	om.Add("1", "one")
	om.Add("2", "two")
	om.Add("3", "three")
	om.Add("4", "four")
	om.Add("5", "five")

	iter := om.Iterator()
	for data := range iter.Loop() {
		if data.Key == "3" {
			iter.Break()
			break
		}
		fmt.Printf("%s > %v\n", data.Key, data.Val)
	}
}
