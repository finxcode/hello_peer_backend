package global

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

type CustomErrors struct {
	BusinessError   CustomError
	ValidateError   CustomError
	TokenError      CustomError
	BadRequestError CustomError
	UserError       CustomError
}

var Errors = CustomErrors{
	BusinessError:   CustomError{40000, "业务错误"},
	ValidateError:   CustomError{42200, "请求参数错误"},
	TokenError:      CustomError{40100, "登录授权失效"},
	BadRequestError: CustomError{400, "请求syntax错误"},
	UserError:       CustomError{40101, "用户不存在"},
}
