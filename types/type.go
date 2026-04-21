package types

import (
	"encoding"
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	reflectx "github.com/octohelm/x/reflect"
)

// Type 定义对 reflect.Type 与 go/types.Type 的统一抽象。
type Type interface {
	// Unwrap 返回底层的 reflect.Type 或 types.Type。
	Unwrap() any

	Name() string
	PkgPath() string
	String() string
	Kind() reflect.Kind
	Implements(u Type) bool
	AssignableTo(u Type) bool
	ConvertibleTo(u Type) bool
	Comparable() bool

	Key() Type
	Elem() Type
	Len() int

	NumField() int
	Field(i int) StructField
	FieldByName(name string) (StructField, bool)
	FieldByNameFunc(match func(string) bool) (StructField, bool)

	NumMethod() int
	Method(i int) Method
	MethodByName(name string) (Method, bool)

	IsVariadic() bool

	NumIn() int
	In(i int) Type
	NumOut() int
	Out(i int) Type
}

// Method 表示统一抽象下的方法元数据。
type Method interface {
	PkgPath() string
	Name() string
	Type() Type
}

// StructField 表示统一抽象下的结构体字段元数据。
type StructField interface {
	PkgPath() string
	Name() string
	Tag() reflect.StructTag
	Type() Type
	Anonymous() bool
}

// TryNew 尝试为 Type 创建一个新的零值。
//
// 当前仅支持基于 reflect.Type 的 RType。
func TryNew(u Type) (reflect.Value, bool) {
	switch t := u.(type) {
	case *RType:
		return reflectx.New(t.Type), true
	}
	return reflect.Value{}, false
}

var rtypeEncodingTextMarshaler = FromRType(reflect.TypeFor[encoding.TextMarshaler]())

// EncodingTextMarshalerTypeReplacer 在类型实现 encoding.TextMarshaler 时返回 string 替代类型。
func EncodingTextMarshalerTypeReplacer(u Type) (Type, bool) {
	switch t := u.(type) {
	case *RType:
		return FromRType(reflect.TypeFor[string]()), t.Implements(rtypeEncodingTextMarshaler)
	case *TType:
		return FromTType(types.Typ[types.String]), t.Implements(rtypeEncodingTextMarshaler)
	}
	return nil, false
}

// EachField 遍历结构体字段，并按标签规则展开匿名嵌入字段。
func EachField(typ Type, tagForName string, each func(field StructField, fieldDisplayName string, omitempty bool) bool, tagsForKeepingNested ...string) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name()
		fieldTag := field.Tag()

		fieldDisplayName, omitempty, keepNested := FieldDisplayName(fieldTag, tagForName, fieldName)

		if !ast.IsExported(fieldName) || fieldDisplayName == "-" {
			continue
		}

		fieldType := Deref(field.Type())

		if field.Anonymous() {
			switch fieldType.Kind() {
			case reflect.Struct:
				if !keepNested {
					for _, tag := range tagsForKeepingNested {
						if _, ok := fieldTag.Lookup(tag); ok {
							keepNested = true
							break
						}
					}
				}

				if !keepNested {
					EachField(fieldType, tagForName, each)
					continue
				}
			case reflect.Interface:
				continue
			}
		}

		if !each(field, fieldDisplayName, omitempty) {
			break
		}
	}
}

// Deref 递归解引用指针类型。
func Deref(typ Type) Type {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

// FullTypeName 返回包含包路径和指针层级的完整类型名。
func FullTypeName(typ Type) string {
	if typ == nil {
		return "nil"
	}

	buf := &strings.Builder{}

	for typ.Kind() == reflect.Pointer {
		buf.WriteByte('*')
		typ = typ.Elem()
	}

	if name := typ.Name(); name != "" {
		if pkgPath := typ.PkgPath(); pkgPath != "" {
			buf.WriteString(pkgPath)
			buf.WriteRune('.')
		}
		buf.WriteString(name)
		return buf.String()
	}

	buf.WriteString(typ.String())
	return buf.String()
}

// FieldDisplayName 根据结构体标签计算字段显示名及相关标志。
func FieldDisplayName(structTag reflect.StructTag, namedTagKey string, defaultName string) (string, bool, bool) {
	jsonTag, exists := structTag.Lookup(namedTagKey)
	if !exists {
		return defaultName, false, exists
	}
	omitempty := strings.Index(jsonTag, "omitempty") > 0
	idxOfComma := strings.IndexRune(jsonTag, ',')
	if jsonTag == "" || idxOfComma == 0 {
		return defaultName, omitempty, true
	}
	if idxOfComma == -1 {
		return jsonTag, omitempty, true
	}
	return jsonTag[0:idxOfComma], omitempty, true
}
