package nex

import (
	"fmt"
)

type CountSet struct {
	Size   int
	Values []int
}

func (cs CountSet) Next() (int, int) {

	i := 0
	for j, x := range cs.Values {
		if x != i {
			return i, j
		}
		i++
	}
	return i, len(cs.Values)

}

func (cs CountSet) Add() (int, CountSet, error) {

	if len(cs.Values) == cs.Size {
		return -1, cs, fmt.Errorf("overflow")
	}
	i, j := cs.Next()
	cs.Values = append(cs.Values[:j], append([]int{i}, cs.Values[j:]...)...)
	return i, cs, nil

}

func (cs CountSet) Remove(i int) CountSet {
	for j, x := range cs.Values {
		if i == x {
			var tail []int
			if j < len(cs.Values)-1 {
				tail = cs.Values[j+1:]
			}
			cs.Values = append(cs.Values[:j], tail...)
			return cs
		}
	}
	return cs
}
