package json

import "github.com/json-iterator/go"

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// 序列化
	Marshal = json.Marshal
	// 反序列化
	Unmarshal = json.Unmarshal
	//// 序列化。。。
	//MarshalIndent = json.MarshalIndent
	//// 解码
	//NewDecoder = json.NewDecoder
	//// 编码
	//NewEncoder = json.NewEncoder
)
