package structs

import (
	"reflect"
)

// MapField retrieves struct field as map[name/tag]*Field from <pointer>, and returns the map.
//
// The parameter <pointer> should be type of struct/*struct.
//
// The parameter <priority> specifies the priority tag array for retrieving from high to low.
//
// The parameter <recursive> specifies whether retrieving the struct field recursively.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func MapField(pointer interface{}, priority []string, recursive bool) map[string]*TagField {
	fieldMap := make(map[string]*TagField)
	fields := ([]*Field)(nil)
	if v, ok := pointer.(reflect.Value); ok {
		fields = Fields(v.Interface())
	} else {
		fields = Fields(pointer)
	}
	tag := ""
	name := ""
	for _, field := range fields {
		name = field.Name()
		// Only retrieve exported attributes.
		if name[0] < byte('A') || name[0] > byte('Z') {
			continue
		}
		fieldMap[name] = &TagField{
			Field: field,
			Tag:   tag,
		}
		tag = ""
		for _, p := range priority {
			tag = field.Tag(p)
			if tag != "" {
				break
			}
		}
		if tag != "" {
			fieldMap[tag] = &TagField{
				Field: field,
				Tag:   tag,
			}
		}
		if recursive {
			rv := reflect.ValueOf(field.Value())
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			if kind == reflect.Struct {
				for k, v := range MapField(rv, priority, true) {
					if _, ok := fieldMap[k]; !ok {
						fieldMap[k] = v
					}
				}
			}
		}
	}
	return fieldMap
}
