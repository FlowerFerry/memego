package memego

import "C"
import (
	"errors"
	"unsafe"
)

/*

#include <stdlib.h>
#include "meme/string.h"

MemeInteger_t mmgoinvoke_match_cond_byte_fn(MemeByte_t byte, void* userdata);
*/
import "C"

type String struct {
	data C.mms_stack_t
}

type MatchCondByteFn func(b byte, userdata unsafe.Pointer) int

type InvokeMatchParameter struct {
	condFn   MatchCondByteFn
	userdata unsafe.Pointer
}

// 接管 C.mms_stack_t 类型
func CreateStringWithTakeOver(data C.mms_stack_t) *String {
	return &String{data: data}
}

func CreateString() *String {
	s := &String{}
	C.MemeStringStack_init(&s.data, C.MMS__OBJECT_SIZE)
	return s
}

func CreateStringByGoStr(str string) (*String, int) {
	b := []byte(str)
	s := &String{}
	result := C.MemeStringStack_initByU8bytes(
		&s.data, C.MMS__OBJECT_SIZE,
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)))
	if result != 0 {
		return nil, int(result)
	}
	return s, 0
}

//! \brief 创建一个字符串对象
//!
//! 该函数不安全的原因之一是当在子动态库中使用后，分配的对象传递到父进程里并生命周期比该动态库的生命周期长的情况下，会引发问题。
//! \return 返回字符串对象指针
//func CreateStringUnsafeByGoStr(str string) *String {
//	cstr := C.CString(str)
//
//	s := &String{}
//	userdata := &StringUserData{ptr: cstr, len: C.size_t(len(str))}
//	C.MemeStringStack_initTakeOverUserObject(
//		&s.data, C.MMS__OBJECT_SIZE,
//		unsafe.Pointer(userdata), unsafe.Pointer(C.MmgoUserData_destroy),
//		unsafe.Pointer(C.MmgoUserData_cstr), unsafe.Pointer(C.MmgoUserData_size))
//	return s
//}

func CreateStringByBytes(bytes []byte) (*String, int) {
	//cstr := C.CBytes(bytes)
	//defer C.free(unsafe.Pointer(cstr))
	s := &String{}
	result := C.MemeStringStack_initByU8bytes(
		&s.data, C.MMS__OBJECT_SIZE,
		(*C.uchar)(unsafe.Pointer(&bytes[0])), C.MemeInteger_t(len(bytes)))
	if result != 0 {
		return nil, int(result)
	}
	return s, 0
}

func CreateStringByOther(str *String) (*String, int) {
	s := CreateString()
	result := s.AssignByOther(str)
	if result != 0 {
		return nil, result
	}
	return s, 0
}

func CreateStringByCStr(cstr *C.char, len C.size_t) (*String, int) {
	s := &String{}
	result := C.MemeStringStack_initByU8bytes(
		&s.data, C.MMS__OBJECT_SIZE, (*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len))
	if result != 0 {
		return nil, int(result)
	}
	return s, 0
}

func (s *String) Destroy() {
	C.MemeStringStack_unInit(&s.data, C.MMS__OBJECT_SIZE)
}

func (s *String) AssignByGoStr(str string) int {
	other, rc := CreateStringByGoStr(str)
	if rc != 0 {
		return rc
	}
	defer other.Destroy()
	return s.AssignByOther(other)

	//cstr := C.CString(str)
	//defer C.free(unsafe.Pointer(cstr))
	//result := C.MemeStringStack_assignByU8bytes(
	//	&s.data, C.MMS__OBJECT_SIZE,
	//	(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len(str)))
	//return int(result)
}

func (s *String) AssignByBytes(bytes []byte) int {
	other, rc := CreateStringByBytes(bytes)
	if rc != 0 {
		return rc
	}
	defer other.Destroy()
	return s.AssignByOther(other)

	//result := C.MemeStringStack_assignByU8bytes(
	//	&s.data, C.MMS__OBJECT_SIZE,
	//	(*C.uchar)(unsafe.Pointer(&bytes[0])), C.MemeInteger_t(len(bytes)))
	//return int(result)
}

func (s *String) AssignByOther(str *String) int {
	result := C.MemeStringStack_assign(
		&s.data, C.MMS__OBJECT_SIZE, str.ToMms())
	return int(result)
}

func (s *String) AssignByCStr(cstr *C.char, len C.size_t) int {
	result := C.MemeStringStack_assignByU8bytes(
		&s.data, C.MMS__OBJECT_SIZE, (*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len))
	return int(result)
}

