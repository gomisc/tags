package tags

import (
	"reflect"
)

// DestinationSetter callback устанавливающий значения из объекта по тэгу
type DestinationSetter func(key string, value interface{})

// DirectGetter стандартный геттер значения полей
func DirectGetter(setter DestinationSetter) TagProcessor {
	return TagProcFunc(func(field reflect.Value, tv string) error {
		if field.IsZero() {
			return nil
		}

		typ := field.Type()

		switch typ.Kind() {
		case reflect.String:
			setter(tv, field.Interface().(string))
		case reflect.Int:
			setter(tv, field.Interface().(int))
		case reflect.Int8:
			setter(tv, field.Interface().(int8))
		case reflect.Int16:
			setter(tv, field.Interface().(int16))
		case reflect.Int32:
			setter(tv, field.Interface().(int32))
		case reflect.Int64:
			setter(tv, field.Interface().(int64))
		case reflect.Uint:
			setter(tv, field.Interface().(uint))
		case reflect.Uint8:
			setter(tv, field.Interface().(uint8))
		case reflect.Uint16:
			setter(tv, field.Interface().(uint16))
		case reflect.Uint32:
			setter(tv, field.Interface().(uint32))
		case reflect.Uint64:
			setter(tv, field.Interface().(uint64))
		case reflect.Bool:
			setter(tv, field.Interface().(bool))
		case reflect.Float32:
			setter(tv, field.Interface().(float32))
		case reflect.Float64:
			setter(tv, field.Interface().(float64))
		case reflect.Slice, reflect.Map, reflect.Struct:
			setter(tv, field.Interface())
		case reflect.Ptr:
			setter(tv, field.Addr().Interface())
		default:
			return nil
		}

		return nil
	})
}
