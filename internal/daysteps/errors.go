package daysteps

import "errors"

var (
	errInvalidData     = errors.New("data is incorrect")
	errInvalidSteps    = errors.New("steps count is invalid")
	errInvalidDuration = errors.New("duration is invalid")
)
