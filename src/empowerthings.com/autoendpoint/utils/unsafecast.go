package utils

import (
	"reflect"
	"unsafe"
)

const (
	BYTES_IN_INT8  = 1
	BYTES_IN_INT16 = 2
	BYTES_IN_INT32 = 4
	BYTES_IN_INT64 = 8
)

func UnsafeCastInt8ToBytes(val int8) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT8, Cap: BYTES_IN_INT8}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCastInt16ToBytes(val int16) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT16, Cap: BYTES_IN_INT16}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCastInt32ToBytes(val int32) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT32, Cap: BYTES_IN_INT32}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCastInt64ToBytes(val int64) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT64, Cap: BYTES_IN_INT64}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}


func UnsafeCastUInt8ToBytes(val uint8) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT8, Cap: BYTES_IN_INT8}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCastUInt16ToBytes(val uint16) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT16, Cap: BYTES_IN_INT16}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCastUInt32ToBytes(val uint32) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT32, Cap: BYTES_IN_INT32}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCastUInt64ToBytes(val uint64) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT64, Cap: BYTES_IN_INT64}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}
