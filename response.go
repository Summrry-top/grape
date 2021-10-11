package grape

const (
	Code200 = 200
	Code400 = 400
	Code404 = 404
	Code500 = 500

	Msg404 = "资源不存在"
)

var NIl = make([]int, 0)

type Response struct {
	Code int
	Msg  string
	Data interface{}
}

func Err404() *Response {
	return &Response{
		Code: Code404,
		Msg:  Msg404,
		Data: NIl,
	}
}
