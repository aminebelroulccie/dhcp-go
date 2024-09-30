package nex

import (
	"fmt"
)

func CheckDupes(objs []Object) error {

	m := make(map[string]bool)

	for _, o := range objs {

		_, ok := m[o.Key()]
		if ok {
			return fmt.Errorf("duplicate object '%s'", o.Key())
		}
		m[o.Key()] = true

	}

	return nil

}
