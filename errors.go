package rational

import "errors"

var (
	// ErrWrongPartsNumber indicates parsing error
	// because of wrong parts number in raw string.
	ErrWrongPartsNumber error = errors.New("got wrong string parts number")

	// ErrParseNum indicates parsing numerator error.
	ErrParseNum error = errors.New("failed parse numerator")

	// ErrParseDen indicates parsing denumerator error.
	ErrParseDen error = errors.New("failed parse denumerator")
)
