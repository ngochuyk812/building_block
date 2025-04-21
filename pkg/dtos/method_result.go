package dtos

type MethodResult[R any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Result  R      `json:"result"`
}

func (m *MethodResult[R]) AddError(code int, mess string) {
	m.Code = code
	m.Message = mess
	m.Success = false
}
