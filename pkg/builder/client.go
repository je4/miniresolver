package builder

import (
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc/resolver"
	"time"
)

func RegisterResolver(miniresolverClient pb.MiniResolverClient, checkTimeout, notFoundTimeout time.Duration, logger zLogger.ZLogger) {
	resolver.Register(NewMiniResolverResolverBuilder(miniresolverClient, checkTimeout, notFoundTimeout, logger))
}
