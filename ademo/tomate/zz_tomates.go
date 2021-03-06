package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	json "encoding/json"
)

// Tomates implements a typed slice of *Tomate
type Tomates struct{ items []*Tomate }

// NewTomates creates a new typed slice of *Tomate
func NewTomates() *Tomates {
	return &Tomates{items: []*Tomate{}}
}

// Push appends every *Tomate
func (t *Tomates) Push(x ...*Tomate) *Tomates {
	t.items = append(t.items, x...)
	return t
}

// Unshift prepends every *Tomate
func (t *Tomates) Unshift(x ...*Tomate) *Tomates {
	t.items = append(x, t.items...)
	return t
}

// Pop removes then returns the last *Tomate.
func (t *Tomates) Pop() *Tomate {
	var ret *Tomate
	if len(t.items) > 0 {
		ret = t.items[len(t.items)-1]
		t.items = append(t.items[:0], t.items[len(t.items)-1:]...)
	}
	return ret
}

// Shift removes then returns the first *Tomate.
func (t *Tomates) Shift() *Tomate {
	var ret *Tomate
	if len(t.items) > 0 {
		ret = t.items[0]
		t.items = append(t.items[:0], t.items[1:]...)
	}
	return ret
}

// Index of given *Tomate. It must implements Ider interface.
func (t *Tomates) Index(s *Tomate) int {
	ret := -1
	for i, item := range t.items {
		if s.GetID() == item.GetID() {
			ret = i
			break
		}
	}
	return ret
}

// Contains returns true if s in is t.
func (t *Tomates) Contains(s *Tomate) bool {
	return t.Index(s) > -1
}

// RemoveAt removes a *Tomate at index i.
func (t *Tomates) RemoveAt(i int) bool {
	if i >= 0 && i < len(t.items) {
		t.items = append(t.items[:i], t.items[i+1:]...)
		return true
	}
	return false
}

// Remove removes given *Tomate
func (t *Tomates) Remove(s *Tomate) bool {
	if i := t.Index(s); i > -1 {
		t.RemoveAt(i)
		return true
	}
	return false
}

// InsertAt adds given *Tomate at index i
func (t *Tomates) InsertAt(i int, s *Tomate) *Tomates {
	if i < 0 || i >= len(t.items) {
		return t
	}
	res := []*Tomate{}
	res = append(res, t.items[:0]...)
	res = append(res, s)
	res = append(res, t.items[i:]...)
	t.items = res
	return t
}

// Splice removes and returns a slice of *Tomate, starting at start, ending at start+length.
// If any s is provided, they are inserted in place of the removed slice.
func (t *Tomates) Splice(start int, length int, s ...*Tomate) []*Tomate {
	var ret []*Tomate
	for i := 0; i < len(t.items); i++ {
		if i >= start && i < start+length {
			ret = append(ret, t.items[i])
		}
	}
	if start >= 0 && start+length <= len(t.items) && start+length >= 0 {
		t.items = append(
			t.items[:start],
			append(s,
				t.items[start+length:]...,
			)...,
		)
	}
	return ret
}

// Slice returns a copied slice of *Tomate, starting at start, ending at start+length.
func (t *Tomates) Slice(start int, length int) []*Tomate {
	var ret []*Tomate
	if start >= 0 && start+length <= len(t.items) && start+length >= 0 {
		ret = t.items[start : start+length]
	}
	return ret
}

// Reverse the slice.
func (t *Tomates) Reverse() *Tomates {
	for i, j := 0, len(t.items)-1; i < j; i, j = i+1, j-1 {
		t.items[i], t.items[j] = t.items[j], t.items[i]
	}
	return t
}

// Len of the slice.
func (t *Tomates) Len() int {
	return len(t.items)
}

// Set the slice.
func (t *Tomates) Set(x []*Tomate) *Tomates {
	t.items = append(t.items[:0], x...)
	return t
}

// Get the slice.
func (t *Tomates) Get() []*Tomate {
	return t.items
}

// At return the item at index i.
func (t *Tomates) At(i int) *Tomate {
	return t.items[i]
}

