package myerror

// BadRequestError 400:BadRequestエラー
type BadRequestError struct {
	Message string
	Err     error
}

func (e *BadRequestError) Error() string {
	return "Bad Request Error"
}

// UnauthorizedError 401:Unauthorizedエラー
type UnauthorizedError struct {
	Err error
}

func (e *UnauthorizedError) Error() string {
	return "Unauthorized Error"
}

// NotFoundError 404:NotFoundエラー
type NotFoundError struct {
	Err error
}

func (e *NotFoundError) Error() string {
	return "Not Found Error"
}

// InternalServerError 500:InternalServerエラー
type InternalServerError struct {
	Err error
}

func (e *InternalServerError) Error() string {
	return "Internal Server Error"
}
