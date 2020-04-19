// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs

import (
	"reflect"
)

// Field is alias of structs.Field.
type TagField struct {
	*Field
	// Retrieved tag name. There might be more than one tags in the field,
	// but only one can be retrieved according to calling function rules.
	Tag string
}

// TagFields retrieves struct tags as []*Field from <pointer>, and returns it.
//
// The parameter <pointer> should be type of struct/*struct.
//
// The parameter <recursive> specifies whether retrieving the struct field recursively.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func TagFields(pointer interface{}, priority []string, recursive bool) []*TagField {
	return doTagFields(pointer, priority, recursive, map[string]struct{}{})
}

// doTagFields retrieves the tag and corresponding attribute name from <pointer>. It also filters repeated
// tag internally.
// The parameter <pointer> should be type of struct/*struct.
func doTagFields(pointer interface{}, priority []string, recursive bool, tagMap map[string]struct{}) []*TagField {
	var fields []*Field
	if v, ok := pointer.(reflect.Value); ok {
		fields = Fields(v.Interface())
	} else {
		rv := reflect.ValueOf(pointer)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		// If pointer is type of **struct and nil, then automatically create a temporary struct,
		// which is used for structs.Fields.
		if kind == reflect.Ptr && (!rv.IsValid() || rv.IsNil()) {
			fields = Fields(reflect.New(rv.Type().Elem()).Elem().Interface())
		} else {
			fields = Fields(pointer)
		}
	}
	tag := ""
	name := ""
	tagFields := make([]*TagField, 0)
	for _, field := range fields {
		name = field.Name()
		// Only retrieve exported attributes.
		if name[0] < byte('A') || name[0] > byte('Z') {
			continue
		}
		tag = ""
		for _, p := range priority {
			tag = field.Tag(p)
			if tag != "" {
				break
			}
		}
		if tag != "" {
			// Filter repeated tag.
			if _, ok := tagMap[tag]; ok {
				continue
			}
			tagFields = append(tagFields, &TagField{
				Field: field,
				Tag:   tag,
			})
		}
		if recursive {
			rv := reflect.ValueOf(field.Value())
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			if kind == reflect.Struct {
				tagFields = append(tagFields, doTagFields(rv, priority, recursive, tagMap)...)
			}
		}
	}
	return tagFields
}

// TagMapName retrieves struct tags as map[tag]attribute from <pointer>, and returns it.
//
// The parameter <pointer> should be type of struct/*struct.
//
// The parameter <recursive> specifies whether retrieving the struct field recursively.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func TagMapName(pointer interface{}, priority []string, recursive bool) map[string]string {
	fields := TagFields(pointer, priority, recursive)
	tagMap := make(map[string]string, len(fields))
	for _, v := range fields {
		tagMap[v.Tag] = v.Name()
	}
	return tagMap
}

// TagMapField retrieves struct tags as map[tag]*Field from <pointer>, and returns it.
//
// The parameter <pointer> should be type of struct/*struct.
//
// The parameter <recursive> specifies whether retrieving the struct field recursively.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func TagMapField(pointer interface{}, priority []string, recursive bool) map[string]*TagField {
	fields := TagFields(pointer, priority, recursive)
	tagMap := make(map[string]*TagField, len(fields))
	for _, v := range fields {
		tagMap[v.Tag] = v
	}
	return tagMap
}
