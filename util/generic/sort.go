package generic

import (
	"sort"
)

type sortable struct {
	length int
	value  handle
	fn     mapFunc
	cmp    lessFunc
	idx    indexFunc
	sw     swapFunc
}

func (s *sortable) Len() int {
	return s.length
}

func (s *sortable) Less(i, j int) bool {
	vi := s.fn(s.idx(s.value, i))
	vj := s.fn(s.idx(s.value, j))
	return s.cmp(vi, vj)
}

func (s *sortable) Swap(i, j int) {
	s.sw(s.value, i, j)
}

type reverseSortable struct {
	*sortable
}

func (s *reverseSortable) Less(i, j int) bool {
	return !s.sortable.Less(i, j)
}

// Sort sorts an array or slice of structs or pointer to
// structs by comparing the given key, which must be a
// an exported struct field or an exported method with no
// arguments and just one return value. If the key is
// prefixed by the character '-', the sorting is performed
// in descending order. If there are any errors, Sort panics
// since they can't be anything but programming errors.
func Sort(data interface{}, key string) {
	descending := false
	if key != "" && key[0] == '-' {
		descending = true
		key = key[1:]
	}
	fn, val, elem, typ, err := sliceMapper(data, key)
	if err != nil {
		panic(err)
	}
	if fn == nil {
		// Empty slice
		return
	}
	cmp, err := lessComparator(typ)
	if err != nil {
		panic(err)
	}
	srt := &sortable{val.Len(), getHandle(val), fn, cmp, indexer(elem), swapper(elem)}
	if descending {
		sort.Sort(&reverseSortable{srt})
	} else {
		sort.Sort(srt)
	}
}