package grpchelper

import (
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	miniresolver "github.com/je4/miniresolver/v2/pkg/resolver"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc/resolver"
	"time"
)

func RegisterResolver(miniresolverClient pb.MiniResolverClient, checkTimeout, notFoundTimeout time.Duration, logger zLogger.ZLogger) {
	resolver.Register(miniresolver.NewMiniResolverResolverBuilder(miniresolverClient, checkTimeout, notFoundTimeout, logger))
}
