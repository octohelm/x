package anyjson

// NullOp 表示 merge 遇到 null 时的处理语义。
type NullOp int

const (
	// NullIgnore 表示 merge 时忽略补丁中的 null。
	NullIgnore NullOp = iota
	// NullAsRemover 表示 merge 时将 null 视为删除操作。
	NullAsRemover
)

// WithNullOp 配置 merge 时对 null 的处理方式。
func WithNullOp(op NullOp) *NullOption {
	return &NullOption{nullOp: op}
}

// NullOption 表示 merge 时与 null 语义相关的选项。
type NullOption struct {
	nullOp NullOp
}

func (n *NullOption) ApplyToMerge(m *merger) {
	m.nullOp = n.nullOp
}

// WithEmptyObjectAsNull 配置将空对象视为 null。
func WithEmptyObjectAsNull() *EmptyObjectAsNullOption {
	return &EmptyObjectAsNullOption{
		emptyObjectAsNull: true,
	}
}

// EmptyObjectAsNullOption 表示是否将空对象视为 null 的选项。
type EmptyObjectAsNullOption struct {
	emptyObjectAsNull bool
}

func (o *EmptyObjectAsNullOption) ApplyToMerge(m *merger) {
	m.emptyObjectAsNull = o.emptyObjectAsNull
}

func (o *EmptyObjectAsNullOption) ApplyToDiff(d *differ) {
	d.emptyObjectAsNull = o.emptyObjectAsNull
}

// WithArrayMergeKey 配置数组元素按指定键做对象级 merge 或 diff。
func WithArrayMergeKey(key string) *ArrayMergeKeyOption {
	return &ArrayMergeKeyOption{
		arrayMergeKey: key,
	}
}

// ArrayMergeKeyOption 表示数组合并键的选项。
type ArrayMergeKeyOption struct {
	arrayMergeKey string
}

func (o *ArrayMergeKeyOption) ApplyToMerge(m *merger) {
	m.arrayMergeKey = o.arrayMergeKey
}

func (o *ArrayMergeKeyOption) ApplyToDiff(d *differ) {
	d.arrayMergeKey = o.arrayMergeKey
}
