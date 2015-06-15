# OrderedMap #
Provides a map that stores its items in a specific order that can be used
in a protected and concurrent fashion. Map key must be a string, but the data
can be anything.

[![GoDoc](https://godoc.org/github.com/emperorcow/orderedmap?status.svg)](http://godoc.org/github.com/emperorcow/orderedmap)
[![Build Status](https://drone.io/github.com/emperorcow/orderedmap/status.png)](https://drone.io/github.com/emperorcow/orderedmap/latest)
[![Coverage Status](https://coveralls.io/repos/emperorcow/orderedmap/badge.svg?branch=master)](https://coveralls.io/r/emperorcow/orderedmap?branch=master)

# Usage #
You'll first need to make a new map and add some data.  Create a new map using the New() function and then add some data in using the Set(Key, Val) function, which takes a key as a string, and any data type as the value:
```
om := orderedmap.New()
om.Set("one", TestData{ID: 1, Name: "one"})
om.Set("two", TestData{ID: 2, Name: "two"})
om.Set("three", TestData{ID: 3, Name: "three"})
```

You can get data from the map using the GetKey or GetIndex.  GetKey pulls data using the key string, GetIndex uses the current location, in order, from the map. 

```
datakey, ok := om.GetKey("two")

dataidx, ok := om.GetIndex(1)
```

There are also many other things, you can do, like delete by key, get the size, etc.  See the godoc for more information

