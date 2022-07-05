package ginex

var (
	ErrUncaughtErr = ErrorResponse{
		code:    10005,
		message: "内部错误，请联系技术人员处理",
	}
)

type ErrorResponse struct {
	code    int
	message string
}

func (e ErrorResponse) Code() int {
	return e.code
}

func (e ErrorResponse) Message() string {
	return e.message
}
