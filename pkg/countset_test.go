package nex

import (
	"reflect"
	"testing"
)

func TestCountsetBasic(t *testing.T) {

	cs := CountSet{Size: 4}

	i, cs, err := cs.Add()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]int(cs.Values), []int{0}) {
		t.Errorf("%v != [0]", cs.Values)
	}
	if i != 0 {
		t.Errorf("%d != 0", i)
	}

	i, cs, err = cs.Add()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 1}) {
		t.Errorf("%v != [0 1]", cs.Values)
	}
	if i != 1 {
		t.Errorf("%d != 1", i)
	}

	i, cs, err = cs.Add()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 1, 2}) {
		t.Errorf("%v != [0 1 2]", cs.Values)
	}
	if i != 2 {
		t.Errorf("%d != 2", i)
	}

	cs = cs.Remove(3) // does not exist, should not change anything
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 1, 2}) {
		t.Errorf("%v != [0 1 2]", cs.Values)
	}

	cs = cs.Remove(1)
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 2}) {
		t.Errorf("%v != [0 2]", cs.Values)
	}

	i, cs, err = cs.Add()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 1, 2}) {
		t.Errorf("%v != [0 1 2]", cs.Values)
	}
	if i != 1 {
		t.Errorf("%d != 1", i)
	}

	i, cs, err = cs.Add()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 1, 2, 3}) {
		t.Errorf("%v != [0 1 2 3]", cs.Values)
	}
	if i != 3 {
		t.Errorf("%d != 3", i)
	}

	cs = cs.Remove(2)
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 1, 3}) {
		t.Errorf("%v != [0 1 3]", cs.Values)
	}

	cs = cs.Remove(1)
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 3}) {
		t.Errorf("%v != [0 1 3]", cs.Values)
	}

	i, cs, err = cs.Add()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 1, 3}) {
		t.Errorf("%v != [0 1 3]", cs.Values)
	}
	if i != 1 {
		t.Errorf("%d != 1", i)
	}

	i, cs, err = cs.Add()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]int(cs.Values), []int{0, 1, 2, 3}) {
		t.Errorf("%v != [0 1 2 3]", cs.Values)
	}
	if i != 2 {
		t.Errorf("%d != 2", i)
	}

	i, cs, err = cs.Add()
	if err == nil {
		t.Fatal("expected overflow error")
	}

}
