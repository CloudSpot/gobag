package validate

import (
	"errors"
	"fmt"
)

var ErrRequired = errors.New("missing or invalid")

type Validator func() error

func Compose(validators ...Validator) error {
	var err error
	for _, v := range validators {
		if err = v(); err != nil {
			break
		}
	}

	return err
}

func MakeValidator(key string, fn func() error) Validator {
	return Validator(func() error {
		if err := fn(); err != nil {
			return fmt.Errorf("%s: %v", key, err)
		}

		return nil
	})
}

func RequiredString(key string, v string) Validator {
	return MakeValidator(key, func() error {
		if v == "" {
			return ErrRequired
		}

		return nil
	})
}