// clear
func (s *String) Clear() {
	C.MemeStringStack_reset(&s.data, C.MMS__OBJECT_SIZE)
}

func (s *String) IsEmpty() bool {
	return C.MemeString_isEmpty(s.ToMms()) != 0
}

// get btye size
func (s *String) Size() int {
	return int(C.MemeString_byteSize(s.ToMms()))
}

func (s *String) SizeToInteger() C.MemeInteger_t {
	return C.MemeString_byteSize(s.ToMms())
}

func (s *String) CStr() *C.char {
	return C.MemeString_cStr(s.ToMms())
}

func (s *String) IsSharedStorage() bool {
	return C.MemeString_isSharedStorageTypes(s.ToMms()) != 0
}

func (s *String) Equal(str *String) bool {
	result := C.int(0)
	C.MemeString_isEqualWithOther(
		s.ToMms(), str.ToMms(), &result)
	return result != 0
}

func (s *String) NotEqual(str *String) bool {
	return !s.Equal(str)
}

func (s *String) EqualGoStr(str string) bool {
	b := []byte(str)
	result := C.int(0)
	C.MemeString_isEqual(
		s.ToMms(), (*C.char)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)), &result)
	return result != 0
}

func (s *String) NotEqualGoStr(str string) bool {
	return !s.EqualGoStr(str)
}

func (s *String) EqualBytes(bytes []byte) bool {
	result := C.int(0)
	C.MemeString_isEqual(
		s.ToMms(), (*C.char)(unsafe.Pointer(&bytes[0])), C.MemeInteger_t(len(bytes)), &result)
	return result != 0
}

func (s *String) NotEqualBytes(bytes []byte) bool {
	return !s.EqualBytes(bytes)
}

func (s *String) EqualCStr(cstr *C.char, len C.size_t) bool {
	result := C.int(0)
	C.MemeString_isEqual(
		s.ToMms(), (*C.char)(unsafe.Pointer(cstr)), C.MemeInteger_t(len), &result)
	return result != 0
}

func (s *String) NotEqualCStr(cstr *C.char, len C.size_t) bool {
	return !s.EqualCStr(cstr, len)
}

func (s *String) ToMms() C.mms_t {
	return C.mms_t(unsafe.Pointer(&s.data))
}

func (s *String) ToString() string {
	return C.GoStringN(s.CStr(), C.int(s.SizeToInteger()))
}

func (s *String) ToBytes() []byte {
	return C.GoBytes(unsafe.Pointer(s.CStr()), C.int(s.SizeToInteger()))
}

// to en unpper
func (s *String) ToEnUpper() *String {
	stack := C.MemeStringStack_toEnUpper(&s.data, C.MMS__OBJECT_SIZE)
	return &String{data: stack}
}

// to en lower
func (s *String) ToEnLower() *String {
	stack := C.MemeStringStack_toEnLower(&s.data, C.MMS__OBJECT_SIZE)
	return &String{stack}
}

// trim space
func (s *String) TrimSpace() *String {
	stack := C.MemeStringStack_trimSpace(&s.data, C.MMS__OBJECT_SIZE)
	return &String{stack}
}

// trim left space
func (s *String) TrimLeftSpace() *String {
	stack := C.MemeStringStack_trimLeftSpace(&s.data, C.MMS__OBJECT_SIZE)
	return &String{stack}
}

// trim right space
func (s *String) TrimRightSpace() *String {
	stack := C.MemeStringStack_trimRightSpace(&s.data, C.MMS__OBJECT_SIZE)
	return &String{stack}
}

func (s *String) TrimByCondByteFn(fn MatchCondByteFn, user unsafe.Pointer) *String {
	param := &InvokeMatchParameter{
		condFn:   fn,
		userdata: user,
	}
	stack := C.MemeStringStack_trimByCondByteFunc(
		&s.data, C.MMS__OBJECT_SIZE,
		(*C.MemeString_MatchCondByteFunc_t)(unsafe.Pointer(C.mmgoinvoke_match_cond_byte_fn)),
		unsafe.Pointer(param))
	return &String{stack}
}

// mid
func (s *String) Mid(start int, length int) *String {
	stack := C.MemeStringStack_mid(
		&s.data, C.MMS__OBJECT_SIZE, C.MemeInteger_t(start), C.MemeInteger_t(length))
	return &String{stack}
}

func (s *String) At(index int) (byte, bool) {
	p := C.MemeString_at(s.ToMms(), C.MemeInteger_t(index))
	if p == nil {
		return 0, false
	}
	return byte(*p), true
}

// storage type
func (s *String) StorageType() Storage_t {
	return Storage_t(C.MemeString_storageType(s.ToMms()))
}

