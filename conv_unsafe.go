package conv

import "unsafe"

// UnsafeStrToBytes converts string to []byte without memory copy.
// Note that, if you completely sure you will never use <s> in the feature,
// you can use this unsafe function to implement type conversion in high performance.
func UnsafeStrToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeBytesToStr converts []byte to string without memory copy.
// Note that, if you completely sure you will never use <s> in the feature,
// you can use this unsafe function to implement type conversion in high performance.
func UnsafeBytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
