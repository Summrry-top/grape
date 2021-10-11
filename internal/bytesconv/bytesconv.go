package bytesconv

import "unsafe"

// 在不分配内存的情况下将字符串转换为字节片
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// 在不分配内存的情况下将字节片转换为字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
