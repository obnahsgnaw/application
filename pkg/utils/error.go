package utils

import (
	"errors"
	"fmt"
)

// NewWrappedError new a wrapped error
func NewWrappedError(msg string, err error) error {
	if err == nil {
		return errors.New(msg)
	}

	return fmt.Errorf(msg+": %w", err)
}

func TitledError(title, msg string, err error) error {
	titleMsg := ToStr(title, ", ", msg)
	if err != nil {
		return NewWrappedError(titleMsg, err)
	}
	return errors.New(titleMsg)
}
