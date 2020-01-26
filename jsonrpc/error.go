package jsonrpc

type ErrorCode int

const (
	ParseError           ErrorCode = -32700
	InvalidRequest       ErrorCode = -32600
	MethodNotFound       ErrorCode = -32601
	InvalidParams        ErrorCode = -32602
	InternalError        ErrorCode = -32603
	ServerErrorStart     ErrorCode = -32099
	ServerErrorEnd       ErrorCode = -32000
	ServerNotInitialized ErrorCode = -32002
	UnknownErrorCode     ErrorCode = -32001
)

type ResponseError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewError(code ErrorCode, message string, data interface{}) *ResponseError {
	return &ResponseError{Code: code, Message: message, Data: data}
}
