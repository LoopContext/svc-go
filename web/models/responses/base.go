package responses

type ResponseBase struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type Response struct {
	ResponseBase
	Error    []ResponseError `json:"errors"`
	Response interface{}     `json:"response"`
}

type ResponseError struct {
	ResponseBase
}
