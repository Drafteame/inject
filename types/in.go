package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Drafteame/inject/utils"
)

const (
	tag            = "inject"
	nameOption     = "name"
	optionalOption = "optional"
)

type Container interface {
	Get(name Symbol) (any, error)
}

// In is a struct that should be embedded to other struct to denote that is a valid input for an invoker function and
// his fields should be filled from the dependency container threes.
type In struct{}

// injectInField  is the configuration that each In struct fields should follow to be filled.
type injectInField struct {
	fieldName  string
	injectName Symbol
	optional   bool
	container  Container
}

func BuildIn(cont Container, in reflect.Value) error {
	if !utils.EmbedsType(in.Type(), reflect.TypeOf(In{})) {
		return fmt.Errorf("inject: struct doesn't embed `inject.In` struct")
	}

	itype := in.Type()
	if itype.Kind() == reflect.Ptr {
		itype = itype.Elem()
	}

	nfields := itype.NumField()
	injectFields := make([]injectInField, 0)

	for i := 0; i < nfields; i++ {
		if itype.Field(i).Anonymous {
			continue
		}

		injectField, err := buildInjectInField(itype.Field(i))
		if err != nil {
			return err
		}

		injectField.container = cont

		injectFields = append(injectFields, injectField)
	}

	return fillInStruct(cont, in, injectFields)
}

// fillInStruct It iterates over the `conf` array. For each element in the array, it calls a function based on the value
// of `inject.source`. The functions called are either `fillInStructFromDeps` or `fillInStructFromNamedDeps`. Both
// functions return an error if something goes wrong, and this error is returned by the caller (`fillInStruct`). If no
// errors occur, then nil is returned.
func fillInStruct(cont Container, in reflect.Value, conf []injectInField) error {
	for _, inject := range conf {
		if err := fillStructFieldFromBuilder(cont, in, inject); err != nil {
			return err
		}
	}

	return nil
}

// fillStructFieldFromBuilder It checks if the dependency exists. If it doesn't exist, it returns an error. It builds the
// dependency using `builder`. It sets the field of the struct with name `conf.fieldName` to be equal to `out`. Returns
// nil (no error).
func fillStructFieldFromBuilder(cont Container, in reflect.Value, conf injectInField) error {
	val, err := cont.Get(conf.injectName)
	if err != nil {
		if conf.optional {
			return nil
		}

		return err
	}

	invalue := in

	if invalue.Kind() == reflect.Ptr {
		invalue = invalue.Elem()
	}

	field := invalue.FieldByName(conf.fieldName)

	field.Set(reflect.ValueOf(val))
	return nil
}

// buildInjectInField We get the tags of the field. If there is a `name` tag, we set the source to `sourceNamedDeps`. If
// there is an `optional` tag, we set it to true. We create a new injectInField struct and return it.
func buildInjectInField(field reflect.StructField) (injectInField, error) {
	ftags := getFieldTags(field)
	injectName := ""
	optional := false

	val, ok := ftags[nameOption]

	if !ok {
		return injectInField{}, fmt.Errorf("inject: missing name tag of inject dependency on field `%s`", field.Name)
	}

	injectName = val

	if _, ok := ftags[optionalOption]; ok {
		optional = true
	}

	inject := injectInField{
		fieldName:  field.Name,
		injectName: Symbol(injectName),
		optional:   optional,
	}

	return inject, nil
}

// getFieldTags It gets the tag value from the field. If there is no tag, it returns an empty map. It splits the tag by
// commas and trims spaces from each part of the split result. It creates a map to store all tags and their values (if
// any). For each part of the split result:
//  1. It splits again by equal sign (`=`) and trims spaces from each part of this second split result too; if there is
//     no equal sign, it will be treated as an empty string for that side of the split operation; so `"a=b"` will be
//     splitted into `["a", "b"]`, but `"a="` will be splitted into `["a", ""]`.
//  2. The first element in this second split operation is considered to be a key for our map; if it's an empty string,
//     we skip this iteration because we don't want to add keys with empty names to our map; otherwise, we add it as a
//     key in our map with its value being either.
func getFieldTags(field reflect.StructField) map[string]string {
	value, exists := field.Tag.Lookup(tag)
	if !exists {
		return map[string]string{}
	}

	tags := strings.Split(value, ",")

	tagValues := make(map[string]string)

	for _, tagValue := range tags {
		tagValue = strings.TrimSpace(tagValue)
		aux := strings.Split(tagValue, "=")

		tagOptionName := strings.TrimSpace(aux[0])

		if tagOptionName == "" {
			continue
		}

		if len(aux) > 1 {
			tagValues[tagOptionName] = strings.TrimSpace(aux[1])
		} else {
			tagValues[tagOptionName] = ""
		}
	}

	return tagValues
}
