package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func getValueForMember(document *Document, memberType MemberType, memberNames ...string) (interface{}, error) {
	var search map[string]interface{}
	switch memberType {
	case MemberTypeAttribute:
		search = document.Data.Attributes
	case MemberTypeMeta:
		search = document.Data.Meta
	}
	for i, name := range memberNames {
		value, ok := search[name]
		if !ok {
			return "", fmt.Errorf("not ok")
		}
		if i == len(memberNames)-1 {
			return value, nil
		}
		search = search[name].(map[string]interface{})
	}
	return "", fmt.Errorf("memberNames was empty")
}

func Unmarshal(data []byte, v interface{}) error {
	document := NewDocument()
	if err := json.Unmarshal(data, document); err != nil {
		return err
	}

	if err := iterateStruct(document, v, func(value reflect.Value, memberType MemberType, memberNames ...string) error {
		fieldKind := value.Kind()
		// TODO this sets ID for all nexted primary tag fields
		if memberType == MemberTypePrimary {
			if fieldKind != reflect.String {
				return fmt.Errorf("ID must be a string")
			}
			value.SetString(document.Data.ID)
			return nil
		}

		// get raw value
		v, err := getValueForMember(document, memberType, memberNames...)
		if err != nil {
			return nil
		}

		// set raw value
		if err := unmarshal(v, &value); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func RegisterUnmarshaler(t reflect.Type, u unmarshalerFunc) {
	customUnmarshalers[t] = u
}

type unmarshalerFunc = func(interface{}, *reflect.Value)

var customUnmarshalers = make(map[reflect.Type]unmarshalerFunc)

func unmarshal(v interface{}, rv *reflect.Value) error {
	switch rv.Kind() {
	case reflect.String:
		if val, ok := v.(string); ok {
			rv.SetString(val)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, ok := v.(float64); ok {
			// TODO resourceValue.OverflowInt(val)
			rv.SetInt(int64(val))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if val, ok := v.(float64); ok {
			// TODO resourceValue.OverflowInt(val)
			rv.SetUint(uint64(val))
		}
	case reflect.Float32, reflect.Float64:
		if val, ok := v.(float64); ok {
			// TODO resourceValue.OverflowInt(val)
			rv.SetFloat(val)
		}
	default:
		cu, ok := customUnmarshalers[rv.Type()]
		if !ok {
			return fmt.Errorf("Type not supported, must implement custom unmarshaller")
		}
		cu(v, rv)
	}
	return nil
}

func getMember(field reflect.StructField) (MemberType, string, error) {
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		return "", "", fmt.Errorf("tag: %s, not specified", tagKey)
	}
	if tag == "" {
		return "", "", fmt.Errorf("tag: %s, was empty", tagKey)
	}
	tagParts := strings.Split(tag, ",")
	if len(tagParts) != 2 {
		return "", "", fmt.Errorf("tag: %s, was not formatted properly", tagKey)
	}
	memberType, err := NewMemberType(tagParts[0])
	if err != nil {
		return "", "", err
	}
	return memberType, tagParts[1], nil
}
