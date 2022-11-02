package tags

import (
	"reflect"

	"git.eth4.dev/golibs/errors"
)

// TagProcessor - интерфейс процессора тэгов
type TagProcessor interface {
	Process(f reflect.Value, tv string) error
}

// TagProcFunc - процессор тэга
type TagProcFunc func(f reflect.Value, tv string) error

// Process - выполняет процессинг тэга для поля
func (p TagProcFunc) Process(f reflect.Value, tv string) error {
	return p(f, tv)
}

// ValuePair - значение тэгов поля структуры
type ValuePair interface {
	Name() string
	TagsValues() map[string]string
	Field() reflect.StructField
	FieldValue() reflect.Value
}

type valuePair struct {
	tv map[string]string
	fv reflect.Value
	f  reflect.StructField
}

// nolint: gocritic
func newValuePair(v reflect.Value, f reflect.StructField, tags ...string) *valuePair {
	var pair *valuePair

	for i := range tags {
		tn := tags[i]

		if tv, ok := f.Tag.Lookup(tn); ok {
			if pair == nil {
				pair = &valuePair{
					f:  f,
					fv: v,
					tv: make(map[string]string),
				}
			}

			pair.tv[tn] = tv
		}
	}

	return pair
}

// Name возвращает имя поля значения
func (tv *valuePair) Name() string {
	return tv.f.Name
}

// TagsValues - возвращает значения тэгов поля
func (tv *valuePair) TagsValues() map[string]string {
	return tv.tv
}

// Field - возвращает интерфейс поля структуры
func (tv *valuePair) Field() reflect.StructField {
	return tv.f
}

// FieldValue - возвращает интерфейс значения поля структуры
func (tv *valuePair) FieldValue() reflect.Value {
	return tv.fv
}

// FieldTagsIterator - итератор по полям структуры со значениями указанных тэгов
func FieldTagsIterator(in interface{}, tags ...string) (<-chan ValuePair, error) {
	sv, err := checkStructValue(in)
	if err != nil {
		return nil, errors.Wrap(err, "check input interface")
	}

	return recurseFieldsValueIterator(*sv, tags), nil
}
