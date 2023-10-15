package reminder

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func moduleNameInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-service-name", "telegram-reminder")

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
