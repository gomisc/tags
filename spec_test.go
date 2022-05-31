package tags

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTagsProcessorSetter(t *testing.T) {
	type customSpec struct {
		Default TagProcessor
		Env     TagProcessor
		Valid   TagProcessor
	}

	type S2 struct {
		Float   float64       `env:"TEST_FLOAT"`
		Timeout time.Duration `env:"TEST_DURATION"`
	}

	type S1 struct {
		S2      *S2
		String  string `env:"TEST_STRING" default:"something"`
		Integer int    `env:"TEST_INTEGER"`
	}

	spec, err := ParseSpec(customSpec{
		Default: DirectSetter(),
		Env:     SetterFromEnv(),
		Valid: TagProcFunc(func(f reflect.Value, tv string) error {
			return nil
		}),
	})
	assert.NoError(t, err)

	os.Setenv("TEST_FLOAT", "1.25")
	os.Setenv("TEST_DURATION", "25s")
	os.Setenv("TEST_INTEGER", "100")
	os.Setenv("TEST_STRING", "something string here")

	s1 := &S1{}
	assert.NoError(t, spec.Apply(s1))
	assert.Equal(t, s1, &S1{
		String:  "something string here",
		Integer: 100,
		S2: &S2{
			Float:   1.25,
			Timeout: time.Second * 25,
		},
	})
}

func TestTagsProcessorGetter(t *testing.T) {
	type customSpec struct {
		Form TagProcessor
	}

	t.Run("getter to map", func(t *testing.T) {
		type str struct {
			String string            `form:"string_param"`
			Int    int               `form:"int_param"`
			Uint   uint              `form:"uint_param"`
			Slice  []string          `form:"slice_param"`
			Map    map[string]string `form:"map_param"`
		}

		m := make(map[string]interface{})

		spec, err := ParseSpec(customSpec{
			Form: DirectGetter(func(key string, value interface{}) {
				m[key] = value
			}),
		})
		assert.NoError(t, err, "parsing spec")
		assert.NoError(t, spec.Apply(&str{
			String: "some string",
			Int:    10,
			Uint:   20,
			Slice:  []string{"item1"},
			Map:    map[string]string{"key1": "value1"},
		}))
		assert.EqualValuesf(t,
			map[string]interface{}{
				"string_param": "some string",
				"int_param":    10,
				"uint_param":   uint(20),
				"slice_param":  []string{"item1"},
				"map_param":    map[string]string{"key1": "value1"},
			},
			m,
			"comapare results",
		)
	})
}
