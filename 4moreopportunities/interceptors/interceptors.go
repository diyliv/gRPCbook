package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type interceptor struct{}

func NewInterceptor() *interceptor {
	return &interceptor{}
}

func (i *interceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//тут может идти какая-то бизнес логика
	// info.FullMethod
	// информация о RPC методе
	log.Printf("Method Info :%v\n", info.FullMethod)

	// мы вызываем обработчик, чтобы завершить нормалное выполнение RPC-вызова
	idunno, err := handler(ctx, req)
	log.Printf("Some Logic After Call: %v\n", idunno)
	return idunno, err
}
