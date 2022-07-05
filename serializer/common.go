package serializer

type Response struct {
	Errno  int         `json:"errno"`
	ErrMsg string      `json:"errmsg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func BuildResponseOk(code int) Response {
	return Response{
		Errno:  code,
		ErrMsg: "",
	}
}

func BuildResponseOkWithData(code int, data interface{}) Response {
	return Response{
		Errno:  code,
		ErrMsg: "",
		Data:   data,
	}
}

func BuildResponseErr(code int, err error) Response {
	return Response{
		Errno:  code,
		ErrMsg: err.Error(),
	}
}
