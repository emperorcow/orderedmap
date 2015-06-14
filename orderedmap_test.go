package orderedmap

import (
	"reflect"
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
		t.Error("Map does not contain three items")
	}

	tmp := om.GetOrder()
	if tmp[1] != "two" {
		t.Logf("Order: %v\n", tmp)
		t.Error("Index two is not the correct object name")
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

	test, ok := om.GetIndex(1)
	gotten := test.(TestData)

	if !ok {
		t.Error("Unable to get item from map by index")
	}

	if gotten.ID != 2 || gotten.Name != "two" {
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
