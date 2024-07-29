package anyjson

type NullOp int

const (
	NullIgnore NullOp = iota
	NullAsRemover
)

func WithNullOp(op NullOp) *NullOption {
	return &NullOption{nullOp: op}
}

type NullOption struct {
	nullOp NullOp
}

func (n *NullOption) ApplyToMerge(m *merger) {
	m.nullOp = n.nullOp
}

func WithEmptyObjectAsNull() *EmptyObjectAsNullOption {
	return &EmptyObjectAsNullOption{
		emptyObjectAsNull: true,
	}
}

type EmptyObjectAsNullOption struct {
	emptyObjectAsNull bool
}

func (o *EmptyObjectAsNullOption) ApplyToMerge(m *merger) {
	m.emptyObjectAsNull = o.emptyObjectAsNull
}

func (o *EmptyObjectAsNullOption) ApplyToDiff(d *differ) {
	d.emptyObjectAsNull = o.emptyObjectAsNull
}

func WithArrayMergeKey(key string) *ArrayMergeKeyOption {
	return &ArrayMergeKeyOption{
		arrayMergeKey: key,
	}
}

type ArrayMergeKeyOption struct {
	arrayMergeKey string
}

func (o *ArrayMergeKeyOption) ApplyToMerge(m *merger) {
	m.arrayMergeKey = o.arrayMergeKey
}

func (o *ArrayMergeKeyOption) ApplyToDiff(d *differ) {
	d.arrayMergeKey = o.arrayMergeKey
}
