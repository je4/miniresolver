package resolver

import (
	"github.com/je4/utils/v2/pkg/zLogger"
	"google.golang.org/grpc/resolver"
	"time"
)

func RegisterResolver(miniresolverClient *MiniResolver, checkTimeout, notFoundTimeout time.Duration, logger zLogger.ZLogger) {
	resolver.Register(NewMiniResolverResolverBuilder(miniresolverClient, checkTimeout, notFoundTimeout, logger))
}