func (s *String) IndexOf(str *String) int {
	return int(C.MemeString_indexOfWithOther(
		s.ToMms(),
		C.MemeInteger_t(0),
		str.ToMms(),
		C.MemeFlag_CaseSensitive))
}

func (s *String) IndexOfWithStartIndex(str *String, startIndex int) int {
	return int(C.MemeString_indexOfWithOther(
		s.ToMms(),
		C.MemeInteger_t(startIndex),
		str.ToMms(),
		C.MemeFlag_CaseSensitive))
}

func (s *String) IndexOfGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_indexOfWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) IndexOfWithStartIndexGoStr(str string, startIndex int) int {
	b := []byte(str)
	return int(C.MemeString_indexOfWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) IndexOfCStr(cstr *C.char, len C.size_t) int {
	return int(C.MemeString_indexOfWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

func (s *String) IndexOfWithStartIndexCStr(cstr *C.char, len C.size_t, startIndex int) int {
	return int(C.MemeString_indexOfWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

func (s *String) LastIndexOfGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_lastIndexOfWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) LastIndexOfWithStartIndexGoStr(str string, startIndex int) int {
	b := []byte(str)
	return int(C.MemeString_lastIndexOfWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) LastIndexOfCStr(cstr *C.char, len C.size_t) int {
	return int(C.MemeString_lastIndexOfWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

func (s *String) LastIndexOfWithStartIndexCStr(cstr *C.char, len C.size_t, startIndex int) int {
	return int(C.MemeString_lastIndexOfWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

// match count
func (s *String) MatchCountGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_matchCountWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) MatchCountWithStartIndexGoStr(str string, startIndex int) int {
	b := []byte(str)
	return int(C.MemeString_matchCountWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) MatchCountCStr(cstr *C.char, len C.size_t) int {
	return int(C.MemeString_matchCountWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

func (s *String) MatchCountWithStartIndexCStr(cstr *C.char, len C.size_t, startIndex int) int {
	return int(C.MemeString_matchCountWithUtf8bytes(
		s.ToMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

// start match
func (s *String) StartMatchGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_startsMatchWithUtf8bytes(
		s.ToMms(),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

// end match
func (s *String) EndMatchGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_endsMatchWithUtf8bytes(
		s.ToMms(),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

// swap
func (s *String) Swap(other *String) {
	C.MemeString_swap(
		s.ToMms(),
		other.ToMms())
}

func (s *String) SplitGoStr(str string) ([]*String, error) {
	stacks := make([]C.mms_stack_t, 4)
	stacksCount := C.MemeInteger_t(0)
	b := []byte(str)
	var list []*String
	for index := C.MemeInteger_t(0); index != C.MemeInteger_t(-1); {
		stacksCount = C.MemeInteger_t(len(stacks))
		result := C.MemeString_split(
			s.ToMms(),
			(*C.char)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
			C.MemeFlag_KeepEmptyParts,
			C.MemeFlag_AllSensitive,
			(*C.mms_stack_t)(unsafe.Pointer(&stacks[0])), &stacksCount, &index)
		if result != 0 {
			return nil, errors.New("split failed")
		}
		for i := C.MemeInteger_t(0); i < stacksCount; i++ {
			list = append(list, &String{stacks[index]})
		}
	}
	return list, nil
}

func (s *String) Split(str *String) ([]*String, error) {
	stacks := make([]C.mms_stack_t, 4)
	stacksCount := C.MemeInteger_t(0)
	var list []*String
	for index := C.MemeInteger_t(0); index != C.MemeInteger_t(-1); {
		stacksCount = C.MemeInteger_t(len(stacks))
		result := C.MemeString_split(
			s.ToMms(),
			str.CStr(), str.SizeToInteger(),
			C.MemeFlag_KeepEmptyParts,
			C.MemeFlag_AllSensitive,
			(*C.mms_stack_t)(unsafe.Pointer(&stacks[0])), &stacksCount, &index)
		if result != 0 {
			return nil, errors.New("split failed")
		}
		for i := C.MemeInteger_t(0); i < stacksCount; i++ {
			list = append(list, &String{stacks[index]})
		}
	}
	return list, nil
}

//export mmgoinvoke_match_cond_byte_fn
func mmgoinvoke_match_cond_byte_fn(b C.MemeByte_t, userdata unsafe.Pointer) C.MemeInteger_t {
	param := (*InvokeMatchParameter)(userdata)
	if param.condFn != nil {
		return C.MemeInteger_t(param.condFn(byte(b), param.userdata))
	}
	return C.MemeInteger_t(0)
}
