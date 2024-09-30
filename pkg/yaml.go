package nex

import (
	"fmt"
	"io"
	"os"

	"github.com/mergetb/yaml/v3"
)

type Kind string

const (
	NetworkKind    = "Network"
	MemberListKind = "MemberList"
)

type MetaObject struct {
	Kind   string
	Object interface{}
}

func (o *MetaObject) UnmarshalYAML(unmarshal func(interface{}) error) error {

	obj := struct{ Kind string }{}
	err := unmarshal(&obj)
	if err != nil {
		return err
	}
	o.Kind = obj.Kind

	switch obj.Kind {
	case "Network":
		net := &Network{}
		err := unmarshal(&net)
		if err != nil {
			return err
		}
		o.Object = net
	case "MemberList":
		ml := &MemberList{}
		err := unmarshal(&ml)
		if err != nil {
			return err
		}
		o.Object = ml
	default:
		return fmt.Errorf("unknown kind: %s", obj.Kind)
	}

	return nil

}

func ReadSpec(file string) ([]MetaObject, error) {

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var objects []MetaObject
	dc := yaml.NewDecoder(f)

	err = nil
	for err == nil {

		var object MetaObject
		err = dc.Decode(&object)
		if err == nil {
			objects = append(objects, object)
		}

	}

	if err != io.EOF {
		return nil, err
	}

	return objects, nil

}