// Filter return a new Tomates with all items satisfying f.
func (t *Tomates) Filter(filters ...func(*Tomate) bool) *Tomates {
	ret := NewTomates()
	for _, i := range t.items {
		ok := true
		for _, f := range filters {
			ok = ok && f(i)
			if !ok {
				break
			}
		}
		if ok {
			ret.Push(i)
		}
	}
	return ret
}

// Map return a new Tomates of each items modified by f.
func (t *Tomates) Map(mappers ...func(*Tomate) *Tomate) *Tomates {
	ret := NewTomates()
	for _, i := range t.items {
		val := i
		for _, m := range mappers {
			val = m(val)
			if val == nil {
				break
			}
		}
		if val != nil {
			ret.Push(val)
		}
	}
	return ret
}

// First returns the first value or default.
func (t *Tomates) First() *Tomate {
	var ret *Tomate
	if len(t.items) > 0 {
		ret = t.items[0]
	}
	return ret
}

// Last returns the last value or default.
func (t *Tomates) Last() *Tomate {
	var ret *Tomate
	if len(t.items) > 0 {
		ret = t.items[len(t.items)-1]
	}
	return ret
}

// Empty returns true if the slice is empty.
func (t *Tomates) Empty() bool {
	return len(t.items) == 0
}

// NotEmpty returns true if the slice is not empty.
func (t *Tomates) NotEmpty() bool {
	return len(t.items) > 0
}

// Transact execute one op.
func (t *Tomates) Transact(F ...func(*Tomates)) {
	for _, f := range F {
		f(t)
	}
}

//UnmarshalJSON JSON unserializes Tomates
func (t *Tomates) UnmarshalJSON(b []byte) error {
	var items []*Tomate
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}
	t.items = items
	return nil
}

//MarshalJSON JSON serializes Tomates
func (t *Tomates) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.items)
}

// TomatesContract are the requirements of Tomates
type TomatesContract interface {
	Push(x ...*Tomate) *Tomates
	Unshift(x ...*Tomate) *Tomates
	Pop() *Tomate
	Shift() *Tomate
	Index(s *Tomate) int
	Contains(s *Tomate) bool
	RemoveAt(i int) bool
	Remove(s *Tomate) bool
	InsertAt(i int, s *Tomate) *Tomates
	Splice(start int, length int, s ...*Tomate) []*Tomate
	Slice(start int, length int) []*Tomate
	Reverse() *Tomates
	Set(x []*Tomate) *Tomates
	Get() []*Tomate
	At(i int) *Tomate
	Filter(filters ...func(*Tomate) bool) *Tomates
	Map(mappers ...func(*Tomate) *Tomate) *Tomates
	First() *Tomate
	Last() *Tomate
	Transact(...func(*Tomates))
	Len() int
	Empty() bool
	NotEmpty() bool
}

// FilterTomates provides filters for a struct.
var FilterTomates = struct {
	ByID     func(...string) func(*Tomate) bool
	NotID    func(...string) func(*Tomate) bool
	ByColor  func(...string) func(*Tomate) bool
	NotColor func(...string) func(*Tomate) bool
}{
	ByID: func(all ...string) func(*Tomate) bool {
		return func(o *Tomate) bool {
			for _, v := range all {
				if o.ID == v {
					return true
				}
			}
			return false
		}
	},
	NotID: func(all ...string) func(*Tomate) bool {
		return func(o *Tomate) bool {
			for _, v := range all {
				if o.ID == v {
					return false
				}
			}
			return true
		}
	},
	ByColor: func(all ...string) func(*Tomate) bool {
		return func(o *Tomate) bool {
			for _, v := range all {
				if o.Color == v {
					return true
				}
			}
			return false
		}
	},
	NotColor: func(all ...string) func(*Tomate) bool {
		return func(o *Tomate) bool {
			for _, v := range all {
				if o.Color == v {
					return false
				}
			}
			return true
		}
	},
}

// SetterTomates provides sets properties.
var SetterTomates = struct {
	SetID    func(string) func(*Tomate) *Tomate
	SetColor func(string) func(*Tomate) *Tomate
}{
	SetID: func(v string) func(*Tomate) *Tomate {
		return func(o *Tomate) *Tomate {
			o.ID = v
			return o
		}
	},
	SetColor: func(v string) func(*Tomate) *Tomate {
		return func(o *Tomate) *Tomate {
			o.Color = v
			return o
		}
	},
}
