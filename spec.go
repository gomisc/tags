package tags

import (
	"reflect"
	"strings"

	"git.eth4.dev/golibs/errors"
)

const (
	// ErrWrongSpec - ошибка спецификации
	ErrWrongSpec = errors.Const("Wrong tags specifiacation")
	// ErrWrongInterface - ошибка переданного параметра в процессинг
	ErrWrongInterface = errors.Const("Wrong struct interface")
)

// Spec - спецификация описания тэгов и их процессинга
type Spec interface {
	Apply(in interface{}) error
}

type tagSpec struct {
	tags  map[string]TagProcessor
	order []string
}

// ParseSpec - парсит и возвращает спецификацию тэгов
func ParseSpec(in interface{}) (Spec, error) {
	v := reflect.ValueOf(in)
	t := v.Type()

	spec := &tagSpec{
		tags:  make(map[string]TagProcessor),
		order: []string{},
	}

	switch {
	case t.Kind() == reflect.Ptr && v.CanInterface():
		for i := 0; i < t.Elem().NumField(); i++ {
			f := t.Elem().Field(i)
			spec.order = append(spec.order, f.Name)
		}

		return spec, nil
	case t.Kind() == reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			name := strings.ToLower(f.Name)
			spec.order = append(spec.order, name)

			fv := reflect.Indirect(v.Field(i))
			tp, ok := fv.Interface().(TagProcessor)

			if ok {
				spec.tags[name] = tp
			}
		}

		return spec, nil
	default:
		return nil, ErrWrongSpec
	}
}

// Apply - осуществляет процессинг тэгов и применяет их к переданному объекту
func (s *tagSpec) Apply(in interface{}) error {
	sv, err := checkStructValue(in)
	if err != nil {
		return errors.Wrap(err, "check input interface")
	}

	for pair := range recurseFieldsValueIterator(*sv, s.order) {
		for i := range s.order {
			name := s.order[i]
			proc := s.tags[name]
			tv := pair.TagsValues()

			if tagVal, ok := tv[name]; ok {
				if err = proc.Process(pair.FieldValue(), tagVal); err != nil {
					return errors.Ctx().
						Str("field", pair.Name()).
						Any("value", in).
						Just(err)
				}
			}
		}
	}

	return nil
}

func recurseFieldsValueIterator(sv reflect.Value, tags []string) <-chan ValuePair {
	vchan := make(chan ValuePair)

	go func() {
		walkStructFields(vchan, tags, sv)
		close(vchan)
	}()

	return vchan
}

func walkStructFields(res chan ValuePair, tags []string, sv reflect.Value) {
	st := sv.Type()

	for i := 0; i < sv.NumField(); i++ {
		sf := sv.Field(i)

		if !sf.CanSet() {
			continue
		}

		for sf.Kind() == reflect.Ptr {
			if sf.IsNil() {
				if sf.Type().Elem().Kind() != reflect.Struct {
					break
				}

				sf.Set(reflect.New(sf.Type().Elem()))
			}

			sf = sf.Elem()
		}

		if pair := newValuePair(sf, st.Field(i), tags...); pair != nil {
			res <- pair
		}

		if sf.Kind() == reflect.Struct {
			walkStructFields(res, tags, sf)
		}
	}
}

func checkStructValue(in interface{}) (*reflect.Value, error) {
	v := reflect.ValueOf(in)

	if v.Kind() != reflect.Ptr {
		return nil, ErrWrongInterface
	}

	s := v.Elem()
	if s.Kind() != reflect.Struct {
		return nil, ErrWrongInterface
	}

	return &s, nil
}
