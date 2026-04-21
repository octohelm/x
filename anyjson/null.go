package anyjson

// Null 表示 JSON null。
type Null struct{}

// Value 返回 nil。
func (*Null) Value() any {
	return nil
}

// String 返回 null 的 JSON 文本表示。
func (v *Null) String() string {
	return ToString(v)
}

func (v *Null) UnmarshalJSON(bytes []byte) error {
	return nil
}

func (v *Null) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
