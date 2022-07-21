package errors

import (
	"errors"
	"fmt"
)

func New(text string) error {
	return errors.New(text)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err, target error) bool {
	return errors.As(err, target)
}

func Wrap(err error, label ...map[string]any) error {
	if err == nil {
		return nil
	}

	if label == nil {
		return err
	}

	return fmt.Errorf("%w %v", err, label[0])
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
