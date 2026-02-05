package v2_test

import (
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/octohelm/x/cmp"
	"github.com/octohelm/x/ptr"
	. "github.com/octohelm/x/testing/v2"
)

// 自定义错误类型用于测试
type ErrTest struct {
	Message string
}

func (e *ErrTest) Error() string {
	return "err test: " + e.Message
}

type ErrAnother struct {
	Code int
}

func (e *ErrAnother) Error() string {
	return "another error"
}

type User struct {
	Name  string
	Age   int
	Email string
}

func TestAPIs(t *testing.T) {
	// 基础断言示例
	t.Run("基础断言", func(t *testing.T) {
		value := 1

		t.Run("相等断言", func(t *testing.T) {
			Then(t, "值应该等于2",
				Expect(value+1, Equal(2)),
			)
		})

		t.Run("不相等断言", func(t *testing.T) {
			Then(t, "值不应该等于 3",
				Expect(value,
					NotEqual(3),
				),
			)
		})

		t.Run("使用 cmp 包的比较函数", func(t *testing.T) {
			Then(t, "值应该大于0",
				Expect(value,
					Be(cmp.Gt(0)),
				),
			)

			Then(t, "值应该小于等于1",
				Expect(value,
					Be(cmp.Lte(1)),
				),
			)
		})
	})

	// 错误处理示例
	t.Run("错误断言", func(t *testing.T) {
		t.Run("期望没有错误", func(t *testing.T) {
			Then(t, "函数应该成功执行",
				ExpectMust(func() error {
					// 模拟成功操作
					return nil
				}),
			)
		})

		t.Run("期望特定错误", func(t *testing.T) {
			testErr := &ErrTest{Message: "something wrong"}

			Then(t, "应该返回ErrTest类型错误",
				ExpectDo(
					func() error {
						return testErr
					},
					ErrorAs(ptr.Ptr(&ErrTest{})),
				),
			)
		})

		t.Run("错误链的Is断言", func(t *testing.T) {
			Then(t, "错误链中包含 os.ErrNotExist",
				ExpectDo(
					func() error {
						return errors.Join(os.ErrNotExist, &ErrTest{Message: "wrapped"})
					},
					ErrorIs(os.ErrNotExist),
				),
			)
		})

		t.Run("错误排除断言", func(t *testing.T) {
			Then(t, "错误不是特定的类型",
				ExpectDo(
					func() error {
						return &ErrAnother{Code: 500}
					},
					ErrorNotAs(ptr.Ptr(&ErrTest{})),
					ErrorNotIs(os.ErrNotExist),
				),
			)
		})

		t.Run("期望错误消息匹配", func(t *testing.T) {
			t.Run("错误消息应该包含特定字符串", func(t *testing.T) {
				Then(t, "错误消息应该包含'something'",
					ExpectDo(
						func() error {
							return &ErrTest{Message: "something went wrong"}
						},
						ErrorMatch(regexp.MustCompile("something")),
					),
				)
			})

			t.Run("错误消息应该匹配正则表达式", func(t *testing.T) {
				Then(t, "错误消息应该匹配正则模式",
					ExpectDo(
						func() error {
							return &ErrTest{Message: "error code: 404 not found"}
						},
						ErrorMatch(regexp.MustCompile(`code: \d+`)),
					),
				)
			})

			t.Run("排除匹配的消息", func(t *testing.T) {
				Then(t, "错误消息不应该包含'ignore'",
					ExpectDo(
						func() error {
							return &ErrTest{Message: "something important"}
						},
						ErrorNotMatch(regexp.MustCompile("ignore")),
					),
				)
			})

			t.Run("组合使用错误匹配", func(t *testing.T) {
				Then(t, "错误应该是特定类型且消息匹配",
					ExpectDo(
						func() error {
							return &ErrTest{Message: "validation failed: email invalid"}
						},
						ErrorAs(ptr.Ptr(&ErrTest{})),
						ErrorMatch(regexp.MustCompile("validation failed")),
						ErrorMatch(regexp.MustCompile("email invalid")),
						ErrorNotMatch(regexp.MustCompile("password")),
					),
				)
			})
		})
	})

	// 复杂数据类型示例
	t.Run("复杂数据类型", func(t *testing.T) {
		t.Run("结构体比较", func(t *testing.T) {
			user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}

			Then(t, "用户信息应该匹配",
				Expect(user,
					Equal(User{Name: "Alice", Age: 30, Email: "alice@example.com"}),
				),
			)
		})

		t.Run("切片比较", func(t *testing.T) {
			numbers := []int{1, 2, 3, 4, 5}

			Then(t, "切片应该相等",
				Expect(numbers,
					Equal([]int{1, 2, 3, 4, 5}),
				),
			)
		})

		t.Run("Map 比较", func(t *testing.T) {
			scores := map[string]int{"Alice": 95, "Bob": 87}

			Then(t, "Map应该相等",
				Expect(scores,
					Equal(map[string]int{"Alice": 95, "Bob": 87}),
				),
			)
		})
	})

	// 条件组合断言示例
	t.Run("组合断言", func(t *testing.T) {
		value := 42

		Then(t, "多个条件同时满足",
			Expect(value,
				Equal(42),
				Be(cmp.Gt(40)),
				Be(cmp.Lt(50)),
				NotEqual(0),
			),
		)
	})

	// Must函数使用示例
	t.Run("Must函数使用", func(t *testing.T) {
		t.Run("Must处理错误", func(t *testing.T) {
			// 如果这里发生错误，测试会立即失败
			result := MustValue(t, func() (string, error) {
				return "success", nil
			})

			Then(t, "结果应该是success",
				Expect(result, Equal("success")),
			)
		})

		t.Run("MustValues多个返回值", func(t *testing.T) {
			a, b := MustValues(t, func() (int, string, error) {
				return 100, "ok", nil
			})

			Then(t, "多个返回值都正确",
				Expect(a, Equal(100)),
				Expect(b, Equal("ok")),
			)
		})
	})

	// 快照测试示例
	t.Run("快照测试", func(t *testing.T) {
		t.Run("JSON 数据快照", func(t *testing.T) {
			data := map[string]any{
				"id":   1231,
				"name": "test",
				"tags": []string{"go", "testing"},
			}

			raw := MustValue(t, func() ([]byte, error) {
				return json.Marshal(data, json.Deterministic(true))
			})

			Then(t, "JSON 快照匹配",
				ExpectMustValue(
					func() (Snapshot, error) {
						return SnapshotOf(
							SnapshotFileFromRaw("test.json", raw),
						), nil
					},
					MatchSnapshot("json_test"),
				),
			)
		})

		t.Run("多文件快照", func(t *testing.T) {
			Then(t, "多文件快照匹配",
				ExpectMustValue(
					func() (Snapshot, error) {
						// 创建多个快照文件
						file1 := SnapshotFileFromRaw("config.yaml", []byte("env: production\nversion: 1.0.0"))
						file2 := SnapshotFileFromRaw("data.json", []byte(`{"count": 100}`))
						file3 := SnapshotFileFromRaw("readme.md", []byte("# Test Document"))

						return SnapshotOf(file1, file2, file3), nil
					},
					MatchSnapshot("multi_file_test"),
				),
			)
		})
	})

	// 嵌套测试示例
	t.Run("嵌套测试场景", func(t *testing.T) {
		// 模拟复杂测试场景
		t.Run("用户注册流程", func(t *testing.T) {
			var user User

			t.Run("GIVEN 新用户信息", func(t *testing.T) {
				user = User{Name: "Bob", Age: 25}

				Then(t, "用户信息初始化正确",
					Expect(user.Name, Equal("Bob")),
					Expect(user.Age, Be(cmp.Gte(18))),
				)
			})

			t.Run("WHEN 设置邮箱", func(t *testing.T) {
				user.Email = "bob@example.com"

				Then(t, "邮箱设置成功",
					Expect(user.Email, Equal("bob@example.com")),
				)
			})

			t.Run("THEN 完整用户信息", func(t *testing.T) {
				Then(t, "用户信息完整",
					Expect(user,
						Equal(User{Name: "Bob", Age: 25, Email: "bob@example.com"}),
					),
				)
			})
		})
	})

	// 边界条件测试
	t.Run("边界条件", func(t *testing.T) {
		t.Run("nil值处理", func(t *testing.T) {
			var ptr *int

			Then(t, "指针应该为nil",
				Expect(ptr,
					Be(cmp.Nil[*int]()),
				),
			)
		})

		t.Run("零值处理", func(t *testing.T) {
			empty := ""

			Then(t, "字符串应该为零值",
				Expect(empty,
					Be(cmp.Zero[string]()),
				),
			)
		})

		t.Run("非零值断言", func(t *testing.T) {
			value := "hello"

			Then(t, "字符串应该非零",
				Expect(value,
					Be(cmp.NotZero[string]()),
				),
			)
		})
	})
}

// 测试 Must 在失败时的行为
func TestMustFailure(t *testing.T) {
	// 这个测试展示了当 Must 遇到错误时的行为
	t.Run("Must 失败示例", func(t *testing.T) {
		t.Skip()

		// 注意：这会导致测试失败，演示了Must的错误处理
		// 实际使用中应该捕获预期的错误
		Must(t, func() error {
			return errors.New("expected error")
		})
	})
}
