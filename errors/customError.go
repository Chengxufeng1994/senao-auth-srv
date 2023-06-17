package errors

type CustomError struct {
	Code    int
	Success bool
	Reason  string
}

func (c *CustomError) Error() string {
	return c.Reason
}

func New(code int, success bool, reason string) *CustomError {
	return &CustomError{
		Code:    code,
		Success: success,
		Reason:  reason,
	}
}
