package anyjson

var PatchKey = "$patch"

type PatchOp string

const (
	PatchOpReplace = "replace"
	PatchOpDelete  = "delete"
)

func IsPatchObject(obj *Object) (PatchOp, bool) {
	if obj == nil {
		return "", false
	}

	if str, ok := obj.Get(PatchKey); ok {
		if s, ok := str.Value().(string); ok {
			switch x := PatchOp(s); x {
			case PatchOpReplace, PatchOpDelete:
				return x, true
			}
		}
	}

	return "", false
}
