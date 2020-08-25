package memdb

import "reflect"

// Indexer is an interface used for defining indexes. Indexes are efficient for table lookups.


type Indexer interface {
	// FromArgs is used to build the exact key from the list of arguments. Using a veriadic function.
	FromArgs(args ...interface{}) ([]byte, error)
}

type PrefixIndexer interface {
	PrefixFromArgs(args ...interface{}) ([]byte, error)
}

type SingleIndexer interface {
	FromObject(raw interface{}) (bool, []byte, error)
}

type MultiIndexer interface {
	FromObject(raw interface{}) (bool, []byte, error)
}

type StringFieldIndex struct {
	Field string
	Lowercase bool
}

func (s *StringFieldIndex) FromObject(obj interface{}) (bool, []byte, error) {
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Deference the pointer
	fv := v.FieldByName(s.Field)
	isPtr := fv.Kind() == reflect.Ptr

}
