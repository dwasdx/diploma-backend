package models

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

type StringNull string

func (s *StringNull) IsEmpty() bool {
	if len(*s) > 0 {
		return false
	}

	return true
}

func (s *StringNull) SqlValue() interface{} {
	if !s.IsEmpty() {
		return s
	}

	return nil
}

func (s *StringNull) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	//v := fmt.Sprintf("%v", value)
	//*s = StringNull(v)

	return nil
}

type NullString sql.NullString

func (ns *NullString) IsEmpty() bool {
	if len(ns.String) > 0 {
		return false
	}

	return true
}

func (ns *NullString) SqlValue() interface{} {
	if !ns.IsEmpty() {
		return ns.String
	}

	return nil
}

func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid || ns.IsEmpty() {
		return []byte("null"), nil
	}

	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}
