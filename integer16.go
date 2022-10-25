package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Int16 struct {
	sql.NullInt16
}

func NewInt16(i int16, valid bool) Int16 {
	return Int16{
		NullInt16: sql.NullInt16{
			Int16: i,
			Valid: valid,
		},
	}
}

func SetInt16(i int16) Int16 {
	return NewInt16(i, true)
}

func SetInt16Ptr(i *int16) Int16 {
	if i == nil {
		return NewInt16(0, false)
	}
	return NewInt16(*i, true)
}

func (i Int16) ToInt16() int16 {
	if i.Valid == false {
		return 0
	}

	return i.Int16
}

func (i Int16) IsNull() bool {
	return i.Valid == false
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int16) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &i.Int16); err != nil {
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
			n, err := strconv.ParseInt(str, 10, 16)
			if err != nil {
				return fmt.Errorf("null: couldn't convert string to int: %w", err)
			}
			i.Int16 = int16(n)
			i.Valid = true
			return nil
		}
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int16) UnmarshalText(text []byte) error {
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
	i.Int16 = int16(val)
	i.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (i Int16) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int16), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int16) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int16), 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int16) SetValid(n int16) {
	i.Int16 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int16) Ptr() *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

// Equal returns true if both ints have the same value or are both null.
func (i Int16) Equal(other Int16) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Int16 == other.Int16)
}
