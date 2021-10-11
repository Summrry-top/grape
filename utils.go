package grape

//func filterFlags(content string) string {
//	for i, char := range content {
//		if char == ' ' || char == ';' {
//			return content[:i]
//		}
//	}
//	return content
//}

func lastChar(str string) uint8 {
	size := len(str)
	if size == 0 {
		panic("The length of the string can't be 0")
	}
	return str[size-1]
}
