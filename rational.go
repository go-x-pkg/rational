package rational

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Rational number numerator/denominator.
type Rational struct {
	Num uint64 `json:"num" yaml:"num" bson:"num"`
	Den uint64 `json:"den" yaml:"den" bson:"den"`
}

const (
	sep     byte = ':'
	sepAlt  byte = '/'
	sepAlt2 byte = 'x'
)

// IsNil checks if Rational is nil.
func (r Rational) IsNil() bool { return r.Num == 0 && r.Den == 0 }

// String returns string representation.
func (r Rational) String() string {
	if r.Num == 0 {
		return "0"
	}

	return fmt.Sprintf("%d%c%d", r.Num, sep, r.Den)
}

// Reverse returns reversed Rational.
func (r Rational) Reverse() Rational {
	return Rational{r.Den, r.Num}
}

// Float64 returns float64 representation.
func (r Rational) Float64() float64 {
	if r.Den == 0 {
		return 0
	}

	return (float64(r.Num) / float64(r.Den))
}

// Percent returns percentage representation.
// e.g. 1/4 is 25%.
func (r Rational) Percent() float64 { return r.Float64() * 100 }

// Unmarshal is generic Unmarshaler. For example to use with
// GetBSON.
func (r *Rational) Unmarshal(fn func(interface{}) error) error {
	var raw string

	if err := fn(&raw); err != nil {
		return fmt.Errorf("error unmarshal rational: %w", err)
	}

	v, err := NewRational(raw)
	if err != nil {
		return fmt.Errorf("error parse rational: %w", err)
	}

	*r = v

	return nil
}

// MarshalJSON implements json Marshaler.
func (r Rational) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", r.String())), nil
}

// MarshalYAML implements yaml Marshaler.
func (r Rational) MarshalYAML() (interface{}, error) { return r.String(), nil }

// UnmarshalJSON implements json Unmarshaler.
func (r *Rational) UnmarshalJSON(data []byte) error {
	return r.Unmarshal(func(r interface{}) error { return json.Unmarshal(data, r) })
}

// UnmarshalYAML implements yaml Unmarshaler.
func (r *Rational) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return r.Unmarshal(unmarshal)
}

// NewRational parses raw string and returns Rational.
func NewRational(v string) (Rational, error) {
	r := Rational{}

	v = strings.ToLower(v)
	v = strings.ReplaceAll(v, " ", "-")
	v = strings.ReplaceAll(v, "_", "-")
	v = strings.ReplaceAll(v, "(", "")
	v = strings.ReplaceAll(v, ")", "")

	if v == "0" {
		return Rational0, nil
	}

	var ss []string

	if strings.ContainsRune(v, rune(sep)) {
		ss = strings.Split(v, string(sep))
	} else if strings.ContainsRune(v, rune(sepAlt)) {
		ss = strings.Split(v, string(sepAlt))
	} else if strings.ContainsRune(v, rune(sepAlt2)) {
		ss = strings.Split(v, string(sepAlt2))
	} else { // no separator
		v, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return r, fmt.Errorf("%s: %w", err, ErrParseNum)
		}

		r.Num = v
		r.Den = 1

		return r, nil
	}

	if len(ss) != 2 {
		return r, ErrWrongPartsNumber
	}

	numS := ss[0]
	denS := ss[1]

	num, err := strconv.ParseUint(numS, 10, 64)
	if err != nil {
		return r, fmt.Errorf("%s: %w", err, ErrParseNum)
	}

	r.Num = num

	den, err := strconv.ParseUint(denS, 10, 64)
	if err != nil {
		return r, fmt.Errorf("%s: %w", err, ErrParseDen)
	}

	r.Den = den

	return r, nil
}
