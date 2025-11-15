package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	pb "agent-platform/gen/go"
	"agent-platform/internal/response"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// setupGateway 设置 gRPC-Gateway HTTP服务器
func setupGateway(grpcAddress string, httpPort string, logger *zap.Logger) error {
	// 使用 context.Background() 作为长期运行的 context
	// 不使用 cancel，因为 Gateway 需要在整个程序生命周期内运行
	ctx := context.Background()

	// 创建 gRPC-Gateway mux，使用自定义 Marshaler 和错误处理器
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, NewCustomMarshaler()),
		runtime.WithErrorHandler(customErrorHandler),
	)

	// 设置 gRPC 连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// 注册所有服务到 gRPC-Gateway
	if err := pb.RegisterAgentServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts); err != nil {
		return fmt.Errorf("failed to register AgentService: %w", err)
	}

	if err := pb.RegisterConversationServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts); err != nil {
		return fmt.Errorf("failed to register ConversationService: %w", err)
	}

	if err := pb.RegisterToolServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts); err != nil {
		return fmt.Errorf("failed to register ToolService: %w", err)
	}

	if err := pb.RegisterKnowledgeBaseServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts); err != nil {
		return fmt.Errorf("failed to register KnowledgeBaseService: %w", err)
	}

	// 添加 CORS 支持
	handler := cors(mux)

	// 启动 HTTP 服务器
	httpAddr := fmt.Sprintf(":%s", httpPort)
	logger.Info("Starting HTTP Gateway", zap.String("address", httpAddr))

	go func() {
		if err := http.ListenAndServe(httpAddr, handler); err != nil {
			logger.Fatal("Failed to serve HTTP Gateway", zap.Error(err))
		}
	}()

	return nil
}

// cors 添加 CORS 中间件
func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// customErrorHandler 自定义错误处理器，包装错误响应为 {code, message, data} 格式
func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	// 从 gRPC 错误中提取状态码和消息
	st := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(st.Code())

	// 包装为统一错误响应格式
	errResp := response.Error(int(st.Code()), st.Message())

	// 设置响应头和状态码
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	// 编码并写入响应
	_ = json.NewEncoder(w).Encode(errResp)
}
