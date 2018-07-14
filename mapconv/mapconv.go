package mapconv

import (
	"errors"
	"fmt"
	"time"

	"github.com/danielkrainas/gobag/iconv"
)

type Mapper interface {
	Map() map[string]interface{}
}

type Unmapper interface {
	Unmap(m map[string]interface{}) error
}

type MapFunc func() map[string]interface{}

type UnmapFunc func(map[string]interface{}) error

type ValueParser func(m map[string]interface{}) error

var errMissing = errors.New("missing required value")

func Parser(key string, required bool, f func(v interface{}) error) ValueParser {
	return func(m map[string]interface{}) error {
		var err error
		value, ok := m[key]
		if !ok && required {
			err = errMissing
		} else if ok {
			err = f(value)
		}

		if err != nil {
			return fmt.Errorf("%s: %v", key, err)
		}

		return nil
	}
}

func Int(key string, required bool, ref *int) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.Int(field)
		if err == nil {
			*ref = v
		}

		return nil
	})
}

func Int16(key string, required bool, ref *int16) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.Int16(field)
		if err == nil {
			*ref = v
		}

		return nil
	})
}

func Int32(key string, required bool, ref *int32) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.Int32(field)
		if err == nil {
			*ref = v
		}

		return nil
	})
}

func Int64(key string, required bool, ref *int64) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.Int64(field)
		if err == nil {
			*ref = v
		}

		return nil
	})
}

func Float64(key string, required bool, ref *float64) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.Float64(field)
		if err == nil {
			*ref = v
		}

		return nil
	})
}

func Float32(key string, required bool, ref *float32) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.Float32(field)
		if err == nil {
			*ref = v
		}

		return nil
	})
}

func String(key string, required bool, ref *string) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.String(field)
		if err == nil {
			*ref = v
		}

		return err
	})
}

func UnixTime(key string, required bool, ref *time.Time) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.Int64(field)
		if err == nil {
			*ref = time.Unix(v, 0)
		}

		return err
	})
}

func ParseUnmap(key string, required bool, unmap UnmapFunc) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.Map(field)
		if err == nil {
			err = unmap(v)
		}

		return err
	})
}

func ParseUnmapArray(key string, required bool, unmap UnmapFunc) ValueParser {
	return Parser(key, required, func(field interface{}) error {
		v, err := iconv.MapArray(field)
		if err == nil {
			for _, m := range v {
				if err = unmap(m); err != nil {
					break
				}
			}
		}

		return err
	})
}

func Compose(m map[string]interface{}, parsers ...ValueParser) error {
	var err error
	for _, p := range parsers {
		if err = p(m); err != nil {
			break
		}
	}

	return err
}
