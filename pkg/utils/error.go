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

	return fmt.Errorf(msg+" %w", err)
}
