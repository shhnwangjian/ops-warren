package lib

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
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

// If 三元表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// SplitNullString 根据空行切分字符串
func SplitNullString(s string) (slice []string) {
	var buffer bytes.Buffer
	lines := strings.Split(s, "\n")
	for _, v := range lines {
		if len(v) == 0 {
			if len(buffer.String()) > 0 {
				slice = append(slice, buffer.String())
				buffer.Reset()
			}
		} else {
			buffer.WriteString(v)
			buffer.WriteString("\n")
		}
	}
	return
}

// RemoveSuffix 判断dest是否存在suffix结尾，删除
func RemoveSuffix(dest, suffix string) string {
	result := strings.HasSuffix(dest, suffix)
	if result {
		dest = strings.TrimSuffix(dest, suffix)
	}
	return dest
}

// RemovePrefix 判断dest开头是否存在prefix，存在删除prefix
func RemovePrefix(dest, prefix string) string {
	result := strings.HasPrefix(dest, prefix)
	if result {
		dest = strings.TrimPrefix(dest, prefix)
	}
	return dest
}
