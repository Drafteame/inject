package inject

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Drafteame/inject/dependency"
	"github.com/Drafteame/inject/utils"
)

const (
	tag             = "inject"
	nameOption      = "name"
	optionalOption  = "optional"
	sourceDeps      = "deps"
	sourceNamedDeps = "named-deps"
)

// In is a struct that should be embedded to other struct to denote that is a valid input for an invoker function and
// his fields should be filled from the dependency container threes.
type In struct{}

// IsIn validate if the given element embeds the struct `inject.In`.
func IsIn(i any) bool {
	return utils.EmbedsType(i, reflect.TypeOf(In{}))
}

// injectInField  is the configuration that each In struct fields should follow to be filled.
type injectInField struct {
	fieldName       string
	source          string
	injectName      string
	injectType      reflect.Type
	optional        bool
	sharedContainer *shared
}

func buildInStruct(cont *container, in reflect.Value) error {
	if !IsIn(in.Type()) {
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

		injectField := buildInjectInField(itype.Field(i))
		injectField.sharedContainer = cont.shared

		injectFields = append(injectFields, injectField)
	}

	return fillInStruct(cont, in, injectFields)
}

// fillInStruct It iterates over the `conf` array. For each element in the array, it calls a function based on the value
// of `inject.source`. The functions called are either `fillInStructFromDeps` or `fillInStructFromNamedDeps`. Both
// functions return an error if something goes wrong, and this error is returned by the caller (`fillInStruct`). If no
// errors occur, then nil is returned.
func fillInStruct(cont *container, in reflect.Value, conf []injectInField) error {
	for _, inject := range conf {
		switch inject.source {
		default:
			if err := fillInStructFromDeps(cont, in, inject); err != nil {
				return err
			}
		case sourceNamedDeps:
			if err := fillInStructFromNamedDeps(cont, in, inject); err != nil {
				return err
			}
		}
	}

	return nil
}

// fillInStructFromDeps It gets the builder from the container by type. It calls `fillStructFromBuilder` with that
// builder and a boolean indicating whether it was found or not. `fillStructFromBuilder` will call `builder()` to get an
// instance of the dependency, and then set that value in the struct field using reflection. If there is no builder for
// that dependency, it will return an error.
func fillInStructFromDeps(cont *container, in reflect.Value, conf injectInField) error {
	builder, ok := cont.deps[conf.injectType]
	return fillStructFromBuilder(builder, ok, in, conf)
}

// fillInStructFromNamedDeps It gets the builder from the container by name. It calls `fillStructFromBuilder` with that
// builder and a boolean indicating whether it was found or not. `fillStructFromBuilder` will call `builder()` to get an
// instance of the dependency, and then set that value in the struct field using reflection. If there is no builder for
// that dependency, it will return an error.
func fillInStructFromNamedDeps(cont *container, in reflect.Value, conf injectInField) error {
	builder, ok := cont.depsByName[conf.injectName]
	return fillStructFromBuilder(builder, ok, in, conf)
}

// fillStructFromBuilder It checks if the dependency exists. If it doesn't exist, it returns an error. It builds the
// dependency using `builder`. It sets the field of the struct with name `conf.fieldName` to be equal to `out`. Returns
// nil (no error).
func fillStructFromBuilder(builder dependency.Builder, exist bool, in reflect.Value, conf injectInField) error {
	if !exist {
		if conf.optional {
			return nil
		}

		nameErrMsg := ""
		if conf.injectName != "" {
			nameErrMsg = fmt.Sprintf(" and name %s", conf.injectName)
		}

		return fmt.Errorf("inject: can't provide dependency of type `%v`%s to In receiver", conf.injectType, nameErrMsg)
	}

	out, err := builder.WithSharedContainer(conf.sharedContainer).Build()
	if err != nil {
		return err
	}

	invalue := in

	if invalue.Kind() == reflect.Ptr {
		invalue = invalue.Elem()
	}

	field := invalue.FieldByName(conf.fieldName)

	field.Set(reflect.ValueOf(out))
	return nil
}

// buildInjectInField We get the tags of the field. If there is a `name` tag, we set the source to `sourceNamedDeps`. If
// there is an `optional` tag, we set it to true. We create a new injectInField struct and return it.
func buildInjectInField(field reflect.StructField) injectInField {
	ftags := getFieldTags(field)
	injectType := field.Type
	injectName := ""
	source := sourceDeps
	optional := false

	if val, ok := ftags[nameOption]; ok {
		source = sourceNamedDeps
		injectName = val
	}

	if _, ok := ftags[optionalOption]; ok {
		optional = true
	}

	inject := injectInField{
		fieldName:  field.Name,
		source:     source,
		injectName: injectName,
		injectType: injectType,
		optional:   optional,
	}

	return inject
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
