package util

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
)

func SetFromEnv(value interface{}, prefix string) error {
	return setFromEnv(reflect.ValueOf(value), prefix)
}

func setFromEnv(v reflect.Value, name string) (err error) {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	var parse func(v reflect.Value, prefix string) (err error)
	switch v.Kind() {
	case reflect.Struct:
		parse = setStructFromEnv
	case reflect.Map:
		parse = setMapFromEnv
	case reflect.Array:
		parse = setArrayFromEnv
	case reflect.Slice:
		parse = setSliceFromEnv
	default:
		parse = setPrimitiveFromEnv

	}

	if e := parse(v, name); e != nil {
		err = fmt.Errorf("%s (%s) - %v", name, typeString(v.Type()), e)
	}

	return
}

// NOTE: Assumes value is of kind reflect.Struct
func setStructFromEnv(v reflect.Value, prefix string) (err error) {
	for i := 0; i < v.NumField(); i++ {
		subName := fmt.Sprintf("%s.%s", prefix, v.Type().Field(i).Name)
		if e := setFromEnv(v.Field(i), subName); e != nil {
			err = multierror.Append(err, e)
		}
	}
	return
}

// NOTE: Assumes value is of kind reflect.Map
func setMapFromEnv(v reflect.Value, prefix string) (err error) {
	var regex *regexp.Regexp

	// Build regex to match environment variable names
	regexString := fmt.Sprintf("^(%s\\.[\\w\\.]+)=.*", regexp.QuoteMeta(prefix))
	if regex, err = regexp.Compile(regexString); err != nil {
		return
	}

	// Search for env variables with mathing name and apply it
	for _, envExpr := range os.Environ() {
		if env := regex.FindSubmatch([]byte(envExpr)); env != nil {
			key := strings.TrimPrefix(string(env[1]), prefix+".")
			if e := setMapIndexFromEnv(v, prefix, key); e != nil {
				envValue := os.Getenv(prefix + "." + key)
				err = multierror.Append(err, fmt.Errorf("[%s]=%s - %v", key, envValue, e))
			}
		}
	}

	return
}

// NOTE: Assumes v is of kind reflect.Map
func setMapIndexFromEnv(v reflect.Value, name string, key string) (err error) {
	canonicalName := fmt.Sprintf("%s.%s", name, key)
	mapValue := reflect.New(v.Type().Elem()).Elem() // Empty map value with matching type
	mapKey := reflect.New(v.Type().Key()).Elem()    // Empty map key with matching type

	if e := stringToValue(mapKey, key); e != nil {
		err = multierror.Append(err, fmt.Errorf("map key - %v", e))
	}
	if e := setPrimitiveFromEnv(mapValue, canonicalName); e != nil {
		err = multierror.Append(err, fmt.Errorf("map value - %v", e))
	}

	if err != nil {
		return
	}

	v.SetMapIndex(mapKey, mapValue)
	return
}

// NOTE: Assumes value is of kind reflect.Array
func setArrayFromEnv(value reflect.Value, name string) (err error) {
	if env, ok := os.LookupEnv(name); ok {
		unmarshaled := make([]interface{}, 0)
		if err = json.Unmarshal([]byte(env), &unmarshaled); err == nil {
			length := int(math.Min(float64(len(unmarshaled)), float64(value.Len())))
			for i := 0; i < length; i++ {
				entry := reflect.ValueOf(unmarshaled[i])
				indexValue := value.Index(i)
				if jsonBytes, e := json.Marshal(entry.Interface()); e != nil {
					err = multierror.Append(err, fmt.Errorf("[%d] - %v", i, e))
					continue
				} else if e = stringToValue(indexValue, string(jsonBytes)); e != nil {
					err = multierror.Append(err, fmt.Errorf("[%d]=%s - %v", i, jsonBytes, e))
					continue
				}
			}
		}
	}

	return
}

// NOTE: Assumes value is of kind reflect.Slice
func setSliceFromEnv(value reflect.Value, name string) (err error) {
	if env, ok := os.LookupEnv(name); ok {
		unmarshaled := make([]interface{}, 0)
		if err = json.Unmarshal([]byte(env), &unmarshaled); err == nil {
			for i, entry := range unmarshaled {
				entry := reflect.ValueOf(entry)
				indexValue := reflect.New(value.Type().Elem()).Elem()

				if jsonBytes, e := json.Marshal(entry.Interface()); e != nil {
					err = multierror.Append(err, fmt.Errorf("[%d] - %v", i, e))
					continue
				} else if e = stringToValue(indexValue, string(jsonBytes)); e != nil {
					err = multierror.Append(err, fmt.Errorf("[%d]=%s - %v", i, jsonBytes, e))
					continue
				}

				if value.Len() > i {
					value.Index(i).Set(indexValue)
				} else {
					value.Set(reflect.Append(value, indexValue))
				}
			}
		}
	}
	return
}

func setPrimitiveFromEnv(v reflect.Value, name string) (err error) {
	if env, ok := os.LookupEnv(name); ok {
		err = stringToValue(v, env)
	}
	return
}

func stringToValue(v reflect.Value, env string) (err error) {
	var value interface{}

	switch v.Kind() {
	case reflect.String:
		value = env
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err = strconv.ParseInt(env, 10, 64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err = strconv.ParseUint(env, 10, 64)
	case reflect.Bool:
		value = strings.ToLower(env) == "true"
	case reflect.Float32, reflect.Float64:
		value, err = strconv.ParseFloat(env, 64)
	case reflect.Complex64, reflect.Complex128:
		value, err = strconv.ParseComplex(env, 128)
	}

	if err == nil && value != nil {
		reflected := reflect.ValueOf(value)
		if reflected.CanConvert(v.Type()) {
			v.Set(reflected.Convert(v.Type()))
		}
	}
	return

}

func typeString(v reflect.Type) string {
	switch v.Kind() {

	case reflect.Pointer:
		return fmt.Sprintf("*%s", typeString((v.Elem())))

	case reflect.Array:
		return fmt.Sprintf("[%d]%s", v.Len(), typeString((v.Elem())))

	case reflect.Slice:
		return fmt.Sprintf("[]%s", typeString((v.Elem())))

	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", typeString((v.Key())), typeString((v.Elem())))

		// Primitives
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return v.Kind().String()

	default:
		if v.Name() != "" {
			return v.Name()
		} else {
			return v.Kind().String()
		}
	}
}
