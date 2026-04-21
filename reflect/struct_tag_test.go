package reflect_test

import (
	"testing"

	"github.com/octohelm/x/cmp"
	reflectx "github.com/octohelm/x/reflect"
	. "github.com/octohelm/x/testing/v2"
)

func TestParseStructTags(t *testing.T) {
	t.Run("解析多个标签并保留转义内容", func(t *testing.T) {
		tags := reflectx.ParseStructTags(`json:"name,omitempty" yaml:"display_name" note:"say \"hi\""`)

		Then(t, "应按标签键拆解原始值",
			Expect(tags["json"], Equal(reflectx.StructTag("name,omitempty"))),
			Expect(tags["yaml"], Equal(reflectx.StructTag("display_name"))),
			Expect(tags["note"], Equal(reflectx.StructTag(`say "hi"`))),
		)
	})

	t.Run("遇到非法片段时保留此前已解析结果", func(t *testing.T) {
		tags := reflectx.ParseStructTags(`json:"name,omitempty" invalid yaml:"ignored"`)

		Then(t, "应停止继续解析但保留先前结果",
			Expect(tags["json"], Equal(reflectx.StructTag("name,omitempty"))),
			Expect(tags["yaml"], Be(cmp.Zero[reflectx.StructTag]())),
		)
	})

	t.Run("空白标签文本返回空映射", func(t *testing.T) {
		Then(t, "应返回空映射",
			Expect(len(reflectx.ParseStructTags("   ")), Equal(0)),
		)
	})
}

func TestStructTagName(t *testing.T) {
	t.Run("带 flag 的标签返回主名称", func(t *testing.T) {
		Then(t, "应忽略逗号后的 flag",
			Expect(reflectx.StructTag("name,omitempty").Name(), Equal("name")),
		)
	})

	t.Run("只有 flag 时返回空名称", func(t *testing.T) {
		Then(t, "前缀为空时名称应为空字符串",
			Expect(reflectx.StructTag(",omitempty").Name(), Equal("")),
		)
	})

	t.Run("不带 flag 时返回原值", func(t *testing.T) {
		Then(t, "应直接返回整个标签值",
			Expect(reflectx.StructTag("name").Name(), Equal("name")),
		)
	})
}

func TestStructTagHasFlag(t *testing.T) {
	t.Run("命中已存在的 flag", func(t *testing.T) {
		Then(t, "存在对应 flag 时应返回 true",
			Expect(reflectx.StructTag("name,omitempty,string").HasFlag("omitempty"), Be(cmp.True())),
			Expect(reflectx.StructTag("name,omitempty,string").HasFlag("string"), Be(cmp.True())),
		)
	})

	t.Run("未命中时返回 false", func(t *testing.T) {
		Then(t, "不存在对应 flag 时应返回 false",
			Expect(reflectx.StructTag("name,omitempty,string").HasFlag("inline"), Be(cmp.False())),
		)
	})

	t.Run("仅匹配完整 flag", func(t *testing.T) {
		Then(t, "不应把子串误判为已命中",
			Expect(reflectx.StructTag("name,omitempty,string").HasFlag("empty"), Be(cmp.False())),
			Expect(reflectx.StructTag("name,omitempty,string").HasFlag("str"), Be(cmp.False())),
			Expect(reflectx.StructTag(",omitempty").HasFlag("omitempty"), Be(cmp.True())),
		)
	})
}
