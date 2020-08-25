package memdb

import "fmt"

type DBSchema struct {
	Tables map[string]*TableSchema
}

type TableSchema struct {
	Name string
	Indexes map[string]*IndexSchema
}

type IndexSchema struct {
	Name string
	AllowMissing bool
	Unique bool
	Indexer Indexer // Index Interface
}


// Validate a DB schema
func (s *DBSchema) Validate() error {
	// if invalid schema is passed
	if s == nil {
		return fmt.Errorf("schema is nil")
	}
	// if there are no tables passed in to the db
	if len(s.Tables) == 0 {
		fmt.Errorf("schema has no table defined")
	}
	// check if table name matches
	for name, table := range s.Tables {
		if name != table.Name {
			fmt.Errorf("table name mis match for %s", name)
		}
	}
	return nil
}

// Validate a Table schema
func (s *TableSchema) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("missing table name")
	}
	if len(s.Indexes) == 0 {
		return fmt.Errorf("missing table indexes for %s", s.Name)
	}
	if _, ok := s.Indexes["id"]; !ok {
		return fmt.Errorf("must have id index")
	}
	if !s.Indexes["id"].Unique {
		return fmt.Errorf("id index must be unique")
	}
	if _, ok := s.Indexes["id"].Indexer.(SingleIndexer); !ok {
		return fmt.Errorf("id index must be a Single Indexer")
	}
	for name, index := range s.Indexes {
		if name != index.Name {
			return fmt.Errorf("index name mismatch for %s", name)
		}
		if err := index.Validate(); err != nil {
			return fmt.Errorf("index %q: %s", name, err)
		}
	}
	return nil
}
// Validate an index schema
func (s *IndexSchema) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("missing index name")
	}
	if s.Indexer == nil {
		return fmt.Errorf("missing index function for '%s'", s.Name)
	}
	switch s.Indexer.(type) {
	case SingleIndexer:
	case MultiIndexer:
	default:
		return fmt.Errorf("indexer for %s must be a SingleIndexer or MultiIndexer", s.Name)
	}
	return nil
}