package grpchelper

import (
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	miniresolver "github.com/je4/miniresolver/v2/pkg/resolver"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc/resolver"
	"strings"
)

func RegisterResolver(miniresolverClient pb.MiniResolverClient, logger zLogger.ZLogger) {
	resolver.Register(miniresolver.NewMiniResolverResolverBuilder(miniresolverClient, logger))
}

func GetAddress(fullMethodName string) string {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/")
	parts := strings.Split(fullMethodName, "/")
	if len(parts) < 2 {
		return ""
	}
	return "miniresolver:" + parts[0]
}
