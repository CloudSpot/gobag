package validate

import (
	"errors"
	"fmt"
)

var ErrRequired = errors.New("missing or invalid")

type Validator interface {
	Validate() error
}

type FieldValidator func() error

func Compose(validators ...FieldValidator) error {
	var err error
	for _, v := range validators {
		if err = v(); err != nil {
			break
		}
	}

	return err
}

func MakeFieldValidator(key string, fn func() error) FieldValidator {
	return FieldValidator(func() error {
		if err := fn(); err != nil {
			return fmt.Errorf("%s: %v", key, err)
		}

		return nil
	})
}

func RequiredString(key string, v string) FieldValidator {
	return MakeFieldValidator(key, func() error {
		if v == "" {
			return ErrRequired
		}

		return nil
	})
}

func IsTrue(key string, v bool, msg string) FieldValidator {
	return MakeFieldValidator(key, func() error {
		if !v {
			return errors.New(msg)
		}

		return nil
	})
}

func Elements(key string, els []interface{}) FieldValidator {
	return MakeFieldValidator(key, func() error {
		for i, el := range els {
			if ev, ok := el.(Validator); ok {
				if err := ev.Validate(); err != nil {
					return fmt.Errorf("#%d: %v", i, err)
				}
			}
		}

		return nil
	})
}
