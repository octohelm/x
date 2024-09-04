package anyjson

type DiffOption interface {
	ApplyToDiff(d *differ)
}

func Diff[T any](template *T, live *T, opts ...DiffOption) (Valuer, error) {
	m := &differ{
		arrayMergeKey:     "name",
		emptyObjectAsNull: true,
	}
	for _, opt := range opts {
		opt.ApplyToDiff(m)
	}

	base, err := FromValue(template)
	if err != nil {
		return nil, err
	}
	patch, err := FromValue(live)
	if err != nil {
		return nil, err
	}

	switch x := any(patch).(type) {
	case *Object:
		if b, ok := any(base).(*Object); ok {
			return m.diffObject(b, x), nil
		} else {
			return patch, nil
		}
	case *Array:
		if b, ok := any(base).(*Array); ok {
			return m.diffArray(b, x), nil
		} else {
			return patch, nil
		}
	default:
		return patch, nil
	}
}

type differ struct {
	arrayMergeKey     string
	emptyObjectAsNull bool
}

func (d *differ) maybeClone(valuer Valuer) Valuer {
	switch x := valuer.(type) {
	case *Object:
		return d.diffObject(&Object{}, x)
	case *Array:
		return d.diffArray(&Array{}, x)
	}
	return valuer
}

func (d *differ) diffObject(left *Object, right *Object) Valuer {
	if right == nil || right.Len() == 0 {
		if d.emptyObjectAsNull && left.Len() == 0 {
			return &Null{}
		}
		return left
	}

	merged := &Object{}

	for key, valuer := range left.KeyValues() {
		if rightValue, ok := right.Get(key); ok {
			switch x := rightValue.(type) {
			case *Array:
				if leftValue, ok := valuer.(*Array); ok {
					valuer = d.diffArray(leftValue, x)
				} else {
					valuer = d.maybeClone(x)
				}
			case *Object:
				if leftValue, ok := valuer.(*Object); ok {
					valuer = d.diffObject(leftValue, x)
				} else {
					valuer = d.maybeClone(x)
				}
			case *Null:
				valuer = &Null{}
			default:
				// when value equal, should remove
				if Equal(rightValue, valuer) {
					valuer = &Null{}
				} else {
					valuer = x
				}
			}
		}

		if _, ok := valuer.(*Null); !ok {
			merged.Set(key, valuer)
		}
	}

	for key, value := range right.KeyValues() {
		value = d.maybeClone(value)

		if _, ok := value.(*Null); ok {
			merged.Delete(key)
		}

		if _, ok := left.Get(key); !ok {
			merged.Set(key, value)
		}
	}

	if d.emptyObjectAsNull && merged.Len() == 0 {
		return &Null{}
	}

	return merged
}

func (d *differ) diffArray(left *Array, right *Array) Valuer {
	if arrayMergeKey := d.arrayMergeKey; arrayMergeKey != "" {
		mergedArr := &Array{}

		processed := map[int]bool{}

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

		for leftItem := range left.Values() {
			if leftItemObj, ok := leftItem.(*Object); ok {
				if mergeKeyValue, ok := leftItemObj.Get(arrayMergeKey); ok {
					if idx, found := findRightItemObjByValue(mergeKeyValue); found != nil {
						processed[idx] = true

						diff := d.diffObject(leftItemObj, found.(*Object))

						switch x := diff.(type) {
						case *Object:
							x.Set(arrayMergeKey, mergeKeyValue)
							mergedArr.Append(x)
						case *Null:
							o := &Object{}
							o.Set(arrayMergeKey, mergeKeyValue)
							o.Set(PatchKey, StringOf(PatchOpDelete))
							mergedArr.Append(o)
						}

						continue
					}

					o := &Object{}
					o.Set(arrayMergeKey, mergeKeyValue)
					o.Set(PatchKey, StringOf(string(PatchOpDelete)))

					mergedArr.Append(o)

					continue
				} else {
					return d.diffFullArray(left, right)
				}
			} else {
				return d.diffFullArray(left, right)
			}
		}

		for idx, item := range right.IndexedValues() {
			if _, ok := processed[idx]; ok {
				continue
			}

			mergedArr.Append(d.maybeClone(item))
		}

		return mergedArr
	}

	mergedArr := &Array{}
	for item := range right.Values() {
		mergedArr.Append(d.maybeClone(item))
	}
	return mergedArr
}

func (d *differ) diffFullArray(left *Array, right *Array) Valuer {
	if left.Len() != right.Len() {
		return d.maybeClone(right)
	}

	allEqual := true

	for i, leftItem := range left.IndexedValues() {
		rightItem, _ := right.Index(i)

		if !Equal(leftItem, rightItem) {
			allEqual = false
			break
		}
	}

	if allEqual {
		return &Null{}
	}
	return d.maybeClone(right)
}
