package tags

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DirectSetter - стандартный сеттер значений по значению тэга
func DirectSetter() TagProcessor {
	return TagProcFunc(convertFromString)
}

// SetterFromEnv - стандартный сеттер значений из переменных окружения
func SetterFromEnv() TagProcessor {
	return TagProcFunc(func(f reflect.Value, key string) error {
		if val := os.Getenv(key); val != "" {
			return convertFromString(f, val)
		}

		return nil
	})
}

func convertFromString(field reflect.Value, valueStr string) error {
	typ := field.Type()

	switch typ.Kind() {
	case reflect.String:
		field.SetString(valueStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)

		if field.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" {
			var d time.Duration

			d, err = time.ParseDuration(valueStr)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(valueStr, 0, typ.Bits())
		}

		if err != nil {
			return err
		}

		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(valueStr, 0, typ.Bits())
		if err != nil {
			return err
		}

		field.SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(valueStr)
		if err != nil {
			return err
		}

		field.SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(valueStr, typ.Bits())
		if err != nil {
			return err
		}

		field.SetFloat(val)
	case reflect.Slice:
		sl := reflect.MakeSlice(typ, 0, 0)
		if typ.Elem().Kind() == reflect.Uint8 {
			sl = reflect.ValueOf([]byte(valueStr))
		} else if strings.TrimSpace(valueStr) != "" {
			values := strings.Split(valueStr, ",")
			sl = reflect.MakeSlice(typ, len(values), len(values))

			for i, val := range values {
				err := convertFromString(sl.Index(i), val)
				if err != nil {
					return err
				}
			}
		}

		field.Set(sl)
	case reflect.Map:
		mp := reflect.MakeMap(typ)

		if strings.TrimSpace(valueStr) != "" {
			pairs := strings.Split(valueStr, ",")
			for _, pair := range pairs {
				kvpair := strings.Split(pair, ":")
				if len(kvpair) != 2 {
					return fmt.Errorf("invalid map item: %q", pair)
				}

				k := reflect.New(typ.Key()).Elem()

				err := convertFromString(k, kvpair[0])
				if err != nil {
					return err
				}

				v := reflect.New(typ.Elem()).Elem()

				err = convertFromString(k, kvpair[1])
				if err != nil {
					return err
				}

				mp.SetMapIndex(k, v)
			}
		}

		field.Set(mp)
	default:
		return nil
	}

	return nil
}
