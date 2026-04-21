package reflect

import (
	"bytes"
	"reflect"
)

var bytesType = reflect.TypeFor[[]byte]()

// IsBytes 判断值、类型或 reflect.Value 是否表示 []byte 语义。
func IsBytes(v any) bool {
	switch v.(type) {
	case []byte:
		return true
	default:
		var t reflect.Type

		switch x := v.(type) {
		case reflect.Type:
			t = x
		case reflect.Value:
			t = x.Type()
		default:
			t = reflect.TypeOf(v)
		}

		return bytesType == t || t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 && t.Elem().PkgPath() == ""
	}
}

// FullTypeName 返回包含包路径和指针层级的完整类型名。
func FullTypeName(rtype reflect.Type) string {
	buf := bytes.NewBuffer(nil)

	for rtype.Kind() == reflect.Pointer {
		buf.WriteByte('*')
		rtype = rtype.Elem()
	}

	if name := rtype.Name(); name != "" {
		if pkgPath := rtype.PkgPath(); pkgPath != "" {
			buf.WriteString(pkgPath)
			buf.WriteRune('.')
		}
		buf.WriteString(name)
		return buf.String()
	}

	buf.WriteString(rtype.String())
	return buf.String()
}

// Deref 递归解引用指针类型，返回最终元素类型。
func Deref(tpe reflect.Type) reflect.Type {
	if tpe.Kind() == reflect.Pointer {
		return Deref(tpe.Elem())
	}
	return tpe
}
