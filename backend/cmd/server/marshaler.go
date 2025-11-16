package main

import (
	"encoding/json"
	"io"

	"agent-platform/internal/response"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

// CustomMarshaler 自定义的 Marshaler，将响应包装为统一格式 {code, message, data}
type CustomMarshaler struct {
	*runtime.JSONPb
}

// NewCustomMarshaler 创建自定义 Marshaler
func NewCustomMarshaler() *CustomMarshaler {
	return &CustomMarshaler{
		JSONPb: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true, // 使用 proto 原始字段名 (snake_case)
				EmitUnpopulated: true, // 输出零值字段
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true, // 忽略未知字段
			},
		},
	}
}

// Marshal 将 protobuf 消息序列化为 JSON，并包装为统一响应格式
func (m *CustomMarshaler) Marshal(v interface{}) ([]byte, error) {
	// 首先使用默认的 JSONPb 序列化
	data, err := m.JSONPb.Marshal(v)
	if err != nil {
		return nil, err
	}

	// 解析为通用对象
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	// 包装为统一响应格式
	wrappedResp := response.Success(obj)

	// 重新序列化
	return json.Marshal(wrappedResp)
}

// Unmarshal 反序列化 JSON 到 protobuf 消息（保持默认行为）
func (m *CustomMarshaler) Unmarshal(data []byte, v interface{}) error {
	return m.JSONPb.Unmarshal(data, v)
}

// NewDecoder 创建解码器（保持默认行为）
func (m *CustomMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return m.JSONPb.NewDecoder(r)
}

// NewEncoder 创建编码器（保持默认行为）
func (m *CustomMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return m.JSONPb.NewEncoder(w)
}

// ContentType 返回内容类型
func (m *CustomMarshaler) ContentType(v interface{}) string {
	return "application/json"
}

// Delimiter 返回流式响应的分隔符
func (m *CustomMarshaler) Delimiter() []byte {
	return m.JSONPb.Delimiter()
}
