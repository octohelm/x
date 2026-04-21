package anyjson

// PatchKey 是补丁对象中声明补丁操作类型的保留键。
var PatchKey = "$patch"

// PatchOp 表示补丁对象支持的操作类型。
type PatchOp string

const (
	// PatchOpReplace 表示用当前对象替换目标对象。
	PatchOpReplace = "replace"
	// PatchOpDelete 表示删除目标对象。
	PatchOpDelete = "delete"
)

// IsPatchObject 判断对象是否为携带补丁操作的补丁对象。
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
