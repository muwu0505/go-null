package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Int32 struct {
	sql.NullInt32
}

func NewInt32(i int32, valid bool) Int32 {
	return Int32{
		NullInt32: sql.NullInt32{
			Int32: i,
			Valid: valid,
		},
	}
}

func SetInt32(i int32) Int32 {
	return NewInt32(i, true)
}

func SetInt32Ptr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}
	return NewInt32(*i, true)
}

func (i Int32) ToInt32() int32 {
	if i.Valid == false {
		return 0
	}

	return i.Int32
}

func (i Int32) IsNull() bool {
	return i.Valid == false
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int32) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &i.Int32); err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			// special case: accept string input
			if typeError.Value != "string" {
				return fmt.Errorf("null: JSON input is invalid type (need int or string): %w", err)
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
			}
			n, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				return fmt.Errorf("null: couldn't convert string to int: %w", err)
			}
			i.Int32 = int32(n)
			i.Valid = true
			return nil
		}
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	val, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	i.Int32 = int32(val)
	i.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (i Int32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int32) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int32) SetValid(n int32) {
	i.Int32 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int32) Ptr() *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

// Equal returns true if both ints have the same value or are both null.
func (i Int32) Equal(other Int32) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Int32 == other.Int32)
}
