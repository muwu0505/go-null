package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Int64 struct {
	sql.NullInt64
}

func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

func SetInt64(i int64) Int64 {
	return NewInt64(i, true)
}

func SetInt64Ptr(i *int64) Int64 {
	if i == nil {
		return NewInt64(0, false)
	}
	return NewInt64(*i, true)
}

func (i Int64) ToInt64() int64 {
	if i.Valid == false {
		return 0
	}

	return i.Int64
}

func (i Int64) IsNull() bool {
	return i.Valid == false
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &i.Int64); err != nil {
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
			n, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return fmt.Errorf("null: couldn't convert string to int: %w", err)
			}
			i.Int64 = n
			i.Valid = true
			return nil
		}
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int64) UnmarshalText(text []byte) error {
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
	i.Int64 = val
	i.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int64), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int64) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int64), 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int64) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int64) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// Equal returns true if both ints have the same value or are both null.
func (i Int64) Equal(other Int64) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Int64 == other.Int64)
}
