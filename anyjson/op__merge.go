package anyjson

type MergeOption interface {
	ApplyToMerge(m *merger)
}

func Merge[T Valuer](base T, patch T, opts ...MergeOption) T {
	m := &merger{}
	for _, opt := range opts {
		opt.ApplyToMerge(m)
	}

	switch x := any(patch).(type) {
	case *Object:
		if b, ok := any(base).(*Object); ok {
			return any(m.mergeObject(b, x)).(T)
		} else {
			return patch
		}
	case *Array:
		if b, ok := any(base).(*Array); ok {
			return any(m.mergeArray(b, x)).(T)
		} else {
			return patch
		}
	default:
		return patch
	}
}

type merger struct {
	arrayMergeKey     string
	nullOp            NullOp
	emptyObjectAsNull bool
}

func (m *merger) maybeClone(valuer Valuer) Valuer {
	switch x := valuer.(type) {
	case *Object:
		return m.mergeObject(&Object{}, x)
	case *Array:
		return m.mergeArray(&Array{}, x)
	}
	return valuer
}

func (m *merger) mergeObject(left *Object, right *Object) Valuer {
	if right == nil || right.Len() == 0 {
		if m.emptyObjectAsNull && left.Len() == 0 {
			return &Null{}
		}
		return left
	}

	merged := &Object{}

	for key, valuer := range left.KeyValues() {
		rightValue, rightExists := right.Get(key)
		if rightExists {
			switch x := rightValue.(type) {
			case *Array:
				if leftValue, ok := valuer.(*Array); ok {
					valuer = m.mergeArray(leftValue, x)
				} else {
					valuer = m.maybeClone(x)
				}
			case *Object:
				if leftValue, ok := valuer.(*Object); ok {
					valuer = m.mergeObject(leftValue, x)
				} else {
					valuer = m.maybeClone(x)
				}
			case *Null:
				if m.nullOp == NullAsRemover {
					// don't merge null valuer
					valuer = &Null{}
					continue
				}
			default:
				valuer = x
			}
		}

		if _, ok := valuer.(*Null); !ok {
			if !rightExists {
				// when right value is not exists , and left prop is bool
				// drop left value
				if _, ok := valuer.(*Boolean); ok {
					merged.Delete(key)
					continue
				}
			}

			merged.Set(key, valuer)
		}

	}

	for key, value := range right.KeyValues() {
		value = m.maybeClone(value)

		if _, ok := value.(*Null); ok {
			if m.nullOp == NullIgnore {
				continue
			} else if m.nullOp == NullAsRemover {
				merged.Delete(key)
				continue
			}
		}

		if _, ok := left.Get(key); !ok {
			merged.Set(key, value)
		}
	}

	if m.emptyObjectAsNull && merged.Len() == 0 {
		return &Null{}
	}

	return merged
}

func (m *merger) mergeArray(left *Array, right *Array) *Array {
	if arrayMergeKey := m.arrayMergeKey; arrayMergeKey != "" {
		mergedArr := &Array{}
		processed := map[int]bool{}

		reset := func() {
			mergedArr = &Array{}
			processed = map[int]bool{}
		}

		findRightItemObjByValue := func(leftItemMergeKeyValue Valuer) (int, Valuer) {
			for i, item := range right.IndexedValues() {
				if itemObject, ok := item.(*Object); ok {
					if itemMergeKeyValue, ok := itemObject.Get(arrayMergeKey); ok {
						if Equal(itemMergeKeyValue, leftItemMergeKeyValue) {
							return i, item
						}
					}
				}
			}
			return 0, nil
		}

	MergeWithMergeKey:
		for leftItem := range left.Values() {
			if leftItemObj, ok := leftItem.(*Object); ok {
				if value, ok := leftItemObj.Get(arrayMergeKey); ok {
					if idx, found := findRightItemObjByValue(value); found != nil {
						obj := found.(*Object)
						processed[idx] = true

						if op, ok := IsPatchObject(obj); ok {
							switch op {
							case PatchOpReplace:
								mergedArr.Append(m.maybeClone(obj))
							case PatchOpDelete:
								continue
							}
						}

						mergedArr.Append(m.mergeObject(leftItemObj, obj))
						continue
					}
				} else {
					reset()
					break MergeWithMergeKey
				}
			} else {
				reset()
				break MergeWithMergeKey
			}

			mergedArr.Append(m.maybeClone(leftItem))
		}

		for idx, item := range right.IndexedValues() {
			if _, ok := processed[idx]; ok {
				continue
			}

			mergedArr.Append(m.maybeClone(item))
		}

		return mergedArr
	}

	mergedArr := &Array{}
	for item := range right.Values() {
		mergedArr.Append(m.maybeClone(item))
	}
	return mergedArr
}
