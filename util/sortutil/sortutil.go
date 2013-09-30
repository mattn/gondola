// Package sortutil contains functions for easily
// sorting slices while avoiding a lot of boilerplate.
//
// Keep in mind though, that this convenience has a significant
// runtime penalty, so you shouldn't use it for long lists.
package sortutil

import (
	"fmt"
	"reflect"
	"sort"
)

type sortable struct {
	value      reflect.Value
	key        string
	descending bool
}

func (s *sortable) Len() int {
	return s.value.Len()
}

func (s *sortable) less(i, j int) bool {
	fi := reflect.Indirect(s.value.Index(i)).FieldByName(s.key)
	fj := reflect.Indirect(s.value.Index(j)).FieldByName(s.key)
	switch fi.Kind() {
	case reflect.Bool:
		return !fi.Bool() && fj.Bool()
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return fi.Int() < fj.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return fi.Uint() < fj.Uint()
	case reflect.Float32, reflect.Float64:
		return fi.Float() < fj.Float()
	case reflect.String:
		return fi.String() < fj.String()
	default:
		panic(fmt.Errorf("can't compare type %s", fi.Type()))
	}
	panic("unreachable")
}

func (s *sortable) Less(i, j int) bool {
	v := s.less(i, j)
	if s.descending {
		return !v
	}
	return v
}

func (s *sortable) Swap(i, j int) {
	vi := s.value.Index(i)
	vj := s.value.Index(j)
	tmp := reflect.New(vi.Type()).Elem()
	tmp.Set(vi)
	vi.Set(vj)
	vj.Set(tmp)
}

// Sort sorts an array or slice of structs or pointer to
// structs by comparing the given key, which must be a
// an exported struct field. If the key is prefixed by
// the character '-', the sorting is performed in
// descending order.
func Sort(data interface{}, key string) (err error) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return fmt.Errorf("can't short type %v, must be slice or array", val.Type())
	}
	descending := false
	if key != "" && key[0] == '-' {
		descending = true
		key = key[1:]
	}
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	sort.Sort(&sortable{val, key, descending})
	return err
}
