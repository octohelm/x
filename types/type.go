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

	// Name 返回类型名；对匿名类型可能为空。
	Name() string
	// PkgPath 返回类型所在包路径；对内建类型或匿名类型可能为空。
	PkgPath() string
	// String 返回类型的字符串表示。
	String() string
	// Kind 返回与 reflect.Kind 对齐的种类信息。
	Kind() reflect.Kind
	// Implements 报告当前类型是否实现 u。
	Implements(u Type) bool
	// AssignableTo 报告当前类型是否可赋值给 u。
	AssignableTo(u Type) bool
	// ConvertibleTo 报告当前类型是否可转换为 u。
	ConvertibleTo(u Type) bool
	// Comparable 报告当前类型的值是否可比较。
	Comparable() bool

	// Key 返回 map 类型的键类型。
	Key() Type
	// Elem 返回指针、切片、数组、map 或 chan 的元素类型。
	Elem() Type
	// Len 返回数组长度。
	Len() int

	// NumField 返回结构体字段数。
	NumField() int
	// Field 返回指定索引的结构体字段。
	Field(i int) StructField
	// FieldByName 按字段名查找结构体字段。
	FieldByName(name string) (StructField, bool)
	// FieldByNameFunc 按匹配函数查找结构体字段。
	FieldByNameFunc(match func(string) bool) (StructField, bool)

	// NumMethod 返回方法数。
	NumMethod() int
	// Method 返回指定索引的方法。
	Method(i int) Method
	// MethodByName 按名称查找方法。
	MethodByName(name string) (Method, bool)

	// IsVariadic 报告函数类型是否为可变参数。
	IsVariadic() bool

	// NumIn 返回函数参数数。
	NumIn() int
	// In 返回指定索引的参数类型。
	In(i int) Type
	// NumOut 返回函数返回值数。
	NumOut() int
	// Out 返回指定索引的返回值类型。
	Out(i int) Type
}

// Method 表示统一抽象下的方法元数据。
type Method interface {
	// PkgPath 返回方法所属包路径。
	PkgPath() string
	// Name 返回方法名。
	Name() string
	// Type 返回方法签名类型。
	Type() Type
}

// StructField 表示统一抽象下的结构体字段元数据。
type StructField interface {
	// PkgPath 返回字段所属包路径；导出字段通常为空。
	PkgPath() string
	// Name 返回字段名。
	Name() string
	// Tag 返回字段结构体标签。
	Tag() reflect.StructTag
	// Type 返回字段类型。
	Type() Type
	// Anonymous 报告字段是否为匿名嵌入字段。
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
