package gke

import (
	"encoding/json"
	"reflect"
)

// Convert between two types by converting to/from JSON. Intended to switch
// between multiple API versions, as they are strict supersets of one another.
// item and out are pointers to structs
func Convert(item, out interface{}) error {
	bytes, err := json.Marshal(item)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}

	// Converting between maps and structs only occurs when autogenerated resources convert the result
	// of an HTTP request. Those results do not contain omitted fields, so no need to set them.
	if _, ok := item.(map[string]interface{}); !ok {
		setOmittedFields(item, out)
	}

	return nil
}

// When converting to a map, we can't use setOmittedFields because FieldByName
// fails. Luckily, we don't use the omitted fields anymore with generated
// resources, and this function is used to bridge from handwritten -> generated.
// Since this is a known type, we can create it inline instead of needing to
// pass an object in.
func ConvertToMap(item interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	bytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func setOmittedFields(item, out interface{}) {
	// Both inputs must be pointers, see https://blog.golang.org/laws-of-reflection:
	// "To modify a reflection object, the value must be settable."
	iVal := reflect.ValueOf(item).Elem()
	oVal := reflect.ValueOf(out).Elem()

	// Loop through all the fields of the struct to look for omitted fields and nested fields
	for i := 0; i < iVal.NumField(); i++ {
		iField := iVal.Field(i)
		if isEmptyValue(iField) {
			continue
		}

		fieldInfo := iVal.Type().Field(i)
		oField := oVal.FieldByName(fieldInfo.Name)

		// Only look at fields that exist in the output struct
		if !oField.IsValid() {
			continue
		}

		// If the field contains a 'json:"="' tag, then it was omitted from the Marshal/Unmarshal
		// call and needs to be added back in.
		if fieldInfo.Tag.Get("json") == "-" {
			oField.Set(iField)
		}

		// If this field is a struct, *struct, []struct, or []*struct, recurse.
		if iField.Kind() == reflect.Struct {
			setOmittedFields(iField.Addr().Interface(), oField.Addr().Interface())
		}
		if iField.Kind() == reflect.Ptr && iField.Type().Elem().Kind() == reflect.Struct {
			setOmittedFields(iField.Interface(), oField.Interface())
		}
		if iField.Kind() == reflect.Slice && iField.Type().Elem().Kind() == reflect.Struct {
			for j := 0; j < iField.Len(); j++ {
				setOmittedFields(iField.Index(j).Addr().Interface(), oField.Index(j).Addr().Interface())
			}
		}
		if iField.Kind() == reflect.Slice && iField.Type().Elem().Kind() == reflect.Ptr &&
			iField.Type().Elem().Elem().Kind() == reflect.Struct {
			for j := 0; j < iField.Len(); j++ {
				setOmittedFields(iField.Index(j).Interface(), oField.Index(j).Interface())
			}
		}
	}
}
