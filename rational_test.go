package rational_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-x-pkg/rational"
)

func TestRationalUnmarshal(t *testing.T) {
	tests := []struct {
		raw    []byte
		e      *rational.Rational
		eIsNil bool
		eErr   error
	}{
		{
			raw: []byte(`{"rational": "24000/1001"}`),
			e:   &rational.Rational{24000, 1001},
		},
		{
			raw: []byte(`{"rational": "1920x1080"}`),
			e:   &rational.Rational{1920, 1080},
		},
		{
			raw: []byte(`{"rational": "16:9"}`),
			e:   &rational.Rational{16, 9},
		},
		{
			raw:    []byte(`{"rational": "0"}`),
			e:      &rational.Rational{0, 0},
			eIsNil: true,
		},
		{
			raw: []byte(`{"rational": "18647"}`),
			e:   &rational.Rational{18647, 1},
		},
		{
			raw:  []byte(`{"rational": "0:0:0"}`),
			eErr: rational.ErrWrongPartsNumber,
		},
		{
			raw:  []byte(`{"rational": "0~~~~xddqwdq"}`),
			eErr: rational.ErrParseNum,
		},
		{
			raw:  []byte(`{"rational": "1xddqwdq"}`),
			eErr: rational.ErrParseDen,
		},
		{
			raw: []byte(`{}`),
			e:   nil,
		},
	}

	for i, tt := range tests {
		func() {
			data := struct {
				Rational *rational.Rational `json:"rational" yaml:"rational" bson:"rational"`
			}{}

			err := json.Unmarshal(tt.raw, &data)
			if err != nil {
				if tt.eErr == nil {
					t.Errorf("#%d: err: %s", i, err)
					return
				} else {
					if !errors.Is(err, tt.eErr) {
						t.Errorf("#%d: unexpected error: got %#v expected %#v", i, err, tt.eErr)
						return
					}
				}
			}

			if data.Rational == nil && tt.e != nil {
				t.Errorf("#%d: bad rational: expected nil", i)
				return
			}

			if data.Rational != nil && tt.e != nil && *data.Rational != *tt.e {
				t.Errorf("#%d: bad rational: got %#v xpected %#v", i, data.Rational, tt.e)
			}

			if _, err := json.Marshal(data); err != nil {
				t.Errorf("#%d: err marshal JSON: %s", i, err)
			}

			if data.Rational != nil {
				// should not panic
				data.Rational.Float64()
				data.Rational.Reverse().Float64()
				data.Rational.Percent()

				if !data.Rational.IsNil() && tt.eIsNil {
					t.Errorf("#%d: rational exp expeted to be nil", i)
				}
			}
		}()
	}
}
