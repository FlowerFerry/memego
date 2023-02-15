package memego

import "C"
import (
	"errors"
	"unsafe"
)

/*

#include <stdlib.h>
#include "meme/string.h"

MemeInteger_t mmsinvoke_match_cond_byte_fn(MemeByte_t byte, void* userdata);
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
	result := s.assignByOther(str)
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

func (s *String) destroy() {
	C.MemeStringStack_unInit(&s.data, C.MMS__OBJECT_SIZE)
}

func (s *String) assignByGoStr(str string) int {
	other, rc := CreateStringByGoStr(str)
	if rc != 0 {
		return rc
	}
	defer other.destroy()
	return s.assignByOther(other)

	//cstr := C.CString(str)
	//defer C.free(unsafe.Pointer(cstr))
	//result := C.MemeStringStack_assignByU8bytes(
	//	&s.data, C.MMS__OBJECT_SIZE,
	//	(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len(str)))
	//return int(result)
}

func (s *String) assignByBytes(bytes []byte) int {
	other, rc := CreateStringByBytes(bytes)
	if rc != 0 {
		return rc
	}
	defer other.destroy()
	return s.assignByOther(other)

	//result := C.MemeStringStack_assignByU8bytes(
	//	&s.data, C.MMS__OBJECT_SIZE,
	//	(*C.uchar)(unsafe.Pointer(&bytes[0])), C.MemeInteger_t(len(bytes)))
	//return int(result)
}

func (s *String) assignByOther(str *String) int {
	result := C.MemeStringStack_assign(
		&s.data, C.MMS__OBJECT_SIZE, str.toMms())
	return int(result)
}

func (s *String) assignByCStr(cstr *C.char, len C.size_t) int {
	result := C.MemeStringStack_assignByU8bytes(
		&s.data, C.MMS__OBJECT_SIZE, (*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len))
	return int(result)
}

// clear
func (s *String) clear() {
	C.MemeStringStack_reset(&s.data, C.MMS__OBJECT_SIZE)
}

func (s *String) isEmpty() bool {
	return C.MemeString_isEmpty(s.toMms()) != 0
}

// get btye size
func (s *String) size() int {
	return int(C.MemeString_byteSize(s.toMms()))
}

func (s *String) sizeToInteger() C.MemeInteger_t {
	return C.MemeString_byteSize(s.toMms())
}

func (s *String) cstr() *C.char {
	return C.MemeString_cStr(s.toMms())
}

func (s *String) isSharedStorage() bool {
	return C.MemeString_isSharedStorageTypes(s.toMms()) != 0
}

func (s *String) equal(str *String) bool {
	result := C.int(0)
	C.MemeString_isEqualWithOther(
		s.toMms(), str.toMms(), &result)
	return result != 0
}

func (s *String) notEqual(str *String) bool {
	return !s.equal(str)
}

func (s *String) equalGoStr(str string) bool {
	b := []byte(str)
	result := C.int(0)
	C.MemeString_isEqual(
		s.toMms(), (*C.char)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)), &result)
	return result != 0
}

func (s *String) notEqualGoStr(str string) bool {
	return !s.equalGoStr(str)
}

func (s *String) equalBytes(bytes []byte) bool {
	result := C.int(0)
	C.MemeString_isEqual(
		s.toMms(), (*C.char)(unsafe.Pointer(&bytes[0])), C.MemeInteger_t(len(bytes)), &result)
	return result != 0
}

func (s *String) notEqualBytes(bytes []byte) bool {
	return !s.equalBytes(bytes)
}

func (s *String) equalCStr(cstr *C.char, len C.size_t) bool {
	result := C.int(0)
	C.MemeString_isEqual(
		s.toMms(), (*C.char)(unsafe.Pointer(cstr)), C.MemeInteger_t(len), &result)
	return result != 0
}

func (s *String) notEqualCStr(cstr *C.char, len C.size_t) bool {
	return !s.equalCStr(cstr, len)
}

func (s *String) toMms() C.mms_t {
	return C.mms_t(unsafe.Pointer(&s.data))
}

func (s *String) toString() string {
	return C.GoStringN(s.cstr(), C.int(s.sizeToInteger()))
}

func (s *String) toBytes() []byte {
	return C.GoBytes(unsafe.Pointer(s.cstr()), C.int(s.sizeToInteger()))
}

// to en unpper
func (s *String) toEnUpper() *String {
	stack := C.MemeStringStack_toEnUpper(&s.data, C.MMS__OBJECT_SIZE)
	return &String{data: stack}
}

// to en lower
func (s *String) toEnLower() *String {
	stack := C.MemeStringStack_toEnLower(&s.data, C.MMS__OBJECT_SIZE)
	return &String{stack}
}

// trim space
func (s *String) trimSpace() *String {
	stack := C.MemeStringStack_trimSpace(&s.data, C.MMS__OBJECT_SIZE)
	return &String{stack}
}

// trim left space
func (s *String) trimLeftSpace() *String {
	stack := C.MemeStringStack_trimLeftSpace(&s.data, C.MMS__OBJECT_SIZE)
	return &String{stack}
}

// trim right space
func (s *String) trimRightSpace() *String {
	stack := C.MemeStringStack_trimRightSpace(&s.data, C.MMS__OBJECT_SIZE)
	return &String{stack}
}

func (s *String) trimByCondByteFn(fn MatchCondByteFn, user unsafe.Pointer) *String {
	param := &InvokeMatchParameter{
		condFn:   fn,
		userdata: user,
	}
	stack := C.MemeStringStack_trimByCondByteFunc(
		&s.data, C.MMS__OBJECT_SIZE,
		(*C.MemeString_MatchCondByteFunc_t)(unsafe.Pointer(C.mmsinvoke_match_cond_byte_fn)),
		unsafe.Pointer(param))
	return &String{stack}
}

// mid
func (s *String) mid(start int, length int) *String {
	stack := C.MemeStringStack_mid(
		&s.data, C.MMS__OBJECT_SIZE, C.MemeInteger_t(start), C.MemeInteger_t(length))
	return &String{stack}
}

func (s *String) at(index int) (byte, bool) {
	p := C.MemeString_at(s.toMms(), C.MemeInteger_t(index))
	if p == nil {
		return 0, false
	}
	return byte(*p), true
}

// storage type
func (s *String) storageType() Storage_t {
	return Storage_t(C.MemeString_storageType(s.toMms()))
}

func (s *String) indexOf(str *String) int {
	return int(C.MemeString_indexOfWithOther(
		s.toMms(),
		C.MemeInteger_t(0),
		str.toMms(),
		C.MemeFlag_CaseSensitive))
}

func (s *String) indexOfWithStartIndex(str *String, startIndex int) int {
	return int(C.MemeString_indexOfWithOther(
		s.toMms(),
		C.MemeInteger_t(startIndex),
		str.toMms(),
		C.MemeFlag_CaseSensitive))
}

func (s *String) indexOfGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_indexOfWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) indexOfWithStartIndexGoStr(str string, startIndex int) int {
	b := []byte(str)
	return int(C.MemeString_indexOfWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) indexOfCStr(cstr *C.char, len C.size_t) int {
	return int(C.MemeString_indexOfWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

func (s *String) indexOfWithStartIndexCStr(cstr *C.char, len C.size_t, startIndex int) int {
	return int(C.MemeString_indexOfWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

func (s *String) lastIndexOfGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_lastIndexOfWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) lastIndexOfWithStartIndexGoStr(str string, startIndex int) int {
	b := []byte(str)
	return int(C.MemeString_lastIndexOfWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) lastIndexOfCStr(cstr *C.char, len C.size_t) int {
	return int(C.MemeString_lastIndexOfWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

func (s *String) lastIndexOfWithStartIndexCStr(cstr *C.char, len C.size_t, startIndex int) int {
	return int(C.MemeString_lastIndexOfWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

// match count
func (s *String) matchCountGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_matchCountWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) matchCountWithStartIndexGoStr(str string, startIndex int) int {
	b := []byte(str)
	return int(C.MemeString_matchCountWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

func (s *String) matchCountCStr(cstr *C.char, len C.size_t) int {
	return int(C.MemeString_matchCountWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(0),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

func (s *String) matchCountWithStartIndexCStr(cstr *C.char, len C.size_t, startIndex int) int {
	return int(C.MemeString_matchCountWithUtf8bytes(
		s.toMms(),
		C.MemeInteger_t(startIndex),
		(*C.uchar)(unsafe.Pointer(cstr)), C.MemeInteger_t(len),
		C.MemeFlag_CaseSensitive))
}

// start match
func (s *String) startMatchGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_startsMatchWithUtf8bytes(
		s.toMms(),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

// end match
func (s *String) endMatchGoStr(str string) int {
	b := []byte(str)
	return int(C.MemeString_endsMatchWithUtf8bytes(
		s.toMms(),
		(*C.uchar)(unsafe.Pointer(&b[0])), C.MemeInteger_t(len(b)),
		C.MemeFlag_CaseSensitive))
}

// swap
func (s *String) swap(other *String) {
	C.MemeString_swap(
		s.toMms(),
		other.toMms())
}

func (s *String) splitGoStr(str string) ([]*String, error) {
	stacks := make([]C.mms_stack_t, 4)
	stacksCount := C.MemeInteger_t(0)
	b := []byte(str)
	var list []*String
	for index := C.MemeInteger_t(0); index != C.MemeInteger_t(-1); {
		stacksCount = C.MemeInteger_t(len(stacks))
		result := C.MemeString_split(
			s.toMms(),
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

func (s *String) split(str *String) ([]*String, error) {
	stacks := make([]C.mms_stack_t, 4)
	stacksCount := C.MemeInteger_t(0)
	var list []*String
	for index := C.MemeInteger_t(0); index != C.MemeInteger_t(-1); {
		stacksCount = C.MemeInteger_t(len(stacks))
		result := C.MemeString_split(
			s.toMms(),
			str.cstr(), str.sizeToInteger(),
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

//export mmsinvoke_match_cond_byte_fn
func mmsinvoke_match_cond_byte_fn(b C.MemeByte_t, userdata unsafe.Pointer) C.MemeInteger_t {
	param := (*InvokeMatchParameter)(userdata)
	if param.condFn != nil {
		return C.MemeInteger_t(param.condFn(byte(b), param.userdata))
	}
	return C.MemeInteger_t(0)
}
