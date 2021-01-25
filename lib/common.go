package lib

import (
	"bytes"
	"fmt"
	"reflect"
)

// YamlToJsonEncode 反射方式，将yaml数据转成json数据格式
func YamlToJsonEncode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return YamlToJsonEncode(buf, v.Elem())
	case reflect.Bool:
		fmt.Fprintf(buf, "%t", v.Bool())
	case reflect.String:
		str := v.String()
		fmt.Fprintf(buf, "%q", str)
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%v", v.Float())
	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := YamlToJsonEncode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case reflect.Map:
		buf.WriteByte('{')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := YamlToJsonEncode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := YamlToJsonEncode(buf, v.MapIndex(key)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

// 三元表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
