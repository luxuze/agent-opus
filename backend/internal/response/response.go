package response

import (
	"encoding/json"
	"net/http"
)

// Response 统一的 HTTP 响应格式
type Response struct {
	Code    int         `json:"code"`    // 业务状态码：0=成功，非0=失败
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// WriteJSON 将响应写入 http.ResponseWriter
func (r *Response) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(r)
}

// ToJSON 转换为 JSON 字符串
func (r *Response) ToJSON() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
