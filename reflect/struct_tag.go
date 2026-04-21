package reflect

import (
	"strconv"
	"strings"
)

// ParseStructTags 将结构体标签文本拆解为键到 StructTag 的映射。
func ParseStructTags(tag string) map[string]StructTag {
	tagFlags := map[string]StructTag{}

	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := tag[:i+1]
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		tagFlags[name] = StructTag(value)
	}

	return tagFlags
}

// StructTag 表示单个标签键对应的原始值。
type StructTag string

// Name 返回标签中的主名称部分。
func (t StructTag) Name() string {
	s := string(t)

	if i := strings.Index(s, ","); i >= 0 {
		if i == 0 {
			return ""
		}
		return s[0:i]
	}

	return s
}

// HasFlag 判断标签中是否包含指定 flag。
func (t StructTag) HasFlag(flag string) bool {
	s := string(t)

	if i := strings.Index(s, ","); i >= 0 {
		for _, part := range strings.Split(s[i+1:], ",") {
			if part == flag {
				return true
			}
		}
	}

	return false
}
