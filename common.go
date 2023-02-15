package memego

// #include <stdlib.h>
// #include "meme/string.h"
import "C"

type Storage_t int

const (
	StorageType_none       = Storage_t(C.MemeString_StorageType_none)
	StorageType_small      = Storage_t(C.MemeString_StorageType_small)
	StorageType_medium     = Storage_t(C.MemeString_StorageType_medium)
	StorageType_large      = Storage_t(C.MemeString_StorageType_large)
	StorageType_user       = Storage_t(C.MemeString_StorageType_user)
	StorageType_viewUnsafe = Storage_t(C.MemeString_UnsafeStorageType_view)
)

type SplitBehavior_t C.MemeFlag_SplitBehavior_t

const (
	KeepEmptyParts SplitBehavior_t = C.MemeFlag_KeepEmptyParts
	SkipEmptyParts SplitBehavior_t = C.MemeFlag_SkipEmptyParts
)
