package anyjson

import "context"

func Transform(ctx context.Context, v Valuer, transform func(v Valuer, keyPath ...any) Valuer) Valuer {
	t := &transformer{
		transform: transform,
	}
	return t.Next(ctx, v, nil)
}

type transformer struct {
	transform func(v Valuer, keyPath ...any) Valuer
}

func (t *transformer) Next(ctx context.Context, v Valuer, keyPath []any) Valuer {
	switch x := v.(type) {
	case *Object:
		o := &Object{}

		for k, v := range x.KeyValues() {
			propValue := t.Next(ctx, v, append(keyPath, k))

			if propValue != nil {
				o.Set(k, propValue)
			}
		}

		return o
	case *Array:
		a := &Array{}

		for i, v := range x.IndexedValues() {
			if itemValue := t.Next(ctx, v, append(keyPath, i)); itemValue != nil {
				a.Append(itemValue)
			}
		}

		return a
	default:
		return t.transform(v, keyPath...)
	}
}
